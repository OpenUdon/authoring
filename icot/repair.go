package icot

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	readinesspkg "github.com/OpenUdon/authoring/readiness"
	"github.com/OpenUdon/authoring/transcript"
)

const (
	RepairStatusPassed       = "passed"
	RepairStatusRepaired     = "repaired"
	RepairStatusExhausted    = "exhausted"
	RepairStatusReviewFailed = "review_failed"
	RepairStatusRepairFailed = "repair_failed"
	RepairStatusNoop         = "noop"
)

var (
	// ErrRepairExhausted reports that review still found blocking issues after
	// the bounded repair attempts were used.
	ErrRepairExhausted = errors.New("authoring repair exhausted")
	// ErrRepairNoop reports that the downstream repair hook made no change for
	// blocking review issues.
	ErrRepairNoop = errors.New("authoring repair made no changes")
)

// ReviewIssue combines a generic readiness issue with optional remediation
// metadata. Downstream products own the codes, messages, and actions.
type ReviewIssue struct {
	Issue       readinesspkg.Issue `json:"issue"`
	Remediation *Remediation       `json:"remediation,omitempty"`
}

// Remediation records product-neutral repair metadata.
type Remediation struct {
	Code    string `json:"code,omitempty"`
	Slot    string `json:"slot,omitempty"`
	Action  string `json:"action,omitempty"`
	Message string `json:"message,omitempty"`
	Applied bool   `json:"applied,omitempty"`
}

// RepairOptions supplies hooks for bounded draft review and repair.
type RepairOptions[S, D, A any] struct {
	State       S
	Documents   []D
	Draft       A
	MaxAttempts int

	Normalize func(*S)
	Review    func(context.Context, S, []D, A) ([]readinesspkg.Issue, error)
	Repair    func(context.Context, *S, []D, []readinesspkg.Issue, int) (bool, error)
	Autosave  func(S) error
	OnEvent   func(transcript.Event)
}

// RepairConfig configures a runtime-backed repair pass.
type RepairConfig[S, D, A any] struct {
	State       S
	Documents   []D
	Draft       A
	MaxAttempts int
	Autosave    func(S) error
	OnEvent     func(transcript.Event)
}

// RepairRuntime is the bound runtime shape used by RunRuntimeRepair.
type RepairRuntime[S, D, A any] interface {
	DraftReviewRuntime[S, D, A]
	DraftRepairRuntime[S, D]
}

// RepairResult captures the final repair-loop state.
type RepairResult[S any] struct {
	State     S                   `json:"state"`
	Status    string              `json:"status"`
	Attempts  int                 `json:"attempts,omitempty"`
	Readiness readinesspkg.Result `json:"readiness"`
	Events    []transcript.Event  `json:"events,omitempty"`
	Completed bool                `json:"completed"`
	Repaired  bool                `json:"repaired,omitempty"`
}

// NormalizeReviewIssues returns deterministic review issue containers.
func NormalizeReviewIssues(issues []ReviewIssue) []ReviewIssue {
	out := make([]ReviewIssue, 0, len(issues))
	for _, issue := range issues {
		normalized := readinesspkg.NormalizeIssues([]readinesspkg.Issue{issue.Issue})
		if len(normalized) == 0 {
			continue
		}
		issue.Issue = normalized[0]
		if issue.Remediation != nil {
			remediation := NormalizeRemediation(*issue.Remediation)
			if emptyRemediation(remediation) {
				issue.Remediation = nil
			} else {
				issue.Remediation = &remediation
			}
		}
		out = append(out, issue)
	}
	slices.SortStableFunc(out, CompareReviewIssue)
	return out
}

// NormalizeRemediation returns deterministic product-neutral remediation
// metadata.
func NormalizeRemediation(remediation Remediation) Remediation {
	remediation.Code = normalizeToken(remediation.Code)
	remediation.Slot = strings.TrimSpace(remediation.Slot)
	remediation.Action = normalizeToken(remediation.Action)
	remediation.Message = strings.TrimSpace(remediation.Message)
	return remediation
}

// CompareReviewIssue orders review issue containers deterministically.
func CompareReviewIssue(a, b ReviewIssue) int {
	if diff := readinesspkg.CompareIssue(a.Issue, b.Issue); diff != 0 {
		return diff
	}
	var ar, br Remediation
	if a.Remediation != nil {
		ar = NormalizeRemediation(*a.Remediation)
	}
	if b.Remediation != nil {
		br = NormalizeRemediation(*b.Remediation)
	}
	return compareStrings(ar.Code, br.Code, ar.Slot, br.Slot, ar.Action, br.Action, ar.Message, br.Message)
}

