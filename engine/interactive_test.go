package engine

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/prompt"
	sharedsession "github.com/OpenUdon/authoring/session"
)

type interactiveState struct {
	Goal string `json:"goal"`
	Op   string `json:"op,omitempty"`
}

type fakeInteractiveExtractor struct {
	draft func(context.Context, DraftRequest[interactiveState, string]) (interactiveState, error)
}

func (f fakeInteractiveExtractor) Kickoff(context.Context, string) (interactiveState, error) {
	return interactiveState{}, nil
}

func (f fakeInteractiveExtractor) Draft(ctx context.Context, req DraftRequest[interactiveState, string]) (interactiveState, error) {
	if f.draft != nil {
		return f.draft(ctx, req)
	}
	return req.Session, nil
}

func (f fakeInteractiveExtractor) Refine(_ context.Context, session interactiveState) (interactiveState, error) {
	return session, nil
}

func (f fakeInteractiveExtractor) Disambiguate(context.Context, string, []string) ([]string, error) {
	return []string{"b", "a"}, nil
}

func TestRunInteractiveOpeningDraftAutosaveTranscript(t *testing.T) {
	var out strings.Builder
	var autosaved []interactiveState
	var transcriptTurns []PromptTurn
	var transcriptEvents []Event
	artifact, err := RunInteractive[interactiveState, string, string](context.Background(), strings.NewReader("List widgets\n"), &out, InteractiveHooks[interactiveState, string, string]{
		Documents:     []string{"a", "b"},
		DefaultMode:   prompt.DefaultsSilent,
		OpeningPrompt: "Describe the goal",
		Extractor: fakeInteractiveExtractor{draft: func(_ context.Context, req DraftRequest[interactiveState, string]) (interactiveState, error) {
			if req.Opening != "List widgets" || len(req.Docs) != 2 {
				t.Fatalf("draft request = %#v", req)
			}
			return interactiveState{Goal: req.Opening, Op: "listWidgets"}, nil
		}},
		ApplyOpeningAnswer: func(state *interactiveState, answer string, _ []string) error {
			state.Goal = answer
			return nil
		},
		Autosave: func(state interactiveState) error {
			autosaved = append(autosaved, state)
			return nil
		},
		CheckReadiness: func(state interactiveState, _ []string) []ReadinessIssue {
			if state.Goal == "" {
				return []ReadinessIssue{{Severity: "blocking", Code: "missing_goal", Message: "goal"}}
			}
			if state.Op == "" {
				return []ReadinessIssue{{Severity: "blocking", Code: "missing_op", Message: "operation"}}
			}
			return nil
		},
		Ready: func(_ interactiveState, issues []ReadinessIssue) bool { return len(issues) == 0 },
		PlanQuestion: func(interactiveState, []string, []ReadinessIssue) InteractiveQuestion {
			return InteractiveQuestion{Prompt: "Operation ID", SuggestedAnswer: "listWidgets"}
		},
		ApplyAnswer: func(state *interactiveState, _ InteractiveQuestion, answer string, _ []string) error {
			state.Op = answer
			return nil
		},
		FinalConfirm: func(_ *PromptSession, state *interactiveState, _ []string, _ *[]Event) (string, error) {
			return state.Goal + ":" + state.Op, nil
		},
		FinalResultSummary: func(artifact string) any { return artifact },
		SaveTranscript: func(turns []PromptTurn, events []Event, _ string) error {
			transcriptTurns = turns
			transcriptEvents = events
			return nil
		},
	})
	if err != nil {
		t.Fatalf("RunInteractive returned error: %v", err)
	}
	if artifact != "List widgets:listWidgets" {
		t.Fatalf("artifact = %q", artifact)
	}
	if len(autosaved) < 2 {
		t.Fatalf("autosaves = %#v", autosaved)
	}
	if len(transcriptTurns) == 0 || len(transcriptEvents) == 0 {
		t.Fatalf("transcript turns=%#v events=%#v", transcriptTurns, transcriptEvents)
	}
	if transcriptTurns[0].Label == "" || transcriptTurns[0].Source == "" {
		t.Fatalf("turn = %#v, want durable prompt-turn metadata", transcriptTurns[0])
	}
	for _, event := range transcriptEvents {
		if event.Kind == "" || event.Type == "" {
			t.Fatalf("event = %#v, want kind and type compatibility fields", event)
		}
	}
	projected := TranscriptEvents(transcriptEvents)
	if len(projected) == 0 || projected[0].Type == "" {
		t.Fatalf("projected events = %#v, want durable transcript projection", projected)
	}
	wantOrder := []string{"readiness", "draft_attempt", "model_draft_call", "draft_success", "final_confirm", "final_generated_artifacts"}
	position := -1
	for wantIndex, want := range wantOrder {
		found := -1
		for i := position + 1; i < len(transcriptEvents); i++ {
			if transcriptEvents[i].Type == want {
				found = i
				break
			}
		}
		if found < 0 {
			t.Fatalf("event %q did not follow %q in transcript: %#v", want, wantOrder[:wantIndex], transcriptEvents)
		}
		position = found
	}
	for _, event := range transcriptEvents {
		if event.Type == "next_question_decision" {
			t.Fatalf("ready draft emitted a phantom next-question decision: %#v", transcriptEvents)
		}
	}
}

