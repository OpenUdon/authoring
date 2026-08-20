package icot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/OpenUdon/authoring/internal/cancellation"
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
	ErrCanceled = cancellation.ErrCanceled
	// ErrNoProgress reports three consecutive rounds with no state or readiness
	// change. It is returned together with ErrNeedsInput.
	ErrNoProgress = errors.New("authoring made no progress for three consecutive rounds")
	// ErrRoundLimit reports that the emergency round fuse was reached before
	// the session became ready. It is returned together with ErrNeedsInput.
	ErrRoundLimit  = errors.New("authoring round limit reached")
	errBlankAnswer = errors.New("frontier question requires operator input")
)

const (
	DefaultNoProgressLimit = 3
	// DefaultMaxRounds is the emergency loop fuse selected by a zero MaxRounds.
	DefaultMaxRounds = 1000
)

// RoundAnswer is one answer collected for a frontier question. A complete
// round is collected before any answer is applied.
type RoundAnswer struct {
	QuestionID string   `json:"question_id"`
	Slots      []string `json:"slots,omitempty"`
	Value      string   `json:"value,omitempty"`
	Source     string   `json:"source,omitempty"`
}

// FrontierText customizes the generic labels used to display a question
// frontier. Heading is a fmt-style string receiving the round number.
type FrontierText struct {
	Heading        string
	Recommendation string
	Rationale      string
	Evidence       string

	// Compatibility aliases are normalized into the fields above.
	RoundHeading        string
	RecommendationLabel string
	RationaleLabel      string
	EvidenceLabel       string
}

// Options supplies product-specific hooks for the generic frontier-round
// engine.
type Options[S, D, A any] struct {
	Session      S
	Documents    []D
	DefaultMode  prompt.DefaultMode
	FrontierText FrontierText
	Interview    *InterviewBinding[S, D]
	MaxRounds    int

	Normalize           func(*S)
	ProgressFingerprint func(S, []D, []session.ReadinessIssue) (string, error)
	Draft               func(context.Context, S, []D, []session.ReadinessIssue, int) (S, error)
	ShouldDraft         func(S, []D, []session.ReadinessIssue, int) bool
	AfterDraft          func(S) error
	RefreshDocs         func(context.Context, S, []D) ([]D, error)
	CheckReadiness      func(S, []D) []session.ReadinessIssue
	Ready               func(S, []session.ReadinessIssue) bool
	PlanFrontier        func(S, []D, []session.ReadinessIssue) []Question
	ApplyRound          func(*S, []RoundAnswer, []D) error
	AfterRound          func(*S, []RoundAnswer, []D) error
	Autosave            func(S) error
	FinalConfirm        func(context.Context, *S, []D, *[]transcript.Event) (A, error)
	SummarizeDraft      func(S) any
	SummarizeResult     func(A) any
	OnDraftError        func(error)
	// OnEvent receives every engine event after its monotonic ID is assigned.
	OnEvent func(transcript.Event)

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
	planFrontierWithError   func(*S, []D, []session.ReadinessIssue) ([]Question, error)
}

// Question is a product-neutral follow-up question plan.
type Question = readinesspkg.Question

