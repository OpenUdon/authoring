package icot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/OpenUdon/authoring/prompt"
	readinesspkg "github.com/OpenUdon/authoring/readiness"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

var (
	// ErrNeedsInput reports that the loop cannot proceed without operator
	// input.
	ErrNeedsInput = errors.New("authoring needs input")
	// ErrCanceled reports that the loop was canceled.
	ErrCanceled = errors.New("authoring canceled")
)

// Options supplies product-specific hooks for the generic progressive loop.
type Options[S, D, A any] struct {
	Session     S
	Documents   []D
	MaxAttempts int
	DefaultMode prompt.DefaultMode

	Normalize       func(*S)
	Draft           func(context.Context, S, []D, []session.ReadinessIssue, int) (S, error)
	ShouldDraft     func(S, []D, []session.ReadinessIssue, int) bool
	AfterDraft      func(S) error
	RefreshDocs     func(context.Context, S, []D) ([]D, error)
	CheckReadiness  func(S, []D) []session.ReadinessIssue
	Ready           func(S, []session.ReadinessIssue) bool
	PlanQuestion    func(S, []D, []session.ReadinessIssue) Question
	ApplyAnswer     func(*S, Question, string, []D) error
	Autosave        func(S) error
	FinalConfirm    func(context.Context, *S, []D, *[]transcript.Event) (A, error)
	SummarizeDraft  func(S) any
	SummarizeResult func(A) any
	OnDraftError    func(error)
}

// Question is a product-neutral follow-up question plan.
type Question = readinesspkg.Question

// Result is the generic loop outcome.
type Result[S, A any] struct {
	Session   S
	Artifact  A
	Events    []transcript.Event
	Turns     []session.PromptTurn
	Completed bool
}

// Run executes the progressive loop.
func Run[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, opts Options[S, D, A]) (Result[S, A], error) {
	var result Result[S, A]
	if err := checkContext(ctx); err != nil {
		return result, err
	}
	prompts := prompt.NewSession(in, out)
	prompts.SetDefaultMode(opts.DefaultMode)
	attempts := opts.MaxAttempts
	if attempts <= 0 {
		attempts = 20
	}
	state := opts.Session
	docs := append([]D(nil), opts.Documents...)
	normalize(opts.Normalize, &state)
	var events []transcript.Event
	record := func(eventType, stage, message string, fields map[string]string) {
		events = append(events, transcript.Event{
			Type:    eventType,
			Stage:   stage,
			Message: strings.TrimSpace(message),
			Fields:  fields,
		})
	}

	for attempt := 1; attempt <= attempts; attempt++ {
		if err := checkContext(ctx); err != nil {
			result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
			return result, err
		}
		issues := readiness(opts.CheckReadiness, state, docs)
		record("readiness", "readiness", "", map[string]string{"issues": fmt.Sprint(len(issues)), "attempt": fmt.Sprint(attempt)})
		if ready(opts.Ready, state, issues) {
			return confirm(ctx, opts, prompts, state, docs, events)
		}
		if opts.RefreshDocs != nil {
			refreshed, err := opts.RefreshDocs(ctx, state, docs)
			if err != nil {
				result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
				return result, err
			}
			docs = append([]D(nil), refreshed...)
		}
		if shouldDraft(opts, state, docs, issues, attempt) {
			record("draft_attempt", "draft", "", map[string]string{"attempt": fmt.Sprint(attempt)})
			draft, err := opts.Draft(ctx, state, docs, issues, attempt)
			if err != nil {
				record("draft_error", "draft", err.Error(), map[string]string{"attempt": fmt.Sprint(attempt)})
				if opts.OnDraftError != nil {
					opts.OnDraftError(err)
				}
			} else {
				state = draft
				normalize(opts.Normalize, &state)
				if opts.Autosave != nil {
					if err := opts.Autosave(state); err != nil {
						result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
						return result, err
					}
				}
				record("draft_success", "draft", "", map[string]string{"attempt": fmt.Sprint(attempt)})
				if opts.SummarizeDraft != nil {
					record("draft_summary", "draft", fmt.Sprint(opts.SummarizeDraft(state)), map[string]string{"attempt": fmt.Sprint(attempt)})
				}
				if opts.AfterDraft != nil {
					if err := opts.AfterDraft(state); err != nil {
						result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
						return result, err
					}
				}
				issues = readiness(opts.CheckReadiness, state, docs)
				record("readiness", "readiness", "", map[string]string{"issues": fmt.Sprint(len(issues)), "attempt": fmt.Sprint(attempt), "after": "draft"})
				if ready(opts.Ready, state, issues) {
					return confirm(ctx, opts, prompts, state, docs, events)
				}
			}
		}
		question := normalizeQuestion(planQuestion(opts.PlanQuestion, state, docs, issues))
		if question.Prompt == "" {
			record("needs_input", "question", "no question planned", map[string]string{"attempt": fmt.Sprint(attempt)})
			result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
			return result, ErrNeedsInput
		}
		answer, source, err := answerQuestion(prompts, question)
		if err != nil {
			record("needs_input", "question", err.Error(), map[string]string{"attempt": fmt.Sprint(attempt), "question": question.ID})
			result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return result, ErrNeedsInput
			}
			return result, err
		}
		record("question_answered", "question", "", map[string]string{"attempt": fmt.Sprint(attempt), "question": question.ID, "source": source})
		if opts.ApplyAnswer != nil {
			if err := opts.ApplyAnswer(&state, question, answer, docs); err != nil {
				result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
				return result, err
			}
		}
		normalize(opts.Normalize, &state)
		if opts.Autosave != nil {
			if err := opts.Autosave(state); err != nil {
				result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
				return result, err
			}
		}
	}
	record("needs_input", "loop", "maximum attempts reached", map[string]string{"attempts": fmt.Sprint(attempts)})
	result = finishResult(state, result.Artifact, events, prompts.Turns(), false)
	return result, ErrNeedsInput
}

