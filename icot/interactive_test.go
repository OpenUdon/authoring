package icot

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/prompt"
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