// Result is the generic loop outcome.
type Result[S, A any] struct {
	Session  S                    `json:"session"`
	Artifact A                    `json:"artifact,omitempty"`
	Events   []transcript.Event   `json:"events,omitempty"`
	Turns    []session.PromptTurn `json:"turns,omitempty"`
	Frontier []Question           `json:"frontier,omitempty"`
	// Answers contains the most recently collected round. It is partial when
	// the result needs input; otherwise it is the last successfully applied
	// round.
	Answers          []RoundAnswer `json:"answers,omitempty"`
	Rounds           int           `json:"rounds,omitempty"`
	NoProgressRounds int           `json:"no_progress_rounds,omitempty"`
	Completed        bool          `json:"completed"`
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
	maxRounds, err := resolveMaxRounds(opts.MaxRounds)
	if err != nil {
		return result, err
	}
	if out == nil {
		out = io.Discard
	}
	if prompts == nil {
		return result, fmt.Errorf("icot prompt session is required")
	}
	noProgressLimit := DefaultNoProgressLimit
	state := opts.Session
	docs := append([]D(nil), opts.Documents...)
	normalize(opts.Normalize, &state)
	events := eventRecorder(opts.OnEvent)
	if opts.Interview != nil {
		result.NoProgressRounds = opts.Interview.noProgressRounds(&state)
	}

	for round := 1; ; round++ {
		if round > maxRounds {
			issues := readiness(opts.CheckReadiness, state, docs)
			if ready(opts.Ready, state, issues) {
				return confirm(ctx, opts, prompts, state, docs, events, result)
			}
			events.Record("round_limit", "loop", ErrRoundLimit.Error(), map[string]string{"rounds": fmt.Sprint(maxRounds)})
			events.Record("needs_input", "loop", ErrRoundLimit.Error(), map[string]string{"rounds": fmt.Sprint(maxRounds)})
			result = finishResult(state, result, events.Events(), prompts.Turns(), false)
			return result, errors.Join(ErrNeedsInput, ErrRoundLimit)
		}
		result.Rounds = round
		if err := checkContext(ctx); err != nil {
			return finishResult(state, result, events.Events(), prompts.Turns(), false), err
		}
		beforeIssues := readiness(opts.CheckReadiness, state, docs)
		beforeFingerprint, err := progressFingerprint(opts.ProgressFingerprint, state, docs, beforeIssues)
		if err != nil {
			return finishResult(state, result, events.Events(), prompts.Turns(), false), fmt.Errorf("compute progress fingerprint before round %d: %w", round, err)
		}
		events.Record("readiness", "readiness", "", map[string]string{"issues": fmt.Sprint(len(beforeIssues)), "round": fmt.Sprint(round)})
		if ready(opts.Ready, state, beforeIssues) {
			return confirm(ctx, opts, prompts, state, docs, events, result)
		}

		issues := beforeIssues
		if opts.RefreshDocs != nil {
			refreshed, err := opts.RefreshDocs(ctx, state, docs)
			if err != nil {
				return finishResult(state, result, events.Events(), prompts.Turns(), false), err
			}
			docs = append([]D(nil), refreshed...)
			issues = readiness(opts.CheckReadiness, state, docs)
		}
		if shouldDraft(opts, state, docs, issues, round) {
			events.Record("draft_attempt", "draft", "", map[string]string{"round": fmt.Sprint(round)})
			draft, err := opts.Draft(ctx, state, docs, issues, round)
			if err != nil {
				events.Record("draft_error", "draft", err.Error(), map[string]string{"round": fmt.Sprint(round)})
				if opts.OnDraftError != nil {
					opts.OnDraftError(err)
				}
			} else {
				state = draft
				normalize(opts.Normalize, &state)
				if opts.Autosave != nil {
					if err := opts.Autosave(state); err != nil {
						return finishResult(state, result, events.Events(), prompts.Turns(), false), err
					}
				}
				events.Record("draft_success", "draft", "", map[string]string{"round": fmt.Sprint(round)})
				if opts.SummarizeDraft != nil {
					events.Record("draft_summary", "draft", fmt.Sprint(opts.SummarizeDraft(state)), map[string]string{"round": fmt.Sprint(round)})
				}
				if opts.AfterDraft != nil {
					if err := opts.AfterDraft(state); err != nil {
						return finishResult(state, result, events.Events(), prompts.Turns(), false), err
					}
				}
			}
			issues = readiness(opts.CheckReadiness, state, docs)
			events.Record("readiness", "readiness", "", map[string]string{"issues": fmt.Sprint(len(issues)), "round": fmt.Sprint(round), "after": "draft"})
			if ready(opts.Ready, state, issues) {
				return confirm(ctx, opts, prompts, state, docs, events, result)
			}
		}

		questions, err := planQuestions(opts, &state, docs, issues)
		if err != nil {
			return finishResult(state, result, events.Events(), prompts.Turns(), false), err
		}
		plan := PlanFrontier(questions)
		result.Frontier = append([]Question(nil), plan.Questions...)
		roundApplied := false
		events.Record("frontier_planned", "question", "", map[string]string{
			"round": fmt.Sprint(round), "questions": fmt.Sprint(len(plan.Questions)),
		})
		if len(plan.Questions) > 0 {
			displayFrontier(out, opts.DefaultMode, round, plan.Questions, opts.FrontierText)
			answers, err := collectRoundAnswers(prompts, plan.Questions)
			if err != nil {
				fields := map[string]string{"round": fmt.Sprint(round)}
				if len(answers) < len(plan.Questions) {
					fields["question"] = plan.Questions[len(answers)].ID
				}
				events.Record("needs_input", "question", err.Error(), fields)
				result.Answers = append([]RoundAnswer(nil), answers...)
				result = finishResult(state, result, events.Events(), prompts.Turns(), false)
				if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, errBlankAnswer) {
					return result, ErrNeedsInput
				}
				return result, err
			}
			for _, answer := range answers {
				if isUserAnswer(answer.Source) && strings.EqualFold(strings.TrimSpace(answer.Value), "cancel") {
					return finishResult(state, result, events.Events(), prompts.Turns(), false), ErrCanceled
				}
			}
			if err := applyRound(opts, &state, answers, docs, plan.Questions); err != nil {
				return finishResult(state, result, events.Events(), prompts.Turns(), false), err
			}
			normalize(opts.Normalize, &state)
			roundApplied = true
			result.Answers = append([]RoundAnswer(nil), answers...)
			events.Record("round_applied", "question", "", map[string]string{"round": fmt.Sprint(round), "answers": fmt.Sprint(len(answers))})
		}

		afterIssues := readiness(opts.CheckReadiness, state, docs)
		// ApplyRound resets the durable interview counter transactionally. Put
		// back its pre-round value only while comparing fingerprints so counter
		// bookkeeping cannot manufacture semantic progress. The authoritative
		// post-round value is assigned immediately below.
		if opts.Interview != nil {
			opts.Interview.setNoProgressRounds(&state, result.NoProgressRounds)
		}
		afterFingerprint, err := progressFingerprint(opts.ProgressFingerprint, state, docs, afterIssues)
		if err != nil {
			return finishResult(state, result, events.Events(), prompts.Turns(), false), fmt.Errorf("compute progress fingerprint after round %d: %w", round, err)
		}
		progress := beforeFingerprint != afterFingerprint
		if progress {
			result.NoProgressRounds = 0
			if opts.Interview != nil {
				opts.Interview.setNoProgressRounds(&state, 0)
			}
		} else {
			result.NoProgressRounds++
			if opts.Interview != nil {
				opts.Interview.setNoProgressRounds(&state, result.NoProgressRounds)
			}
			events.Record("round_no_progress", "loop", "round made no progress", map[string]string{
				"round": fmt.Sprint(round), "consecutive": fmt.Sprint(result.NoProgressRounds),
			})
			if result.NoProgressRounds >= noProgressLimit {
				if roundApplied && opts.Autosave != nil {
					if err := opts.Autosave(state); err != nil {
						return finishResult(state, result, events.Events(), prompts.Turns(), false), err
					}
				}
				events.Record("needs_input", "loop", ErrNoProgress.Error(), map[string]string{"rounds": fmt.Sprint(result.NoProgressRounds)})
				result = finishResult(state, result, events.Events(), prompts.Turns(), false)
				return result, errors.Join(ErrNeedsInput, ErrNoProgress)
			}
		}
		if roundApplied && opts.Autosave != nil {
			if err := opts.Autosave(state); err != nil {
				return finishResult(state, result, events.Events(), prompts.Turns(), false), err
			}
		}
	}
}

