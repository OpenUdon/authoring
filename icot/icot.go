package icot

import (
	"context"
	"encoding/json"
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
	// ErrNoProgress reports three consecutive rounds with no state or readiness
	// change. It is returned together with ErrNeedsInput.
	ErrNoProgress  = errors.New("authoring made no progress for three consecutive rounds")
	errBlankAnswer = errors.New("frontier question requires operator input")
)

const DefaultNoProgressLimit = 3

// RoundAnswer is one answer collected for a frontier question. A complete
// round is collected before any answer is applied.
type RoundAnswer struct {
	QuestionID string   `json:"question_id"`
	Slots      []string `json:"slots,omitempty"`
	Value      string   `json:"value,omitempty"`
	Source     string   `json:"source,omitempty"`
}

// Options supplies product-specific hooks for the generic frontier-round
// engine.
type Options[S, D, A any] struct {
	Session     S
	Documents   []D
	DefaultMode prompt.DefaultMode

	Normalize       func(*S)
	Draft           func(context.Context, S, []D, []session.ReadinessIssue, int) (S, error)
	ShouldDraft     func(S, []D, []session.ReadinessIssue, int) bool
	AfterDraft      func(S) error
	RefreshDocs     func(context.Context, S, []D) ([]D, error)
	CheckReadiness  func(S, []D) []session.ReadinessIssue
	Ready           func(S, []session.ReadinessIssue) bool
	PlanFrontier    func(S, []D, []session.ReadinessIssue) []Question
	ApplyRound      func(*S, []RoundAnswer, []D) error
	Autosave        func(S) error
	FinalConfirm    func(context.Context, *S, []D, *[]transcript.Event) (A, error)
	SummarizeDraft  func(S) any
	SummarizeResult func(A) any
	OnDraftError    func(error)
	onEvent         func(transcript.Event)

	// Deprecated: MaxAttempts is ignored. The engine has no breadth ceiling and
	// stops after three consecutive no-progress rounds.
	MaxAttempts int
	// Deprecated: NoProgressLimit is ignored; v2 always diagnoses the third
	// consecutive no-progress round.
	NoProgressLimit int
	// Deprecated: PlanQuestion and ApplyAnswer are adapted to one-question
	// frontiers for source compatibility. New adapters should use PlanFrontier
	// and ApplyRound.
	PlanQuestion func(S, []D, []session.ReadinessIssue) Question
	ApplyAnswer  func(*S, Question, string, []D) error

	// finalConfirmWithPrompts lets the interactive compatibility adapter share
	// the engine's prompt stream without exposing a second loop.
	finalConfirmWithPrompts func(context.Context, *prompt.Session, *S, []D, *[]transcript.Event) (A, error)
}

// Question is a product-neutral follow-up question plan.
type Question = readinesspkg.Question

// Result is the generic loop outcome.
type Result[S, A any] struct {
	Session          S                    `json:"session"`
	Artifact         A                    `json:"artifact,omitempty"`
	Events           []transcript.Event   `json:"events,omitempty"`
	Turns            []session.PromptTurn `json:"turns,omitempty"`
	Frontier         []Question           `json:"frontier,omitempty"`
	Answers          []RoundAnswer        `json:"answers,omitempty"`
	Rounds           int                  `json:"rounds,omitempty"`
	NoProgressRounds int                  `json:"no_progress_rounds,omitempty"`
	Completed        bool                 `json:"completed"`
}

// PlanFrontier normalizes and deterministically orders a full ready question
// frontier.
func PlanFrontier(questions []Question) readinesspkg.Plan {
	return readinesspkg.EvaluatePlan(questions)
}

// Run executes the dependency-ready frontier-round engine.
func Run[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, opts Options[S, D, A]) (Result[S, A], error) {
	prompts := prompt.NewSession(in, out)
	prompts.SetDefaultMode(opts.DefaultMode)
	return runWithPromptSession(ctx, prompts, out, opts)
}