// RunRepair reviews a draft and invokes the downstream repair hook until it is
// ready, attempts are exhausted, or a hook fails.
func RunRepair[S, D, A any](ctx context.Context, opts RepairOptions[S, D, A]) (RepairResult[S], error) {
	var result RepairResult[S]
	if err := checkContext(ctx); err != nil {
		result.Status = RepairStatusRepairFailed
		return result, err
	}
	attemptLimit := opts.MaxAttempts
	if attemptLimit <= 0 {
		attemptLimit = 3
	}
	state := opts.State
	docs := append([]D(nil), opts.Documents...)
	normalize(opts.Normalize, &state)
	result.State = state
	events := eventRecorder(opts.OnEvent)

	for {
		if err := checkContext(ctx); err != nil {
			result = finishRepairResult(state, result, events.Events(), RepairStatusRepairFailed, false)
			return result, err
		}
		issues, err := reviewDraft(ctx, opts.Review, state, docs, opts.Draft)
		if err != nil {
			events.Record("draft_review_error", "review", err.Error(), map[string]string{"attempt": fmt.Sprint(result.Attempts)})
			result = finishRepairResult(state, result, events.Events(), RepairStatusReviewFailed, false)
			return result, err
		}
		ready := readinesspkg.Evaluate(issues)
		result.Readiness = ready
		events.Record("draft_review", "review", "", reviewFields(ready, result.Attempts))
		if ready.Ready {
			status := RepairStatusPassed
			if result.Repaired {
				status = RepairStatusRepaired
			}
			result = finishRepairResult(state, result, events.Events(), status, true)
			return result, nil
		}
		if result.Attempts >= attemptLimit {
			events.Record("draft_repair_exhausted", "repair", "maximum repair attempts reached", reviewFields(ready, result.Attempts))
			result = finishRepairResult(state, result, events.Events(), RepairStatusExhausted, false)
			return result, ErrRepairExhausted
		}
		if opts.Repair == nil {
			events.Record("draft_repair_exhausted", "repair", "repair hook is not configured", reviewFields(ready, result.Attempts))
			result = finishRepairResult(state, result, events.Events(), RepairStatusExhausted, false)
			return result, ErrRepairExhausted
		}
		result.Attempts++
		events.Record("draft_repair_attempt", "repair", "", reviewFields(ready, result.Attempts))
		changed, err := opts.Repair(ctx, &state, docs, ready.Issues, result.Attempts)
		if err != nil {
			events.Record("draft_repair_error", "repair", err.Error(), reviewFields(ready, result.Attempts))
			result = finishRepairResult(state, result, events.Events(), RepairStatusRepairFailed, false)
			return result, err
		}
		if !changed {
			events.Record("draft_repair_noop", "repair", "repair hook made no changes", reviewFields(ready, result.Attempts))
			result = finishRepairResult(state, result, events.Events(), RepairStatusNoop, false)
			return result, ErrRepairNoop
		}
		result.Repaired = true
		normalize(opts.Normalize, &state)
		if opts.Autosave != nil {
			if err := opts.Autosave(state); err != nil {
				events.Record("draft_repair_error", "repair", err.Error(), reviewFields(ready, result.Attempts))
				result = finishRepairResult(state, result, events.Events(), RepairStatusRepairFailed, false)
				return result, err
			}
		}
		events.Record("draft_repair_success", "repair", "", map[string]string{"attempt": fmt.Sprint(result.Attempts)})
		result.State = state
	}
}

// RunRuntimeRepair runs a bounded repair pass through the M08 review/repair
// runtime hook shapes.
func RunRuntimeRepair[S, D, A any](ctx context.Context, runtime RepairRuntime[S, D, A], config RepairConfig[S, D, A]) (RepairResult[S], error) {
	if runtime == nil {
		return RepairResult[S]{}, fmt.Errorf("icot repair runtime is required")
	}
	opts := RepairOptions[S, D, A]{
		State:       config.State,
		Documents:   append([]D(nil), config.Documents...),
		Draft:       config.Draft,
		MaxAttempts: config.MaxAttempts,
		Review:      runtime.ReviewDraft,
		Repair:      runtime.RepairDraft,
		Autosave:    config.Autosave,
		OnEvent:     config.OnEvent,
	}
	if normalizer, ok := runtime.(NormalizingRuntime[S]); ok {
		opts.Normalize = normalizer.Normalize
	}
	return RunRepair(ctx, opts)
}

func reviewDraft[S, D, A any](ctx context.Context, review func(context.Context, S, []D, A) ([]readinesspkg.Issue, error), state S, docs []D, draft A) ([]readinesspkg.Issue, error) {
	if review == nil {
		return nil, nil
	}
	issues, err := review(ctx, state, docs, draft)
	if err != nil {
		return nil, err
	}
	return readinesspkg.NormalizeIssues(issues), nil
}

func finishRepairResult[S any](state S, result RepairResult[S], events []transcript.Event, status string, completed bool) RepairResult[S] {
	result.State = state
	result.Status = status
	result.Completed = completed
	result.Events = append([]transcript.Event(nil), events...)
	return result
}

func reviewFields(result readinesspkg.Result, attempt int) map[string]string {
	fields := map[string]string{
		"attempt":  fmt.Sprint(attempt),
		"issues":   fmt.Sprint(result.Summary.IssueCount),
		"blocking": fmt.Sprint(result.Summary.BlockingCount),
		"warnings": fmt.Sprint(result.Summary.WarningCount),
	}
	if result.TopIssue != nil {
		if result.TopIssue.Code != "" {
			fields["top_code"] = result.TopIssue.Code
		}
		if result.TopIssue.Slot != "" {
			fields["top_slot"] = result.TopIssue.Slot
		}
	}
	return fields
}

type eventLog struct {
	events  []transcript.Event
	onEvent func(transcript.Event)
}

func eventRecorder(onEvent func(transcript.Event)) *eventLog {
	return &eventLog{onEvent: onEvent}
}

func (log *eventLog) Record(eventType, stage, message string, fields map[string]string) {
	event := transcript.Event{
		Type:    eventType,
		Stage:   stage,
		Message: strings.TrimSpace(message),
		Fields:  fields,
	}
	log.events = append(log.events, event)
	if log.onEvent != nil {
		log.onEvent(event)
	}
}

func (log *eventLog) Events() []transcript.Event {
	return append([]transcript.Event(nil), log.events...)
}

func emptyRemediation(remediation Remediation) bool {
	return remediation.Code == "" && remediation.Slot == "" && remediation.Action == "" && remediation.Message == ""
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), "_"))
}

func compareStrings(values ...string) int {
	for i := 0; i+1 < len(values); i += 2 {
		if values[i] < values[i+1] {
			return -1
		}
		if values[i] > values[i+1] {
			return 1
		}
	}
	return 0
}
