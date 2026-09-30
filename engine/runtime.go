package engine

import (
	"context"
	"fmt"
	"io"

	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

// Runtime is the bound product adapter used by the generic iCoT loop.
//
// Implementations own domain semantics: draft schemas, readiness issue codes,
// question planning, answer application, and artifact generation.
type Runtime[S, D, A any] interface {
	Draft(context.Context, S, []D, []session.ReadinessIssue, int) (S, error)
	Readiness(S, []D) []session.ReadinessIssue
	PlanFrontier(S, []D, []session.ReadinessIssue) []Question
	ApplyRound(*S, []RoundAnswer, []D) error
	WriteArtifacts(context.Context, *S, []D, *[]transcript.Event) (A, error)
}

// InterviewRuntime is the bound product adapter shape that delegates frontier
// projection and atomic settlement to an InterviewBinding.
type InterviewRuntime[S, D, A any] interface {
	Draft(context.Context, S, []D, []session.ReadinessIssue, int) (S, error)
	Readiness(S, []D) []session.ReadinessIssue
	InterviewBinding() InterviewBinding[S, D]
	WriteArtifacts(context.Context, *S, []D, *[]transcript.Event) (A, error)
}

// LegacyRuntime is the pre-frontier source compatibility shape. BindRuntime
// adapts it to one-question rounds; new runtimes should implement Runtime.
type LegacyRuntime[S, D, A any] interface {
	Draft(context.Context, S, []D, []session.ReadinessIssue, int) (S, error)
	Readiness(S, []D) []session.ReadinessIssue
	PlanQuestion(S, []D, []session.ReadinessIssue) Question
	ApplyAnswer(*S, Question, string, []D) error
	WriteArtifacts(context.Context, *S, []D, *[]transcript.Event) (A, error)
}

// NormalizingRuntime can normalize product state after mutations.
type NormalizingRuntime[S any] interface {
	Normalize(*S)
}

// ReadyRuntime can override the default readiness policy.
type ReadyRuntime[S any] interface {
	Ready(S, []session.ReadinessIssue) bool
}

// DocumentRefreshingRuntime can refresh product documents between attempts.
type DocumentRefreshingRuntime[S, D any] interface {
	RefreshDocuments(context.Context, S, []D) ([]D, error)
}

// DraftPolicyRuntime can suppress or allow draft attempts.
type DraftPolicyRuntime[S, D any] interface {
	ShouldDraft(S, []D, []session.ReadinessIssue, int) bool
}

// DraftReviewRuntime defines the generic draft review hook shape.
type DraftReviewRuntime[S, D, A any] interface {
	ReviewDraft(context.Context, S, []D, A) ([]session.ReadinessIssue, error)
}

// DraftRepairRuntime defines the generic draft repair hook shape.
type DraftRepairRuntime[S, D any] interface {
	RepairDraft(context.Context, *S, []D, []session.ReadinessIssue, int) (bool, error)
}

// RuntimeConfig configures a runtime-backed loop.
type RuntimeConfig[S, D any] struct {
	Session             S
	Documents           []D
	DefaultMode         prompt.DefaultMode
	FrontierText        FrontierText
	Autosave            func(S) error
	AfterDraft          func(S) error
	ProgressFingerprint func(S, []D, []session.ReadinessIssue) (string, error)
	OnEvent             func(transcript.Event)
	MaxRounds           int

	// Deprecated: retained for source compatibility and ignored.
	MaxAttempts int
	// Deprecated: retained for source compatibility and ignored.
	NoProgressLimit int
}

// BindRuntime converts a bound runtime into loop options.
func BindRuntime[S, D, A any](runtime any, config RuntimeConfig[S, D]) (Options[S, D, A], error) {
	if runtime == nil {
		return Options[S, D, A]{}, fmt.Errorf("icot runtime is required")
	}
	opts := Options[S, D, A]{
		Session:             config.Session,
		Documents:           append([]D(nil), config.Documents...),
		MaxAttempts:         config.MaxAttempts,
		MaxRounds:           config.MaxRounds,
		DefaultMode:         config.DefaultMode,
		FrontierText:        config.FrontierText,
		Autosave:            config.Autosave,
		AfterDraft:          config.AfterDraft,
		ProgressFingerprint: config.ProgressFingerprint,
		OnEvent:             config.OnEvent,
	}
	switch bound := runtime.(type) {
	case InterviewRuntime[S, D, A]:
		binding := bound.InterviewBinding()
		opts.Interview = &binding
		opts.Draft = bound.Draft
		opts.CheckReadiness = bound.Readiness
		opts.FinalConfirm = bound.WriteArtifacts
	case Runtime[S, D, A]:
		opts.Draft = bound.Draft
		opts.CheckReadiness = bound.Readiness
		opts.PlanFrontier = bound.PlanFrontier
		opts.ApplyRound = bound.ApplyRound
		opts.FinalConfirm = bound.WriteArtifacts
	case LegacyRuntime[S, D, A]:
		opts.Draft = bound.Draft
		opts.CheckReadiness = bound.Readiness
		opts.PlanQuestion = bound.PlanQuestion
		opts.ApplyAnswer = bound.ApplyAnswer
		opts.FinalConfirm = bound.WriteArtifacts
	default:
		return Options[S, D, A]{}, fmt.Errorf("icot runtime does not implement the frontier runtime contract")
	}
	if normalizer, ok := runtime.(NormalizingRuntime[S]); ok {
		opts.Normalize = normalizer.Normalize
	}
	if ready, ok := runtime.(ReadyRuntime[S]); ok {
		opts.Ready = ready.Ready
	}
	if refresher, ok := runtime.(DocumentRefreshingRuntime[S, D]); ok {
		opts.RefreshDocs = refresher.RefreshDocuments
	}
	if policy, ok := runtime.(DraftPolicyRuntime[S, D]); ok {
		opts.ShouldDraft = policy.ShouldDraft
	}
	return opts, nil
}

// RunRuntime runs the generic loop through a bound runtime.
func RunRuntime[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, runtime any, config RuntimeConfig[S, D]) (Result[S, A], error) {
	if err := checkContext(ctx); err != nil {
		return Result[S, A]{}, err
	}
	opts, err := BindRuntime[S, D, A](runtime, config)
	if err != nil {
		return Result[S, A]{}, err
	}
	return Run(ctx, in, out, opts)
}
