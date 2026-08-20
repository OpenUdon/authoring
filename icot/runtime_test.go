package icot

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/interview"
	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

type boundRuntime struct {
	draftErr      error
	refreshErr    error
	writeErr      error
	refreshCalled bool
	normalized    bool
}

func (r *boundRuntime) Normalize(state *fakeState) {
	r.normalized = true
	state.Goal = strings.TrimSpace(state.Goal)
}

func (r *boundRuntime) Draft(context.Context, fakeState, []string, []session.ReadinessIssue, int) (fakeState, error) {
	if r.draftErr != nil {
		return fakeState{}, r.draftErr
	}
	return fakeState{Goal: "drafted", Ready: true}, nil
}

func (r *boundRuntime) Readiness(state fakeState, _ []string) []session.ReadinessIssue {
	if state.Ready {
		return nil
	}
	return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
}

func (r *boundRuntime) Ready(_ fakeState, issues []session.ReadinessIssue) bool {
	return len(issues) == 0
}

func (r *boundRuntime) PlanQuestion(fakeState, []string, []session.ReadinessIssue) Question {
	return Question{ID: "goal", Prompt: "Goal", AllowDefault: true, DefaultAnswer: "defaulted"}
}

func (r *boundRuntime) ApplyAnswer(state *fakeState, _ Question, answer string, _ []string) error {
	state.Goal = answer
	state.Ready = true
	return nil
}

func (r *boundRuntime) RefreshDocuments(context.Context, fakeState, []string) ([]string, error) {
	r.refreshCalled = true
	if r.refreshErr != nil {
		return nil, r.refreshErr
	}
	return []string{"refreshed"}, nil
}

func (r *boundRuntime) ShouldDraft(fakeState, []string, []session.ReadinessIssue, int) bool {
	return true
}

func (r *boundRuntime) WriteArtifacts(_ context.Context, state *fakeState, _ []string, _ *[]transcript.Event) (string, error) {
	if r.writeErr != nil {
		return "", r.writeErr
	}
	return state.Goal, nil
}

func TestRunRuntimeSuccess(t *testing.T) {
	runtime := &boundRuntime{}
	result, err := RunRuntime[fakeState, string, string](context.Background(), nil, nil, runtime, RuntimeConfig[fakeState, string]{
		Session:     fakeState{Goal: " seed "},
		Documents:   []string{"doc"},
		MaxAttempts: 2,
	})
	if err != nil {
		t.Fatalf("RunRuntime returned error: %v", err)
	}
	if !result.Completed || result.Artifact != "drafted" {
		t.Fatalf("result = %#v, want drafted artifact", result)
	}
	if !runtime.normalized {
		t.Fatalf("runtime Normalize was not called")
	}
	if !runtime.refreshCalled {
		t.Fatalf("runtime RefreshDocuments was not called")
	}
}

func TestBindRuntimeRequiresRuntime(t *testing.T) {
	if _, err := BindRuntime[fakeState, string, string](nil, RuntimeConfig[fakeState, string]{}); err == nil {
		t.Fatalf("BindRuntime accepted nil runtime")
	}
	if _, err := RunRuntime[fakeState, string, string](context.Background(), nil, nil, nil, RuntimeConfig[fakeState, string]{}); err == nil {
		t.Fatalf("RunRuntime accepted nil runtime")
	}
	if _, err := RunRuntime[fakeState, string, string](nil, nil, nil, &boundRuntime{}, RuntimeConfig[fakeState, string]{}); err == nil || !strings.Contains(err.Error(), "context is required") {
		t.Fatalf("RunRuntime nil context error = %v", err)
	}
}

func TestRunRuntimePropagatesHookErrors(t *testing.T) {
	writeErr := errors.New("write failed")
	_, err := RunRuntime[fakeState, string, string](context.Background(), nil, nil, &boundRuntime{writeErr: writeErr}, RuntimeConfig[fakeState, string]{MaxAttempts: 1})
	if !errors.Is(err, writeErr) {
		t.Fatalf("RunRuntime error = %v, want write error", err)
	}

	refreshErr := errors.New("refresh failed")
	_, err = RunRuntime[fakeState, string, string](context.Background(), nil, nil, &boundRuntime{refreshErr: refreshErr}, RuntimeConfig[fakeState, string]{MaxAttempts: 1})
	if !errors.Is(err, refreshErr) {
		t.Fatalf("RunRuntime error = %v, want refresh error", err)
	}
}

func TestRunRuntimeDraftErrorFallsBackToQuestion(t *testing.T) {
	result, err := RunRuntime[fakeState, string, string](context.Background(), nil, nil, &boundRuntime{draftErr: errors.New("draft failed")}, RuntimeConfig[fakeState, string]{DefaultMode: prompt.DefaultsSilent})
	if err != nil {
		t.Fatalf("RunRuntime returned error: %v", err)
	}
	if result.Artifact != "defaulted" {
		t.Fatalf("artifact = %q, want defaulted fallback", result.Artifact)
	}
}

type interviewBoundRuntime struct{}

func (interviewBoundRuntime) Draft(_ context.Context, state boundInterviewState, _ []string, _ []session.ReadinessIssue, _ int) (boundInterviewState, error) {
	return state, nil
}

func (interviewBoundRuntime) Readiness(state boundInterviewState, _ []string) []session.ReadinessIssue {
	if state.Values["a"] != "" {
		return nil
	}
	return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
}

func (interviewBoundRuntime) InterviewBinding() InterviewBinding[boundInterviewState, string] {
	return testInterviewBinding(nil)
}

func (interviewBoundRuntime) WriteArtifacts(_ context.Context, state *boundInterviewState, _ []string, _ *[]transcript.Event) (string, error) {
	return state.Values["a"], nil
}

func TestRunRuntimeUsesInterviewBindingAlternative(t *testing.T) {
	result, err := RunRuntime[boundInterviewState, string, string](context.Background(), nil, nil, interviewBoundRuntime{}, RuntimeConfig[boundInterviewState, string]{
		Session:     boundInterviewState{Interview: interview.State{Nodes: []interview.Node{{ID: "a", Prompt: "A", Recommendation: "bound"}}}, Values: map[string]string{}},
		DefaultMode: prompt.DefaultsSilent,
	})
	if err != nil || !result.Completed || result.Artifact != "bound" || result.Session.Interview.Round != 1 {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}