func confirm[S, D, A any](ctx context.Context, opts Options[S, D, A], prompts *prompt.Session, state S, docs []D, events *eventLog, result Result[S, A]) (Result[S, A], error) {
	var zero A
	if opts.FinalConfirm == nil && opts.finalConfirmWithPrompts == nil {
		return finishResult(state, zeroResult(result, zero), events.Events(), prompts.Turns(), true), nil
	}
	events.Record("final_confirm", "confirm", "", nil)
	var artifact A
	var err error
	priorEvents := len(events.events)
	if opts.finalConfirmWithPrompts != nil {
		artifact, err = opts.finalConfirmWithPrompts(ctx, prompts, &state, docs, &events.events)
	} else {
		artifact, err = opts.FinalConfirm(ctx, &state, docs, &events.events)
	}
	events.SequenceFrom(priorEvents)
	if err != nil {
		return finishResult(state, zeroResult(result, zero), events.Events(), prompts.Turns(), false), err
	}
	result.Artifact = artifact
	result.Frontier = nil
	if opts.SummarizeResult != nil {
		events.Record("final_result", "confirm", fmt.Sprint(opts.SummarizeResult(artifact)), nil)
	}
	return finishResult(state, result, events.Events(), prompts.Turns(), true), nil
}

func zeroResult[S, A any](result Result[S, A], artifact A) Result[S, A] {
	result.Artifact = artifact
	return result
}

