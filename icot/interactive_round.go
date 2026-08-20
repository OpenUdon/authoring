package icot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

// RunInteractive adapts the richer extractor and lifecycle hooks to the same
// frontier-round engine used by Run and RunRuntime.
func RunInteractive[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, hooks InteractiveHooks[S, D, A]) (A, error) {
	var zero A
	if err := checkContext(ctx); err != nil {
		return zero, err
	}
	if out == nil {
		out = io.Discard
	}
	prompts := NewPromptSession(in, out)
	prompts.SetDefaultMode(hooks.DefaultMode)
	prompts.SetMessages(hooks.Messages)
	extractor := hooks.Extractor
	if extractor == nil {
		extractor = NoopExtractor[S, D]{}
	}
	_, noopExtractor := extractor.(interface{ noopInteractiveExtractor() })
	state := hooks.Session
	docs := append([]D(nil), hooks.Documents...)
	normalizeInteractive(hooks.Normalize, &state)
	var events []Event
	record := func(kind string, data any) {
		events = append(events, Event{Kind: kind, Type: kind, Data: data})
	}

	opening := strings.TrimSpace(hooks.Opening)
	openingLabel := firstNonEmpty(hooks.OpeningLabel, "Goal")
	if opening == "" {
		if hooks.OpeningPrompt != "" {
			fmt.Fprintln(out, hooks.OpeningPrompt)
		}
		answer, err := prompts.Ask(openingLabel)
		if err != nil {
			return zero, err
		}
		opening = strings.TrimSpace(answer)
		if strings.EqualFold(opening, "cancel") {
			return zero, ErrCanceled
		}
		if hooks.ApplyOpeningAnswer != nil {
			if err := hooks.ApplyOpeningAnswer(&state, opening, docs); err != nil {
				return zero, err
			}
		}
		if !hooks.NoLLM && !noopExtractor {
			record("model_kickoff_call", map[string]any{"opening": opening})
			kickoff, kickoffErr := extractor.Kickoff(ctx, opening)
			if kickoffErr != nil {
				record("model_kickoff_error", kickoffErr.Error())
				if hooks.OnDraftError != nil {
					hooks.OnDraftError(kickoffErr)
				}
			} else if hooks.LooksLikeSession == nil || hooks.LooksLikeSession(kickoff) {
				if hooks.MergeDraft != nil {
					state = hooks.MergeDraft(state, kickoff, docs)
				} else {
					state = kickoff
				}
				normalizeInteractive(hooks.Normalize, &state)
				record("model_kickoff_result", map[string]any{"opening": opening})
			}
		}
		appendInteractiveEvents(&events, hooks.OpeningEvents, state)
		if hooks.RefreshDocuments != nil {
			refreshed, err := hooks.RefreshDocuments(state, docs)
			if err != nil {
				return zero, err
			}
			docs = append([]D(nil), refreshed...)
		}
		normalizeInteractive(hooks.Normalize, &state)
		if hooks.Autosave != nil {
			if err := hooks.Autosave(state); err != nil {
				return zero, err
			}
		}
		record("progressive_question", InteractiveQuestion{Prompt: openingLabel, Slots: []string{"goal"}})
		record("progressive_answer", PromptTurn{Label: openingLabel, Answer: answer, Source: "user"})
	}
	if !hooks.NoLLM && !noopExtractor && len(docs) > 1 && opening != "" {
		ranked, err := extractor.Disambiguate(ctx, opening, docs)
		if err == nil && hooks.RankDocuments != nil {
			docs = hooks.RankDocuments(docs, ranked)
		} else if err != nil && hooks.OnDraftError != nil {
			hooks.OnDraftError(fmt.Errorf("document ranking skipped: %w", err))
		}
	}
	if hooks.DeterministicPrefill != nil {
		hooks.DeterministicPrefill(&state, docs)
		normalizeInteractive(hooks.Normalize, &state)
	}

	opts := Options[S, D, A]{
		Session:             state,
		Documents:           docs,
		DefaultMode:         hooks.DefaultMode,
		FrontierText:        hooks.FrontierText,
		Interview:           hooks.Interview,
		MaxRounds:           hooks.MaxRounds,
		Normalize:           hooks.Normalize,
		ProgressFingerprint: hooks.ProgressFingerprint,
		CheckReadiness:      hooks.CheckReadiness,
		Ready:               hooks.Ready,
		Autosave:            hooks.Autosave,
		AfterDraft:          hooks.AfterDraft,
		OnDraftError:        hooks.OnDraftError,
		OnEvent: func(event transcript.Event) {
			events = append(events, Event{Kind: event.Type, Type: event.Type, Data: event})
			if hooks.OnEvent != nil {
				hooks.OnEvent(event)
			}
		},
	}
	if hooks.RefreshDocuments != nil {
		opts.RefreshDocs = func(_ context.Context, current S, currentDocs []D) ([]D, error) {
			return hooks.RefreshDocuments(current, currentDocs)
		}
	}
	var provisionalFrontier []Question
	var provisionalValid, provisionalMutated bool
	var displayedFrontier []Question
	clearProvisionalFrontier := func() {
		provisionalFrontier = nil
		provisionalValid = false
		provisionalMutated = false
	}
	if !hooks.NoLLM && !noopExtractor {
		opts.Draft = func(ctx context.Context, current S, currentDocs []D, issues []session.ReadinessIssue, _ int) (S, error) {
			clearProvisionalFrontier()
			shouldDraft := true
			if hooks.ShouldDraft != nil {
				shouldDraft = hooks.ShouldDraft(current, currentDocs, issues)
			}
			if shouldDraft {
				drafted, err := interactiveModelDraft(ctx, extractor, hooks, opening, prompts.Turns(), current, currentDocs, issues, "model_draft_call", record)
				if err != nil {
					if hooks.OnDraftError != nil {
						hooks.OnDraftError(err)
					}
				} else {
					current = drafted
				}
			}
			currentIssues := issues
			if hooks.CheckReadiness != nil {
				currentIssues = hooks.CheckReadiness(current, currentDocs)
			}
			questions, err := interactivePlanQuestions(hooks, &current, currentDocs, currentIssues)
			if err != nil {
				clearProvisionalFrontier()
				return current, err
			}
			questions = PlanFrontier(questions).Questions
			provisionalFrontier = append([]Question(nil), questions...)
			provisionalValid = true
			for _, question := range questions {
				if hooks.DraftQuestion != nil {
					drafted, err := hooks.DraftQuestion(ctx, &current, currentDocs, currentIssues, question)
					if err != nil {
						// The outer loop discards the entire provisional draft on
						// error, so its question plan must be discarded as well.
						clearProvisionalFrontier()
						return current, err
					}
					if drafted {
						provisionalMutated = true
						normalizeInteractive(hooks.Normalize, &current)
						record("question_draft_result", map[string]any{"question": question, "frontier": questions, "readiness_issues": currentIssues})
					}
				}
				if hooks.ShouldDraftQuestion != nil && hooks.ShouldDraftQuestion(current, currentDocs, currentIssues, question) {
					drafted, err := interactiveModelDraft(ctx, extractor, hooks, opening, prompts.Turns(), current, currentDocs, currentIssues, "model_question_draft_call", record)
					if err != nil {
						if hooks.OnDraftError != nil {
							hooks.OnDraftError(err)
						}
					} else {
						current = drafted
						provisionalMutated = true
					}
				}
			}
			return current, nil
		}
		opts.ShouldDraft = func(S, []D, []session.ReadinessIssue, int) bool { return true }
	}
	opts.planFrontierWithError = func(current *S, currentDocs []D, issues []session.ReadinessIssue) ([]Question, error) {
		record("readiness_decision", issues)
		questions := append([]Question(nil), provisionalFrontier...)
		if !provisionalValid || provisionalMutated {
			var err error
			questions, err = interactivePlanQuestions(hooks, current, currentDocs, issues)
			if err != nil {
				return nil, err
			}
		}
		questions = PlanFrontier(questions).Questions
		clearProvisionalFrontier()
		displayedFrontier = append([]Question(nil), questions...)
		for _, question := range questions {
			record("next_question_decision", question)
		}
		return questions, nil
	}
	if hooks.Interview == nil {
		opts.ApplyRound = func(current *S, answers []RoundAnswer, currentDocs []D) error {
			if hooks.ApplyRound != nil {
				return hooks.ApplyRound(current, answers, currentDocs)
			}
			if hooks.ApplyAnswer == nil {
				return fmt.Errorf("frontier answer hook is required")
			}
			byID := map[string]Question{}
			for _, question := range displayedFrontier {
				byID[question.ID] = question
			}
			for _, answer := range answers {
				question, ok := byID[answer.QuestionID]
				if !ok {
					return fmt.Errorf("answer references unknown displayed question %q", answer.QuestionID)
				}
				if err := hooks.ApplyAnswer(current, question, answer.Value, currentDocs); err != nil {
					return err
				}
			}
			return nil
		}
	}
	opts.AfterRound = func(current *S, answers []RoundAnswer, currentDocs []D) error {
		if hooks.AfterRound != nil {
			if err := hooks.AfterRound(current, currentDocs); err != nil {
				return err
			}
		}
		if hooks.DeterministicPrefill != nil {
			hooks.DeterministicPrefill(current, currentDocs)
		}
		normalizeInteractive(hooks.Normalize, current)
		for _, answer := range answers {
			record("progressive_answer", answer)
		}
		displayedFrontier = nil
		return nil
	}
	opts.finalConfirmWithPrompts = func(_ context.Context, shared *prompt.Session, current *S, currentDocs []D, _ *[]transcript.Event) (A, error) {
		if hooks.FinalQuestion != nil {
			question := readinessQuestion(hooks.FinalQuestion(*current, currentDocs))
			if question.Prompt != "" {
				record("next_question_decision", question)
			}
		}
		if hooks.FinalConfirm == nil {
			return zero, fmt.Errorf("final confirmation hook is required")
		}
		artifact, err := hooks.FinalConfirm(&PromptSession{session: shared}, current, currentDocs, &events)
		if err == nil && hooks.FinalResultSummary != nil {
			record("final_generated_artifacts", hooks.FinalResultSummary(artifact))
		}
		return artifact, err
	}

	result, err := runWithPromptSession(ctx, prompts.session, out, opts)
	if err != nil {
		if errors.Is(err, ErrNeedsInput) {
			return zero, fmt.Errorf("interactive iCoT requires operator input: %w", err)
		}
		return zero, err
	}
	if hooks.SaveTranscript != nil {
		if saveErr := hooks.SaveTranscript(result.Turns, events, result.Artifact); saveErr != nil {
			return result.Artifact, saveErr
		}
	}
	return result.Artifact, nil
}