func TestRunInteractiveQuestionDraftAndNeedsInput(t *testing.T) {
	_, err := RunInteractive[interactiveState, string, string](context.Background(), strings.NewReader("\n"), nil, InteractiveHooks[interactiveState, string, string]{
		Opening: "goal",
		CheckReadiness: func(interactiveState, []string) []ReadinessIssue {
			return []ReadinessIssue{{Severity: "blocking", Code: "missing", Message: "missing"}}
		},
		Ready: func(_ interactiveState, issues []ReadinessIssue) bool { return len(issues) == 0 },
		PlanQuestion: func(interactiveState, []string, []ReadinessIssue) InteractiveQuestion {
			return InteractiveQuestion{Prompt: "Operation ID"}
		},
		ApplyAnswer: func(*interactiveState, InteractiveQuestion, string, []string) error { return nil },
		MaxAttempts: 1,
	})
	if err == nil || !strings.Contains(err.Error(), "requires operator input") {
		t.Fatalf("RunInteractive error = %v, want needs input", err)
	}

	artifact, err := RunInteractive[interactiveState, string, string](context.Background(), nil, nil, InteractiveHooks[interactiveState, string, string]{
		Opening: "goal",
		Extractor: fakeInteractiveExtractor{draft: func(context.Context, DraftRequest[interactiveState, string]) (interactiveState, error) {
			return interactiveState{Goal: "goal", Op: "drafted"}, nil
		}},
		CheckReadiness: func(state interactiveState, _ []string) []ReadinessIssue {
			if state.Op == "" {
				return []ReadinessIssue{{Severity: "blocking", Code: "missing", Message: "missing"}}
			}
			return nil
		},
		Ready: func(_ interactiveState, issues []ReadinessIssue) bool { return len(issues) == 0 },
		PlanQuestion: func(interactiveState, []string, []ReadinessIssue) InteractiveQuestion {
			return InteractiveQuestion{Prompt: "Operation ID"}
		},
		ApplyAnswer: func(*interactiveState, InteractiveQuestion, string, []string) error { return nil },
		DraftQuestion: func(context.Context, *interactiveState, []string, []ReadinessIssue, InteractiveQuestion) (bool, error) {
			return true, nil
		},
		FinalConfirm: func(_ *PromptSession, state *interactiveState, _ []string, _ *[]Event) (string, error) {
			return state.Op, nil
		},
	})
	if err != nil || artifact != "drafted" {
		t.Fatalf("RunInteractive artifact=%q err=%v", artifact, err)
	}
}

func TestRunInteractiveNilReadyDefaultsToNoReadinessIssues(t *testing.T) {
	artifact, err := RunInteractive[interactiveState, string, string](context.Background(), nil, nil, InteractiveHooks[interactiveState, string, string]{
		Opening:        "goal",
		Session:        interactiveState{Goal: "goal", Op: "list"},
		CheckReadiness: func(interactiveState, []string) []ReadinessIssue { return nil },
		FinalConfirm: func(_ *PromptSession, state *interactiveState, _ []string, _ *[]Event) (string, error) {
			return state.Goal + ":" + state.Op, nil
		},
	})
	if err != nil || artifact != "goal:list" {
		t.Fatalf("RunInteractive artifact=%q err=%v", artifact, err)
	}
}