func displayFrontier(out io.Writer, mode prompt.DefaultMode, round int, questions []Question, text FrontierText) {
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
	text = normalizeFrontierText(text)
	heading := text.Heading
	if strings.Contains(heading, "%") {
		heading = fmt.Sprintf(heading, round)
	}
	fmt.Fprintln(out, heading)
	for i, question := range visible {
		fmt.Fprintf(out, "%d. %s\n", i+1, question.Prompt)
		if question.Recommendation != "" {
			fmt.Fprintf(out, "   %s: %s\n", text.Recommendation, prompt.OneLine(question.Recommendation))
		}
		if question.Rationale != "" {
			fmt.Fprintf(out, "   %s: %s\n", text.Rationale, prompt.OneLine(question.Rationale))
		}
		if len(question.EvidenceRefs) > 0 {
			fmt.Fprintf(out, "   %s: %s\n", text.Evidence, strings.Join(question.EvidenceRefs, ", "))
		}
	}
}

func normalizeFrontierText(text FrontierText) FrontierText {
	text.Heading = firstNonEmpty(text.Heading, text.RoundHeading, "Round %d decisions:")
	text.Recommendation = firstNonEmpty(text.Recommendation, text.RecommendationLabel, "Recommendation")
	text.Rationale = firstNonEmpty(text.Rationale, text.RationaleLabel, "Why")
	text.Evidence = firstNonEmpty(text.Evidence, text.EvidenceLabel, "Evidence")
	text.RoundHeading = text.Heading
	text.RecommendationLabel = text.Recommendation
	text.RationaleLabel = text.Rationale
	text.EvidenceLabel = text.Evidence
	return text
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
		source := firstNonEmpty(question.DefaultSource, DefaultRecommendationSource)
		turns := prompts.Turns()
		if len(turns) > turnCount {
			turnSource := strings.ToLower(strings.TrimSpace(turns[len(turns)-1].Source))
			if turnSource != "" && turnSource != "default" && turnSource != DefaultRecommendationSource {
				source = turnSource
			}
		}
		answers = append(answers, RoundAnswer{
			QuestionID: question.ID,
			Slots:      append([]string(nil), question.Slots...),
			Value:      strings.TrimSpace(value),
			Source:     source,
		})
		if isUserAnswer(source) && strings.EqualFold(strings.TrimSpace(value), "cancel") {
			return answers, nil
		}
	}
	return answers, nil
}