func interactivePlanQuestions[S, D, A any](hooks InteractiveHooks[S, D, A], state *S, docs []D, issues []session.ReadinessIssue) ([]Question, error) {
	if hooks.Interview != nil {
		return hooks.Interview.Plan(state, docs)
	}
	if hooks.PlanFrontier != nil {
		return hooks.PlanFrontier(*state, docs, issues), nil
	}
	if hooks.PlanQuestion != nil {
		question := hooks.PlanQuestion(*state, docs, issues)
		if strings.TrimSpace(question.Prompt) != "" {
			return []Question{question}, nil
		}
	}
	return nil, nil
}

func readinessQuestion(question Question) Question {
	plan := PlanFrontier([]Question{question})
	if len(plan.Questions) == 0 {
		return Question{}
	}
	return plan.Questions[0]
}

func interactiveModelDraft[S, D, A any](ctx context.Context, extractor Extractor[S, D], hooks InteractiveHooks[S, D, A], opening string, turns []PromptTurn, state S, docs []D, issues []session.ReadinessIssue, kind string, record func(string, any)) (S, error) {
	request := DraftRequest[S, D]{
		Opening: opening, Brief: hooks.Brief, Session: state, Docs: docs,
		TranscriptTurns: turns, ReadinessFeedback: append([]ReadinessIssue(nil), issues...),
	}
	record(kind, map[string]any{"opening": request.Opening, "turn_count": len(turns), "readiness_issues": request.ReadinessFeedback})
	draft, err := extractor.Draft(ctx, request)
	if err != nil {
		record("model_draft_error", err.Error())
		return state, err
	}
	if hooks.LooksLikeSession != nil && !hooks.LooksLikeSession(draft) {
		return state, fmt.Errorf("model draft did not match the session contract")
	}
	if hooks.MergeDraft != nil {
		state = hooks.MergeDraft(state, draft, docs)
	} else {
		state = draft
	}
	normalizeInteractive(hooks.Normalize, &state)
	if hooks.DeterministicPrefill != nil {
		hooks.DeterministicPrefill(&state, docs)
		normalizeInteractive(hooks.Normalize, &state)
	}
	if hooks.DraftResultSummary != nil {
		record("model_draft_result", hooks.DraftResultSummary(state))
	}
	appendInteractiveEventsWithRecord(hooks.DraftEvents, state, record)
	return state, nil
}

func normalizeInteractive[S any](normalize func(*S), state *S) {
	if normalize != nil {
		normalize(state)
	}
}

func appendInteractiveEvents[S any](events *[]Event, source func(S) []Event, state S) {
	if source == nil {
		return
	}
	*events = append(*events, NormalizeEvents(source(state))...)
}

func appendInteractiveEventsWithRecord[S any](source func(S) []Event, state S, record func(string, any)) {
	if source == nil {
		return
	}
	for _, event := range NormalizeEvents(source(state)) {
		record(event.Type, event.Data)
	}
}