func TestRunInteractiveLifecycleAndCancellation(t *testing.T) {
	root := t.TempDir()
	draftPath := filepath.Join(root, "draft.json")
	transcriptPath := filepath.Join(root, "transcript.json")
	var saved bool
	artifact, err := RunInteractiveWithLifecycle[interactiveState, string, string](context.Background(), nil, nil, InteractiveHooks[interactiveState, string, string]{
		Opening:        "loaded",
		Session:        interactiveState{Goal: "loaded", Op: "list"},
		CheckReadiness: func(interactiveState, []string) []ReadinessIssue { return nil },
		Ready:          func(interactiveState, []ReadinessIssue) bool { return true },
		FinalConfirm: func(_ *PromptSession, state *interactiveState, _ []string, _ *[]Event) (string, error) {
			return state.Goal + ":" + state.Op, nil
		},
	}, InteractiveLifecycleOptions[interactiveState, string, string]{
		DraftPath:            draftPath,
		TranscriptPath:       transcriptPath,
		DeleteDraftOnSuccess: true,
		LoadDraft: func(string) (interactiveState, bool, error) {
			return interactiveState{Goal: "loaded", Op: "read"}, true, nil
		},
		SaveDraft: func(string, interactiveState) error {
			saved = true
			return nil
		},
		DeleteDraft: func(path string) error {
			return os.WriteFile(path+".deleted", []byte("ok"), 0o600)
		},
	})
	if err != nil || artifact != "loaded:read" {
		t.Fatalf("RunInteractiveWithLifecycle artifact=%q err=%v", artifact, err)
	}
	if saved {
		t.Fatalf("draft saved unexpectedly for already-ready session")
	}
	if _, err := os.Stat(transcriptPath); err != nil {
		t.Fatalf("transcript not saved: %v", err)
	}
	if _, err := os.Stat(draftPath + ".deleted"); err != nil {
		t.Fatalf("draft delete hook not called: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = RunInteractive[interactiveState, string, string](ctx, nil, nil, InteractiveHooks[interactiveState, string, string]{
		Opening: "goal",
		CheckReadiness: func(interactiveState, []string) []ReadinessIssue {
			return []ReadinessIssue{{Severity: "blocking", Code: "missing", Message: "missing"}}
		},
	})
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("RunInteractive canceled error = %v", err)
	}
}

func TestRunInteractiveCachesPlannerAndReplansOnlyAfterMutation(t *testing.T) {
	for _, test := range []struct {
		name      string
		mutate    bool
		wantPlans int
		wantLabel string
	}{
		{name: "unchanged", wantPlans: 1, wantLabel: "Operation initial"},
		{name: "mutated", mutate: true, wantPlans: 2, wantLabel: "Operation revised"},
	} {
		t.Run(test.name, func(t *testing.T) {
			plans := 0
			appliedLabel := ""
			artifact, err := RunInteractive(context.Background(), nil, nil, InteractiveHooks[interactiveState, string, string]{
				Opening:     "goal",
				DefaultMode: prompt.DefaultsSilent,
				Extractor:   fakeInteractiveExtractor{},
				CheckReadiness: func(state interactiveState, _ []string) []ReadinessIssue {
					if state.Op == "applied" {
						return nil
					}
					return []ReadinessIssue{{Code: "missing", Severity: "blocking"}}
				},
				PlanQuestion: func(state interactiveState, _ []string, _ []ReadinessIssue) InteractiveQuestion {
					plans++
					suffix := "initial"
					if state.Op == "revised" {
						suffix = "revised"
					}
					return InteractiveQuestion{ID: "operation", Prompt: "Operation " + suffix, Recommendation: "list"}
				},
				DraftQuestion: func(_ context.Context, state *interactiveState, _ []string, _ []ReadinessIssue, _ InteractiveQuestion) (bool, error) {
					if test.mutate {
						state.Op = "revised"
						return true, nil
					}
					return false, nil
				},
				ApplyAnswer: func(state *interactiveState, question InteractiveQuestion, _ string, _ []string) error {
					appliedLabel = question.Prompt
					state.Op = "applied"
					return nil
				},
				FinalConfirm: func(_ *PromptSession, state *interactiveState, _ []string, _ *[]Event) (string, error) {
					return state.Op, nil
				},
			})
			if err != nil || artifact != "applied" {
				t.Fatalf("artifact=%q err=%v", artifact, err)
			}
			if plans != test.wantPlans || appliedLabel != test.wantLabel {
				t.Fatalf("plans=%d applied label=%q", plans, appliedLabel)
			}
		})
	}
}