func runWithPromptSession[S, D, A any](ctx context.Context, prompts *prompt.Session, out io.Writer, opts Options[S, D, A]) (Result[S, A], error) {
	var result Result[S, A]
	if err := checkContext(ctx); err != nil {
		return result, err
	}
	if out == nil {
		out = io.Discard
	}
	if prompts == nil {
		prompts = prompt.NewSession(nil, out)
		prompts.SetDefaultMode(opts.DefaultMode)
	}
	noProgressLimit := DefaultNoProgressLimit
	state := opts.Session
	docs := append([]D(nil), opts.Documents...)
	normalize(opts.Normalize, &state)
	var events []transcript.Event
	record := func(eventType, stage, message string, fields map[string]string) {
		event := transcript.Event{Type: eventType, Stage: stage, Message: strings.TrimSpace(message), Fields: fields}
		events = append(events, event)
		if opts.onEvent != nil {
			opts.onEvent(event)
		}
	}

	for round := 1; ; round++ {
		result.Rounds = round
		if err := checkContext(ctx); err != nil {
			return finishResult(state, result, events, prompts.Turns(), false), err
		}
		beforeIssues := readiness(opts.CheckReadiness, state, docs)
		beforeFingerprint := roundFingerprint(state, docs, beforeIssues)
		record("readiness", "readiness", "", map[string]string{"issues": fmt.Sprint(len(beforeIssues)), "round": fmt.Sprint(round)})
		if ready(opts.Ready, state, beforeIssues) {
			return confirm(ctx, opts, prompts, state, docs, events, result)
		}

		issues := beforeIssues
		if opts.RefreshDocs != nil {
			refreshed, err := opts.RefreshDocs(ctx, state, docs)
			if err != nil {
				return finishResult(state, result, events, prompts.Turns(), false), err
			}
			docs = append([]D(nil), refreshed...)
			issues = readiness(opts.CheckReadiness, state, docs)
		}
		if shouldDraft(opts, state, docs, issues, round) {
			record("draft_attempt", "draft", "", map[string]string{"round": fmt.Sprint(round)})
			draft, err := opts.Draft(ctx, state, docs, issues, round)
			if err != nil {
				record("draft_error", "draft", err.Error(), map[string]string{"round": fmt.Sprint(round)})
				if opts.OnDraftError != nil {
					opts.OnDraftError(err)
				}
			} else {
				state = draft
				normalize(opts.Normalize, &state)
				if opts.Autosave != nil {
					if err := opts.Autosave(state); err != nil {
						return finishResult(state, result, events, prompts.Turns(), false), err
					}
				}
				record("draft_success", "draft", "", map[string]string{"round": fmt.Sprint(round)})
				if opts.SummarizeDraft != nil {
					record("draft_summary", "draft", fmt.Sprint(opts.SummarizeDraft(state)), map[string]string{"round": fmt.Sprint(round)})
				}
				if opts.AfterDraft != nil {
					if err := opts.AfterDraft(state); err != nil {
						return finishResult(state, result, events, prompts.Turns(), false), err
					}
				}
			}
			issues = readiness(opts.CheckReadiness, state, docs)
			record("readiness", "readiness", "", map[string]string{"issues": fmt.Sprint(len(issues)), "round": fmt.Sprint(round), "after": "draft"})
			if ready(opts.Ready, state, issues) {
				return confirm(ctx, opts, prompts, state, docs, events, result)
			}
		}

		plan := PlanFrontier(planQuestions(opts, state, docs, issues))
		result.Frontier = append([]Question(nil), plan.Questions...)
		record("frontier_planned", "question", "", map[string]string{
			"round": fmt.Sprint(round), "questions": fmt.Sprint(len(plan.Questions)),
		})
		if len(plan.Questions) > 0 {
			displayFrontier(out, opts.DefaultMode, round, plan.Questions)
			answers, err := collectRoundAnswers(prompts, plan.Questions)
			if err != nil {
				fields := map[string]string{"round": fmt.Sprint(round)}
				if len(answers) < len(plan.Questions) {
					fields["question"] = plan.Questions[len(answers)].ID
				}
				record("needs_input", "question", err.Error(), fields)
				result.Answers = append([]RoundAnswer(nil), answers...)
				result = finishResult(state, result, events, prompts.Turns(), false)
				if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, errBlankAnswer) {
					return result, ErrNeedsInput
				}
				return result, err
			}
			for _, answer := range answers {
				if strings.EqualFold(strings.TrimSpace(answer.Value), "cancel") {
					result.Answers = append([]RoundAnswer(nil), answers...)
					return finishResult(state, result, events, prompts.Turns(), false), ErrCanceled
				}
			}
			if err := applyRound(opts, &state, answers, docs, plan.Questions); err != nil {
				return finishResult(state, result, events, prompts.Turns(), false), err
			}
			normalize(opts.Normalize, &state)
			if opts.Autosave != nil {
				if err := opts.Autosave(state); err != nil {
					return finishResult(state, result, events, prompts.Turns(), false), err
				}
			}
			result.Answers = append([]RoundAnswer(nil), answers...)
			record("round_applied", "question", "", map[string]string{"round": fmt.Sprint(round), "answers": fmt.Sprint(len(answers))})
		}

		afterIssues := readiness(opts.CheckReadiness, state, docs)
		progress := beforeFingerprint != roundFingerprint(state, docs, afterIssues)
		if progress {
			result.NoProgressRounds = 0
		} else {
			result.NoProgressRounds++
			record("round_no_progress", "loop", "round made no progress", map[string]string{
				"round": fmt.Sprint(round), "consecutive": fmt.Sprint(result.NoProgressRounds),
			})
			if result.NoProgressRounds >= noProgressLimit {
				record("needs_input", "loop", ErrNoProgress.Error(), map[string]string{"rounds": fmt.Sprint(result.NoProgressRounds)})
				result = finishResult(state, result, events, prompts.Turns(), false)
				return result, errors.Join(ErrNeedsInput, ErrNoProgress)
			}
		}
	}
}

