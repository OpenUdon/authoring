// Package icot preserves the existing authoring API over the neutral engine.
// Generic orchestration has one implementation in package engine.
package icot

import (
	"context"
	"io"

	"github.com/OpenUdon/authoring/engine"
	"github.com/OpenUdon/authoring/readiness"
	"github.com/OpenUdon/authoring/transcript"
)

// ErrNeedsInput retains the existing engine value.
var ErrNeedsInput = engine.ErrNeedsInput

// ErrCanceled retains the existing engine value.
var ErrCanceled = engine.ErrCanceled

// ErrNoProgress retains the existing engine value.
var ErrNoProgress = engine.ErrNoProgress

// ErrRoundLimit retains the existing engine value.
var ErrRoundLimit = engine.ErrRoundLimit

// DefaultNoProgressLimit retains the existing engine value.
const DefaultNoProgressLimit = engine.DefaultNoProgressLimit

// DefaultMaxRounds retains the existing engine value.
const DefaultMaxRounds = engine.DefaultMaxRounds

// RoundAnswer is the source-compatible engine alias.
type RoundAnswer = engine.RoundAnswer

// FrontierText is the source-compatible engine alias.
type FrontierText = engine.FrontierText

// Options is the source-compatible engine alias.
type Options[S, D, A any] = engine.Options[S, D, A]

// Question is the source-compatible engine alias.
type Question = engine.Question

// Result is the source-compatible engine alias.
type Result[S, A any] = engine.Result[S, A]

// PlanFrontier delegates to the single engine implementation.
func PlanFrontier(questions []Question) readiness.Plan { return engine.PlanFrontier(questions) }

// Run delegates to the single engine implementation.
func Run[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, opts Options[S, D, A]) (Result[S, A], error) {
	return engine.Run[S, D, A](ctx, in, out, opts)
}

// DefaultRecommendationSource retains the existing engine value.
const DefaultRecommendationSource = engine.DefaultRecommendationSource

// PromptTurn is the source-compatible engine alias.
type PromptTurn = engine.PromptTurn

// Event is the source-compatible engine alias.
type Event = engine.Event

// PromptTranscript is the source-compatible engine alias.
type PromptTranscript = engine.PromptTranscript

// PromptSession is the source-compatible engine alias.
type PromptSession = engine.PromptSession

// NewPromptSession delegates to the single engine implementation.
func NewPromptSession(in io.Reader, out io.Writer) *PromptSession {
	return engine.NewPromptSession(in, out)
}

// OneLine delegates to the single engine implementation.
func OneLine(value string) string { return engine.OneLine(value) }

// AssertPromptLabelsInOrder delegates to the single engine implementation.
func AssertPromptLabelsInOrder(output string, turns []PromptTurn) error {
	return engine.AssertPromptLabelsInOrder(output, turns)
}

// SavePromptTranscript delegates to the single engine implementation.
func SavePromptTranscript(path, version string, turns []PromptTurn, events []Event, session any) error {
	return engine.SavePromptTranscript(path, version, turns, events, session)
}

// NormalizeEvents delegates to the single engine implementation.
func NormalizeEvents(events []Event) []Event { return engine.NormalizeEvents(events) }

// TranscriptEvent delegates to the single engine implementation.
func TranscriptEvent(event Event) transcript.Event { return engine.TranscriptEvent(event) }

// TranscriptEvents delegates to the single engine implementation.
func TranscriptEvents(events []Event) []transcript.Event { return engine.TranscriptEvents(events) }

// ReadinessIssue is the source-compatible engine alias.
type ReadinessIssue = engine.ReadinessIssue

// InteractiveQuestion is the source-compatible engine alias.
type InteractiveQuestion = engine.InteractiveQuestion

// DraftRequest is the source-compatible engine alias.
type DraftRequest[S, D any] = engine.DraftRequest[S, D]

// Extractor is the source-compatible engine alias.
type Extractor[S, D any] = engine.Extractor[S, D]

// NoopExtractor is the source-compatible engine alias.
type NoopExtractor[S, D any] = engine.NoopExtractor[S, D]

// InteractiveHooks is the source-compatible engine alias.
type InteractiveHooks[S, D, A any] = engine.InteractiveHooks[S, D, A]

// InteractiveLifecycleOptions is the source-compatible engine alias.
type InteractiveLifecycleOptions[S, D, A any] = engine.InteractiveLifecycleOptions[S, D, A]

// RunInteractiveWithLifecycle delegates to the single engine implementation.
func RunInteractiveWithLifecycle[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, hooks InteractiveHooks[S, D, A], opts InteractiveLifecycleOptions[S, D, A]) (A, error) {
	return engine.RunInteractiveWithLifecycle[S, D, A](ctx, in, out, hooks, opts)
}

// RunInteractive delegates to the single engine implementation.
func RunInteractive[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, hooks InteractiveHooks[S, D, A]) (A, error) {
	return engine.RunInteractive[S, D, A](ctx, in, out, hooks)
}

