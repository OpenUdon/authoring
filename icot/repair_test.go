package icot

import (
	"context"
	"errors"
	"strings"
	"testing"

	readinesspkg "github.com/OpenUdon/authoring/readiness"
)

type repairState struct {
	Fixed      bool
	Normalized bool
}

func TestRunRepairSuccess(t *testing.T) {
	var autosaved bool
	result, err := RunRepair(context.Background(), RepairOptions[repairState, string, string]{
		State:       repairState{},
		Documents:   []string{"doc"},
		Draft:       "draft",
		MaxAttempts: 2,
		Normalize: func(state *repairState) {
			state.Normalized = true
		},
		Review: func(_ context.Context, state repairState, _ []string, _ string) ([]readinesspkg.Issue, error) {
			if state.Fixed {
				return nil, nil
			}
			return []readinesspkg.Issue{{Code: "missing", Severity: "blocking"}}, nil
		},
		Repair: func(_ context.Context, state *repairState, _ []string, issues []readinesspkg.Issue, attempt int) (bool, error) {
			if len(issues) != 1 || attempt != 1 {
				t.Fatalf("repair issues=%#v attempt=%d, want first blocking issue", issues, attempt)
			}
			state.Fixed = true
			return true, nil
		},
		Autosave: func(repairState) error {
			autosaved = true
			return nil
		},
	})
	if err != nil {
		t.Fatalf("RunRepair error = %v", err)
	}
	if result.Status != RepairStatusRepaired || !result.Completed || !result.Repaired || result.Attempts != 1 {
		t.Fatalf("result = %#v, want repaired completion", result)
	}
	if !result.State.Fixed || !result.State.Normalized || !autosaved {
		t.Fatalf("state=%#v autosaved=%v, want repaired normalized autosaved state", result.State, autosaved)
	}
	if eventTypes(result.Events) != "draft_review,draft_repair_attempt,draft_repair_success,draft_review" {
		t.Fatalf("events = %s", eventTypes(result.Events))
	}
}

func TestRunRepairExhaustion(t *testing.T) {
	result, err := RunRepair(context.Background(), RepairOptions[repairState, string, string]{
		MaxAttempts: 2,
		Review: func(context.Context, repairState, []string, string) ([]readinesspkg.Issue, error) {
			return []readinesspkg.Issue{{Code: "still_missing", Severity: "blocking"}}, nil
		},
		Repair: func(context.Context, *repairState, []string, []readinesspkg.Issue, int) (bool, error) {
			return true, nil
		},
	})
	if !errors.Is(err, ErrRepairExhausted) {
		t.Fatalf("err = %v, want ErrRepairExhausted", err)
	}
	if result.Status != RepairStatusExhausted || result.Completed || result.Attempts != 2 {
		t.Fatalf("result = %#v, want exhausted after two attempts", result)
	}
	if got := eventTypes(result.Events); !strings.HasSuffix(got, "draft_repair_exhausted") {
		t.Fatalf("events = %s, want exhaustion event", got)
	}
}

func TestRunRepairReviewError(t *testing.T) {
	reviewErr := errors.New("review failed")
	result, err := RunRepair(context.Background(), RepairOptions[repairState, string, string]{
		Review: func(context.Context, repairState, []string, string) ([]readinesspkg.Issue, error) {
			return nil, reviewErr
		},
	})
	if !errors.Is(err, reviewErr) {
		t.Fatalf("err = %v, want review error", err)
	}
	if result.Status != RepairStatusReviewFailed || eventTypes(result.Events) != "draft_review_error" {
		t.Fatalf("result = %#v events=%s, want review failure", result, eventTypes(result.Events))
	}
}

func TestRunRepairNoop(t *testing.T) {
	result, err := RunRepair(context.Background(), RepairOptions[repairState, string, string]{
		Review: func(context.Context, repairState, []string, string) ([]readinesspkg.Issue, error) {
			return []readinesspkg.Issue{{Code: "missing", Severity: "blocking"}}, nil
		},
		Repair: func(context.Context, *repairState, []string, []readinesspkg.Issue, int) (bool, error) {
			return false, nil
		},
	})
	if !errors.Is(err, ErrRepairNoop) {
		t.Fatalf("err = %v, want ErrRepairNoop", err)
	}
	if result.Status != RepairStatusNoop || result.Completed {
		t.Fatalf("result = %#v, want noop failure", result)
	}
	if got := eventTypes(result.Events); got != "draft_review,draft_repair_attempt,draft_repair_noop" {
		t.Fatalf("events = %s", got)
	}
}

func TestRunRuntimeRepairBindsHooks(t *testing.T) {
	runtime := &repairRuntime{}
	result, err := RunRuntimeRepair(context.Background(), runtime, RepairConfig[repairState, string, string]{
		MaxAttempts: 1,
	})
	if err != nil {
		t.Fatalf("RunRuntimeRepair error = %v", err)
	}
	if result.Status != RepairStatusRepaired || !result.State.Fixed || !result.State.Normalized {
		t.Fatalf("result = %#v, want runtime-bound repaired state", result)
	}
}

func TestNormalizeReviewIssues(t *testing.T) {
	issues := NormalizeReviewIssues([]ReviewIssue{
		{
			Issue: readinesspkg.Issue{Code: "warn", Severity: "warning"},
			Remediation: &Remediation{
				Code:    " Fix It ",
				Action:  " Ask User ",
				Message: " apply patch ",
			},
		},
		{Issue: readinesspkg.Issue{Code: "missing", Severity: "blocking", Slot: "goal"}},
		{},
	})
	if len(issues) != 2 {
		t.Fatalf("issues = %#v, want empty issue removed", issues)
	}
	if issues[0].Issue.Code != "missing" || issues[1].Remediation.Code != "fix_it" || issues[1].Remediation.Action != "ask_user" {
		t.Fatalf("issues = %#v, want sorted normalized issues", issues)
	}
}

type repairRuntime struct{}

func (repairRuntime) ReviewDraft(_ context.Context, state repairState, _ []string, _ string) ([]readinesspkg.Issue, error) {
	if state.Fixed {
		return nil, nil
	}
	return []readinesspkg.Issue{{Code: "missing", Severity: "blocking"}}, nil
}

func (repairRuntime) RepairDraft(_ context.Context, state *repairState, _ []string, _ []readinesspkg.Issue, _ int) (bool, error) {
	state.Fixed = true
	return true, nil
}

func (repairRuntime) Normalize(state *repairState) {
	state.Normalized = true
}