func confirm[S, D, A any](ctx context.Context, opts Options[S, D, A], prompts *prompt.Session, state S, docs []D, events []transcript.Event, result Result[S, A]) (Result[S, A], error) {
	var zero A
	if opts.FinalConfirm == nil && opts.finalConfirmWithPrompts == nil {
		return finishResult(state, zeroResult(result, zero), events, prompts.Turns(), true), nil
	}
	event := transcript.Event{Type: "final_confirm", Stage: "confirm"}
	events = append(events, event)
	if opts.onEvent != nil {
		opts.onEvent(event)
	}
	var artifact A
	var err error
	priorEvents := len(events)
	if opts.finalConfirmWithPrompts != nil {
		artifact, err = opts.finalConfirmWithPrompts(ctx, prompts, &state, docs, &events)
	} else {
		artifact, err = opts.FinalConfirm(ctx, &state, docs, &events)
	}
	if opts.onEvent != nil {
		for _, added := range events[priorEvents:] {
			opts.onEvent(added)
		}
	}
	if err != nil {
		return finishResult(state, zeroResult(result, zero), events, prompts.Turns(), false), err
	}
	result.Artifact = artifact
	result.Frontier = nil
	if opts.SummarizeResult != nil {
		event = transcript.Event{Type: "final_result", Stage: "confirm", Message: fmt.Sprint(opts.SummarizeResult(artifact))}
		events = append(events, event)
		if opts.onEvent != nil {
			opts.onEvent(event)
		}
	}
	return finishResult(state, result, events, prompts.Turns(), true), nil
}

func zeroResult[S, A any](result Result[S, A], artifact A) Result[S, A] {
	result.Artifact = artifact
	return result
}

func displayFrontier(out io.Writer, mode prompt.DefaultMode, round int, questions []Question) {
	visible := make([]Question, 0, len(questions))
	for _, question := range questions {
		question = readinesspkg.NormalizeQuestion(question)
		if mode == prompt.DefaultsSilent && !question.Forced && question.Recommendation != "" {
			continue
		}
		visible = append(visible, question)
	}
	if len(visible) == 0 {
		return
	}
	fmt.Fprintf(out, "Round %d decisions:\n", round)
	for i, question := range visible {
		fmt.Fprintf(out, "%d. %s\n", i+1, question.Prompt)
		if question.Recommendation != "" {
			fmt.Fprintf(out, "   Recommendation: %s\n", prompt.OneLine(question.Recommendation))
		}
		if question.Rationale != "" {
			fmt.Fprintf(out, "   Why: %s\n", prompt.OneLine(question.Rationale))
		}
		if len(question.EvidenceRefs) > 0 {
			fmt.Fprintf(out, "   Evidence: %s\n", strings.Join(question.EvidenceRefs, ", "))
		}
	}
}

