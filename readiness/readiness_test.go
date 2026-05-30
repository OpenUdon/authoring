package readiness

import (
	"testing"

	"github.com/OpenUdon/authoring/decision"
)

func TestEvaluateBlockingWarningAndTopIssue(t *testing.T) {
	result := Evaluate([]Issue{
		{Code: "warn", Severity: "warning", Slot: "b"},
		{Code: "missing", Slot: "a"},
		{Code: "info", Severity: "info"},
	})
	if result.Ready {
		t.Fatalf("ready = true, want blocking issue to stop readiness")
	}
	if len(result.Blocking) != 1 || result.Blocking[0].Code != "missing" {
		t.Fatalf("blocking = %#v, want missing issue", result.Blocking)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "warn" {
		t.Fatalf("warnings = %#v, want warning issue", result.Warnings)
	}
	if result.TopIssue == nil || result.TopIssue.Code != "missing" {
		t.Fatalf("top issue = %#v, want missing", result.TopIssue)
	}
	if result.Summary.IssueCount != 3 || result.Summary.BlockingCount != 1 || result.Summary.WarningCount != 1 || result.Summary.InfoCount != 1 {
		t.Fatalf("summary = %#v, want stable counters", result.Summary)
	}
}

func TestReadyAllowsWarnings(t *testing.T) {
	if !Ready([]Issue{{Code: "confirm", Severity: "warning"}}) {
		t.Fatalf("warning-only readiness should pass default policy")
	}
	if Ready([]Issue{{Code: "custom", Severity: "needs review"}}) {
		t.Fatalf("unknown severity should block by default")
	}
}

func TestDecisionIssuesRequireConfirmation(t *testing.T) {
	issues := DecisionIssues([]decision.Record{
		{Stage: "draft", Slot: "operation", Value: "create", Source: "model", Confidence: "high"},
		{Stage: "draft", Slot: "operation", Value: "delete", Source: "user", Confidence: "review"},
		{Stage: "draft", Slot: "region", Value: "us-east-1", Source: "model", Confidence: "high"},
	})
	if len(issues) != 2 {
		t.Fatalf("issues = %#v, want conflict issues for operation only", issues)
	}
	if issues[0].Code != DecisionConflictIssueCode || issues[0].Slot != "operation" {
		t.Fatalf("first issue = %#v, want operation conflict", issues[0])
	}
	for _, issue := range issues {
		if !IsBlocking(issue) || issue.SuggestedAnswer == "" {
			t.Fatalf("issue = %#v, want blocking decision confirmation with suggested value", issue)
		}
	}
}

func TestQuestionPlanningForcedAndDefaults(t *testing.T) {
	issue := Issue{Code: "missing_goal", Severity: "blocking", Slot: "goal", SuggestedAnswer: "default goal"}
	suggested := SuggestedQuestion(issue, "Goal?")
	if suggested.ID != "missing_goal" || !suggested.Required || !suggested.AllowDefault {
		t.Fatalf("suggested = %#v, want required defaulted question", suggested)
	}
	answer, source, ok := DefaultAnswer(suggested)
	if !ok || answer != "default goal" || source != DefaultSuggestedAnswerSource {
		t.Fatalf("default = %q %q %v, want suggested answer", answer, source, ok)
	}
	forced := ForcedQuestion(issue, "Goal?")
	if !forced.Forced || !forced.Required || forced.AllowDefault {
		t.Fatalf("forced = %#v, want required operator question without auto default", forced)
	}
	if _, _, ok := DefaultAnswer(forced); ok {
		t.Fatalf("forced question should not auto-apply default")
	}
}

func TestNormalizeQuestionsDeterministic(t *testing.T) {
	plan := EvaluatePlan([]Question{
		{ID: "z", Prompt: "Z"},
		{ID: "a", Prompt: "A", Forced: true, AllowDefault: true, DefaultAnswer: "x"},
		{ID: "b", Prompt: "B", Required: true},
		{ID: "a", Prompt: "A", Forced: true, AllowDefault: true, DefaultAnswer: "x"},
	})
	if len(plan.Questions) != 4 {
		t.Fatalf("questions = %#v", plan.Questions)
	}
	if plan.TopQuestion == nil || plan.TopQuestion.ID != "a" || !plan.TopQuestion.Forced {
		t.Fatalf("top question = %#v, want forced question first", plan.TopQuestion)
	}
	if plan.Summary.ForcedCount != 2 || plan.Summary.RequiredCount != 3 || plan.Summary.DefaultCount != 0 {
		t.Fatalf("summary = %#v, want forced/required/default counters", plan.Summary)
	}
}