func TestRunInteractiveClearsProvisionalFrontierAfterDraftHookError(t *testing.T) {
	plans := 0
	appliedLabel := ""
	draftErrors := 0
	artifact, err := RunInteractive(context.Background(), nil, nil, InteractiveHooks[interactiveState, string, string]{
		Opening:     "goal",
		DefaultMode: prompt.DefaultsSilent,
		Extractor: fakeInteractiveExtractor{draft: func(_ context.Context, req DraftRequest[interactiveState, string]) (interactiveState, error) {
			state := req.Session
			state.Op = "drafted"
			return state, nil
		}},
		CheckReadiness: func(state interactiveState, _ []string) []ReadinessIssue {
			if state.Op == "applied" {
				return nil
			}
			return []ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(state interactiveState, _ []string, _ []ReadinessIssue) InteractiveQuestion {
			plans++
			label := "Operation initial"
			if state.Op == "drafted" {
				label = "Operation drafted"
			}
			return InteractiveQuestion{ID: "operation", Prompt: label, Recommendation: "list"}
		},
		DraftQuestion: func(context.Context, *interactiveState, []string, []ReadinessIssue, InteractiveQuestion) (bool, error) {
			return false, errors.New("recoverable draft hook failure")
		},
		OnDraftError: func(error) {
			draftErrors++
		},
		ApplyAnswer: func(state *interactiveState, question InteractiveQuestion, _ string, _ []string) error {
			appliedLabel = question.Prompt
			state.Op = "applied"
			return nil
		},
		FinalConfirm: func(_ *PromptSession, state *interactiveState, _ []string, _ *[]Event) (string, error) {
			return state.Op, nil
		},
	})
	if err != nil || artifact != "applied" {
		t.Fatalf("artifact=%q err=%v", artifact, err)
	}
	if plans != 2 || appliedLabel != "Operation initial" || draftErrors != 1 {
		t.Fatalf("plans=%d applied label=%q draft errors=%d", plans, appliedLabel, draftErrors)
	}
}

func TestRunInteractiveTextOverridesAndNilContext(t *testing.T) {
	var out strings.Builder
	artifact, err := RunInteractive(context.Background(), strings.NewReader("goal\nmanual\n"), &out, InteractiveHooks[interactiveState, string, string]{
		OpeningLabel: "Intent",
		FrontierText: FrontierText{
			Heading:        "Choices %d",
			Recommendation: "Suggested",
			Rationale:      "Reason",
			Evidence:       "Sources",
		},
		CheckReadiness: func(state interactiveState, _ []string) []ReadinessIssue {
			if state.Op != "" {
				return nil
			}
			return []ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(interactiveState, []string, []ReadinessIssue) InteractiveQuestion {
			return InteractiveQuestion{ID: "op", Prompt: "Operation", Recommendation: "list", Forced: true, Rationale: "safer", EvidenceRefs: []string{"catalog"}}
		},
		ApplyAnswer: func(state *interactiveState, _ InteractiveQuestion, answer string, _ []string) error {
			state.Op = answer
			return nil
		},
		FinalConfirm: func(_ *PromptSession, state *interactiveState, _ []string, _ *[]Event) (string, error) {
			return state.Op, nil
		},
	})
	if err != nil || artifact != "manual" {
		t.Fatalf("artifact=%q err=%v", artifact, err)
	}
	for _, text := range []string{"Intent: ", "Choices 1", "Suggested: list", "Reason: safer", "Sources: catalog"} {
		if !strings.Contains(out.String(), text) {
			t.Fatalf("output missing %q: %s", text, out.String())
		}
	}
	if _, err := RunInteractive[interactiveState, string, string](nil, nil, nil, InteractiveHooks[interactiveState, string, string]{}); err == nil || !strings.Contains(err.Error(), "context is required") {
		t.Fatalf("nil interactive context error = %v", err)
	}
	if _, err := RunInteractiveWithLifecycle[interactiveState, string, string](nil, nil, nil, InteractiveHooks[interactiveState, string, string]{}, InteractiveLifecycleOptions[interactiveState, string, string]{}); err == nil || !strings.Contains(err.Error(), "context is required") {
		t.Fatalf("nil lifecycle context error = %v", err)
	}
}

func TestSavePromptTranscriptRejectsSensitiveTurnsAndCapableSession(t *testing.T) {
	secret := "must-not-appear"
	path := filepath.Join(t.TempDir(), "prompt.json")
	err := SavePromptTranscript(path, "", []PromptTurn{{Label: "Credential", Answer: secret, Sensitive: true}}, nil, nil)
	if err == nil || strings.Contains(err.Error(), secret) {
		t.Fatalf("sensitive turn error = %v", err)
	}
	err = SavePromptTranscript(path, "", []PromptTurn{{Label: "Credential", Answer: "[redacted]", Sensitive: true, Redacted: true}}, nil, sharedsession.State{
		Answers: []sharedsession.Answer{{Slot: "credential", Value: secret, Sensitive: true}},
	})
	if err == nil || strings.Contains(err.Error(), secret) {
		t.Fatalf("sensitive session error = %v", err)
	}
	if err := SavePromptTranscript(path, "", []PromptTurn{{Label: "Credential", Answer: "[redacted]", Sensitive: true, Redacted: true}}, nil, sharedsession.State{
		Answers: []sharedsession.Answer{{Slot: "credential", Value: "[redacted]", Sensitive: true, Redacted: true}},
	}); err != nil {
		t.Fatalf("redacted prompt transcript rejected: %v", err)
	}
}