func collectRoundAnswers(prompts *prompt.Session, questions []Question) ([]RoundAnswer, error) {
	answers := make([]RoundAnswer, 0, len(questions))
	for _, question := range questions {
		question = readinesspkg.NormalizeQuestion(question)
		turnCount := len(prompts.Turns())
		value, err := answerQuestion(prompts, question)
		if err != nil {
			return answers, err
		}
		if strings.TrimSpace(value) == "" && question.Recommendation == "" {
			return answers, errBlankAnswer
		}
		source := DefaultRecommendationSource
		turns := prompts.Turns()
		if len(turns) > turnCount {
			source = turns[len(turns)-1].Source
		}
		answers = append(answers, RoundAnswer{
			QuestionID: question.ID,
			Slots:      append([]string(nil), question.Slots...),
			Value:      strings.TrimSpace(value),
			Source:     source,
		})
	}
	return answers, nil
}

const DefaultRecommendationSource = readinesspkg.DefaultRecommendationSource

func answerQuestion(prompts *prompt.Session, question Question) (string, error) {
	question = readinesspkg.NormalizeQuestion(question)
	switch {
	case question.Forced:
		return prompts.AskDefaultForced(question.Prompt, question.Recommendation)
	case question.Required && question.Recommendation == "":
		return prompts.AskDefaultRequired(question.Prompt, "")
	default:
		return prompts.AskDefault(question.Prompt, question.Recommendation)
	}
}

func shouldDraft[S, D, A any](opts Options[S, D, A], state S, docs []D, issues []session.ReadinessIssue, round int) bool {
	if opts.Draft == nil {
		return false
	}
	if opts.ShouldDraft == nil {
		return true
	}
	return opts.ShouldDraft(state, docs, issues, round)
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

func planQuestions[S, D, A any](opts Options[S, D, A], state S, docs []D, issues []session.ReadinessIssue) []Question {
	if opts.PlanFrontier != nil {
		return opts.PlanFrontier(state, docs, issues)
	}
	if opts.PlanQuestion != nil {
		question := readinesspkg.NormalizeQuestion(opts.PlanQuestion(state, docs, issues))
		if question.Prompt != "" {
			return []Question{question}
		}
	}
	return nil
}

func applyRound[S, D, A any](opts Options[S, D, A], state *S, answers []RoundAnswer, docs []D, questions []Question) error {
	if opts.ApplyRound != nil {
		return opts.ApplyRound(state, answers, docs)
	}
	if opts.ApplyAnswer == nil {
		return fmt.Errorf("frontier answer hook is required")
	}
	byID := make(map[string]Question, len(questions))
	for _, question := range questions {
		byID[question.ID] = question
	}
	for _, answer := range answers {
		question, ok := byID[answer.QuestionID]
		if !ok && len(questions) == 1 {
			question = questions[0]
		}
		if err := opts.ApplyAnswer(state, question, answer.Value, docs); err != nil {
			return err
		}
	}
	return nil
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

func finishResult[S, A any](state S, result Result[S, A], events []transcript.Event, turns []session.PromptTurn, completed bool) Result[S, A] {
	result.Session = state
	result.Events = append([]transcript.Event(nil), events...)
	result.Turns = session.Normalize(session.State{Turns: turns}).Turns
	result.Completed = completed
	return result
}

func roundFingerprint[S, D any](state S, docs []D, issues []session.ReadinessIssue) string {
	data, err := json.Marshal(struct {
		State  S                        `json:"state"`
		Docs   []D                      `json:"docs"`
		Issues []session.ReadinessIssue `json:"issues"`
	}{State: state, Docs: docs, Issues: issues})
	if err == nil {
		return string(data)
	}
	return fmt.Sprintf("%#v|%#v|%#v", state, docs, issues)
}