func isUserAnswer(source string) bool {
	source = strings.ToLower(strings.TrimSpace(source))
	return source == "user" || source == "operator" || source == "user_input"
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

func planQuestions[S, D, A any](opts Options[S, D, A], state *S, docs []D, issues []session.ReadinessIssue) ([]Question, error) {
	if opts.planFrontierWithError != nil {
		return opts.planFrontierWithError(state, docs, issues)
	}
	if opts.Interview != nil {
		return opts.Interview.Plan(state, docs)
	}
	if opts.PlanFrontier != nil {
		return opts.PlanFrontier(*state, docs, issues), nil
	}
	if opts.PlanQuestion != nil {
		question := readinesspkg.NormalizeQuestion(opts.PlanQuestion(*state, docs, issues))
		if question.Prompt != "" {
			return []Question{question}, nil
		}
	}
	return nil, nil
}

func applyRound[S, D, A any](opts Options[S, D, A], state *S, answers []RoundAnswer, docs []D, questions []Question) error {
	if err := validateRoundAnswers(answers, questions); err != nil {
		return err
	}
	if opts.Interview != nil {
		if err := opts.Interview.Apply(state, answers, docs); err != nil {
			return err
		}
		if opts.AfterRound != nil {
			return opts.AfterRound(state, answers, docs)
		}
		return nil
	}
	if opts.ApplyRound != nil {
		if err := opts.ApplyRound(state, answers, docs); err != nil {
			return err
		}
		if opts.AfterRound != nil {
			return opts.AfterRound(state, answers, docs)
		}
		return nil
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
		if !ok {
			return fmt.Errorf("answer references unknown frontier question %q", answer.QuestionID)
		}
		if err := opts.ApplyAnswer(state, question, answer.Value, docs); err != nil {
			return err
		}
	}
	if opts.AfterRound != nil {
		return opts.AfterRound(state, answers, docs)
	}
	return nil
}

func validateRoundAnswers(answers []RoundAnswer, questions []Question) error {
	if len(answers) != len(questions) {
		return fmt.Errorf("frontier round must answer all %d questions; got %d", len(questions), len(answers))
	}
	questionIDs := make(map[string]bool, len(questions))
	for _, question := range questions {
		id := strings.TrimSpace(question.ID)
		if questionIDs[id] {
			return fmt.Errorf("frontier contains duplicate question ID %q", id)
		}
		questionIDs[id] = true
	}
	seen := map[string]bool{}
	for _, answer := range answers {
		id := strings.TrimSpace(answer.QuestionID)
		if !questionIDs[id] {
			return fmt.Errorf("answer references unknown frontier question %q", answer.QuestionID)
		}
		if seen[id] {
			return fmt.Errorf("frontier round contains multiple answers for question %q", id)
		}
		seen[id] = true
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
		return fmt.Errorf("icot context is required")
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

func resolveMaxRounds(value int) (int, error) {
	if value < 0 {
		return 0, fmt.Errorf("icot max rounds must be nonnegative")
	}
	if value == 0 {
		return DefaultMaxRounds, nil
	}
	return value, nil
}

func progressFingerprint[S, D any](custom func(S, []D, []session.ReadinessIssue) (string, error), state S, docs []D, issues []session.ReadinessIssue) (string, error) {
	if custom != nil {
		return custom(state, append([]D(nil), docs...), append([]session.ReadinessIssue(nil), issues...))
	}
	data, err := json.Marshal(struct {
		State  S                        `json:"state"`
		Docs   []D                      `json:"docs"`
		Issues []session.ReadinessIssue `json:"issues"`
	}{State: state, Docs: docs, Issues: issues})
	if err != nil {
		return "", err
	}
	return string(data), nil
}