// InterviewBinding is the source-compatible engine alias.
type InterviewBinding[S, D any] = engine.InterviewBinding[S, D]

// RepairStatusPassed retains the existing engine value.
const RepairStatusPassed = engine.RepairStatusPassed

// RepairStatusRepaired retains the existing engine value.
const RepairStatusRepaired = engine.RepairStatusRepaired

// RepairStatusExhausted retains the existing engine value.
const RepairStatusExhausted = engine.RepairStatusExhausted

// RepairStatusReviewFailed retains the existing engine value.
const RepairStatusReviewFailed = engine.RepairStatusReviewFailed

// RepairStatusRepairFailed retains the existing engine value.
const RepairStatusRepairFailed = engine.RepairStatusRepairFailed

// RepairStatusNoop retains the existing engine value.
const RepairStatusNoop = engine.RepairStatusNoop

// ErrRepairExhausted retains the existing engine value.
var ErrRepairExhausted = engine.ErrRepairExhausted

// ErrRepairNoop retains the existing engine value.
var ErrRepairNoop = engine.ErrRepairNoop

// ReviewIssue is the source-compatible engine alias.
type ReviewIssue = engine.ReviewIssue

// Remediation is the source-compatible engine alias.
type Remediation = engine.Remediation

// RepairOptions is the source-compatible engine alias.
type RepairOptions[S, D, A any] = engine.RepairOptions[S, D, A]

// RepairConfig is the source-compatible engine alias.
type RepairConfig[S, D, A any] = engine.RepairConfig[S, D, A]

// RepairRuntime is the source-compatible engine alias.
type RepairRuntime[S, D, A any] = engine.RepairRuntime[S, D, A]

// RepairResult is the source-compatible engine alias.
type RepairResult[S any] = engine.RepairResult[S]

// NormalizeReviewIssues delegates to the single engine implementation.
func NormalizeReviewIssues(issues []ReviewIssue) []ReviewIssue {
	return engine.NormalizeReviewIssues(issues)
}

// NormalizeRemediation delegates to the single engine implementation.
func NormalizeRemediation(remediation Remediation) Remediation {
	return engine.NormalizeRemediation(remediation)
}

// CompareReviewIssue delegates to the single engine implementation.
func CompareReviewIssue(a, b ReviewIssue) int { return engine.CompareReviewIssue(a, b) }

// RunRepair delegates to the single engine implementation.
func RunRepair[S, D, A any](ctx context.Context, opts RepairOptions[S, D, A]) (RepairResult[S], error) {
	return engine.RunRepair[S, D, A](ctx, opts)
}

// RunRuntimeRepair delegates to the single engine implementation.
func RunRuntimeRepair[S, D, A any](ctx context.Context, runtime RepairRuntime[S, D, A], config RepairConfig[S, D, A]) (RepairResult[S], error) {
	return engine.RunRuntimeRepair[S, D, A](ctx, runtime, config)
}

// Runtime is the source-compatible engine alias.
type Runtime[S, D, A any] = engine.Runtime[S, D, A]

// InterviewRuntime is the source-compatible engine alias.
type InterviewRuntime[S, D, A any] = engine.InterviewRuntime[S, D, A]

// LegacyRuntime is the source-compatible engine alias.
type LegacyRuntime[S, D, A any] = engine.LegacyRuntime[S, D, A]

// NormalizingRuntime is the source-compatible engine alias.
type NormalizingRuntime[S any] = engine.NormalizingRuntime[S]

// ReadyRuntime is the source-compatible engine alias.
type ReadyRuntime[S any] = engine.ReadyRuntime[S]

// DocumentRefreshingRuntime is the source-compatible engine alias.
type DocumentRefreshingRuntime[S, D any] = engine.DocumentRefreshingRuntime[S, D]

// DraftPolicyRuntime is the source-compatible engine alias.
type DraftPolicyRuntime[S, D any] = engine.DraftPolicyRuntime[S, D]

// DraftReviewRuntime is the source-compatible engine alias.
type DraftReviewRuntime[S, D, A any] = engine.DraftReviewRuntime[S, D, A]

// DraftRepairRuntime is the source-compatible engine alias.
type DraftRepairRuntime[S, D any] = engine.DraftRepairRuntime[S, D]

// RuntimeConfig is the source-compatible engine alias.
type RuntimeConfig[S, D any] = engine.RuntimeConfig[S, D]

// BindRuntime delegates to the single engine implementation.
func BindRuntime[S, D, A any](runtime any, config RuntimeConfig[S, D]) (Options[S, D, A], error) {
	return engine.BindRuntime[S, D, A](runtime, config)
}

// RunRuntime delegates to the single engine implementation.
func RunRuntime[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, runtime any, config RuntimeConfig[S, D]) (Result[S, A], error) {
	return engine.RunRuntime[S, D, A](ctx, in, out, runtime, config)
}