func confirm[S, D, A any](ctx context.Context, opts Options[S, D, A], prompts *prompt.Session, state S, docs []D, events []transcript.Event) (Result[S, A], error) {
	var zero A
	if opts.FinalConfirm == nil {
		return finishResult(state, zero, events, prompts.Turns(), true), nil
	}
	events = append(events, transcript.Event{Type: "final_confirm", Stage: "confirm"})
	artifact, err := opts.FinalConfirm(ctx, &state, docs, &events)
	if err != nil {
		return finishResult(state, zero, events, prompts.Turns(), false), err
	}
	if opts.SummarizeResult != nil {
		events = append(events, transcript.Event{Type: "final_result", Stage: "confirm", Message: fmt.Sprint(opts.SummarizeResult(artifact))})
	}
	return finishResult(state, artifact, events, prompts.Turns(), true), nil
}

func answerQuestion(prompts *prompt.Session, question Question) (string, string, error) {
	if question.AllowDefault && !question.Forced && strings.TrimSpace(question.DefaultAnswer) != "" {
		return strings.TrimSpace(question.DefaultAnswer), firstNonEmpty(question.DefaultSource, "default"), nil
	}
	switch {
	case question.Forced:
		answer, err := prompts.AskDefaultForced(question.Prompt, question.DefaultAnswer)
		return answer, "user", err
	case question.Required:
		answer, err := prompts.AskDefaultRequired(question.Prompt, question.DefaultAnswer)
		return answer, "user", err
	default:
		answer, err := prompts.AskDefault(question.Prompt, question.DefaultAnswer)
		return answer, "user", err
	}
}

func shouldDraft[S, D, A any](opts Options[S, D, A], state S, docs []D, issues []session.ReadinessIssue, attempt int) bool {
	if opts.Draft == nil {
		return false
	}
	if opts.ShouldDraft == nil {
		return true
	}
	return opts.ShouldDraft(state, docs, issues, attempt)
}

func readiness[S, D any](check func(S, []D) []session.ReadinessIssue, state S, docs []D) []session.ReadinessIssue {
	if check == nil {
		return nil
	}
	return session.Normalize(session.State{Readiness: check(state, docs)}).Readiness
}

func ready[S any](check func(S, []session.ReadinessIssue) bool, state S, issues []session.ReadinessIssue) bool {
	if check != nil {
		return check(state, issues)
	}
	return readinesspkg.Ready(issues)
}

func planQuestion[S, D any](plan func(S, []D, []session.ReadinessIssue) Question, state S, docs []D, issues []session.ReadinessIssue) Question {
	if plan == nil {
		return Question{}
	}
	return plan(state, docs, issues)
}

func normalizeQuestion(question Question) Question {
	return readinesspkg.NormalizeQuestion(question)
}

func normalize[S any](fn func(*S), state *S) {
	if fn != nil {
		fn(state)
	}
}

func checkContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", ErrCanceled, err)
	}
	return nil
}

func finishResult[S, A any](state S, artifact A, events []transcript.Event, turns []session.PromptTurn, completed bool) Result[S, A] {
	return Result[S, A]{
		Session:   state,
		Artifact:  artifact,
		Events:    append([]transcript.Event(nil), events...),
		Turns:     session.Normalize(session.State{Turns: turns}).Turns,
		Completed: completed,
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
