package icot

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
	Session     S
	Documents   []D
	MaxAttempts int
	DefaultMode prompt.DefaultMode
	Autosave    func(S) error
	AfterDraft  func(S) error
}

// BindRuntime converts a bound runtime into loop options.
func BindRuntime[S, D, A any](runtime Runtime[S, D, A], config RuntimeConfig[S, D]) (Options[S, D, A], error) {
	if runtime == nil {
		return Options[S, D, A]{}, fmt.Errorf("icot runtime is required")
	}
	opts := Options[S, D, A]{
		Session:        config.Session,
		Documents:      append([]D(nil), config.Documents...),
		MaxAttempts:    config.MaxAttempts,
		DefaultMode:    config.DefaultMode,
		Draft:          runtime.Draft,
		CheckReadiness: runtime.Readiness,
		PlanQuestion:   runtime.PlanQuestion,
		ApplyAnswer:    runtime.ApplyAnswer,
		Autosave:       config.Autosave,
		AfterDraft:     config.AfterDraft,
		FinalConfirm:   runtime.WriteArtifacts,
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
func RunRuntime[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, runtime Runtime[S, D, A], config RuntimeConfig[S, D]) (Result[S, A], error) {
	opts, err := BindRuntime[S, D, A](runtime, config)
	if err != nil {
		return Result[S, A]{}, err
	}
	return Run(ctx, in, out, opts)
}
