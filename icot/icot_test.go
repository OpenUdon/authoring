package icot

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

type fakeState struct {
	Goal  string
	Ready bool
}

func TestRunProgressiveSuccess(t *testing.T) {
	var order []string
	result, err := Run[fakeState, string, string](context.Background(), strings.NewReader(""), nil, Options[fakeState, string, string]{
		Session: fakeState{Goal: "draft"},
		Draft: func(context.Context, fakeState, []string, []session.ReadinessIssue, int) (fakeState, error) {
			order = append(order, "draft")
			return fakeState{Goal: "drafted", Ready: true}, nil
		},
		CheckReadiness: func(state fakeState, _ []string) []session.ReadinessIssue {
			order = append(order, "readiness")
			if state.Ready {
				return nil
			}
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		Ready: func(_ fakeState, issues []session.ReadinessIssue) bool {
			return len(issues) == 0
		},
		FinalConfirm: func(_ context.Context, state *fakeState, _ []string, events *[]transcript.Event) (string, error) {
			order = append(order, "confirm")
			return state.Goal, nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !result.Completed || result.Artifact != "drafted" {
		t.Fatalf("result = %#v, want completed drafted result", result)
	}
	if strings.Join(order, ",") != "readiness,draft,readiness,confirm" {
		t.Fatalf("order = %v", order)
	}
	if !strings.Contains(eventTypes(result.Events), "draft_attempt,draft_success,readiness") {
		t.Fatalf("events = %#v", result.Events)
	}
}

func TestRunNeedsInputWhenRequiredQuestionHasNoInput(t *testing.T) {
	result, err := Run[fakeState, string, string](context.Background(), strings.NewReader(""), nil, Options[fakeState, string, string]{
		CheckReadiness: func(fakeState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(fakeState, []string, []session.ReadinessIssue) Question {
			return Question{ID: "goal", Prompt: "Goal", Required: true}
		},
		MaxAttempts: 1,
	})
	if !errors.Is(err, ErrNeedsInput) {
		t.Fatalf("Run error = %v, want ErrNeedsInput", err)
	}
	if result.Completed {
		t.Fatalf("result completed unexpectedly: %#v", result)
	}
}

func TestRunCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Run[fakeState, string, string](ctx, nil, nil, Options[fakeState, string, string]{})
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("Run error = %v, want ErrCanceled", err)
	}
}

func TestRunDraftErrorThenDefaultedAnswer(t *testing.T) {
	var draftErrors int
	result, err := Run[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{
		DefaultMode: prompt.DefaultsSilent,
		MaxAttempts: 2,
		Draft: func(context.Context, fakeState, []string, []session.ReadinessIssue, int) (fakeState, error) {
			return fakeState{}, errors.New("draft failed")
		},
		OnDraftError: func(error) {
			draftErrors++
		},
		CheckReadiness: func(state fakeState, _ []string) []session.ReadinessIssue {
			if state.Ready {
				return nil
			}
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(fakeState, []string, []session.ReadinessIssue) Question {
			return Question{ID: "goal", Prompt: "Goal", AllowDefault: true, DefaultAnswer: "defaulted", DefaultSource: "test"}
		},
		ApplyAnswer: func(state *fakeState, _ Question, answer string, _ []string) error {
			state.Goal = answer
			state.Ready = true
			return nil
		},
		FinalConfirm: func(_ context.Context, state *fakeState, _ []string, _ *[]transcript.Event) (string, error) {
			return state.Goal, nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if draftErrors != 1 || result.Artifact != "defaulted" {
		t.Fatalf("draftErrors=%d result=%#v, want fallback defaulted result", draftErrors, result)
	}
	if !strings.Contains(eventTypes(result.Events), "draft_error") || !strings.Contains(eventTypes(result.Events), "question_answered") {
		t.Fatalf("events = %#v", result.Events)
	}
}

func TestRunForcedQuestionIgnoresDefault(t *testing.T) {
	result, err := Run[fakeState, string, string](context.Background(), strings.NewReader("manual\n"), nil, Options[fakeState, string, string]{
		MaxAttempts: 2,
		CheckReadiness: func(state fakeState, _ []string) []session.ReadinessIssue {
			if state.Ready {
				return nil
			}
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(fakeState, []string, []session.ReadinessIssue) Question {
			return Question{ID: "goal", Prompt: "Goal", Forced: true, AllowDefault: true, DefaultAnswer: "defaulted"}
		},
		ApplyAnswer: func(state *fakeState, _ Question, answer string, _ []string) error {
			state.Goal = answer
			state.Ready = true
			return nil
		},
		FinalConfirm: func(_ context.Context, state *fakeState, _ []string, _ *[]transcript.Event) (string, error) {
			return state.Goal, nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.Artifact != "manual" || result.Turns[0].Answer != "manual" {
		t.Fatalf("result = %#v turns=%#v, want manual forced answer", result, result.Turns)
	}
}

func TestRunStopsAtMaxAttempts(t *testing.T) {
	var answers int
	_, err := Run[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{
		MaxAttempts: 2,
		CheckReadiness: func(fakeState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(fakeState, []string, []session.ReadinessIssue) Question {
			return Question{ID: "goal", Prompt: "Goal", AllowDefault: true, DefaultAnswer: "still-missing"}
		},
		ApplyAnswer: func(*fakeState, Question, string, []string) error {
			answers++
			return nil
		},
	})
	if !errors.Is(err, ErrNeedsInput) {
		t.Fatalf("Run error = %v, want ErrNeedsInput", err)
	}
	if answers != 2 {
		t.Fatalf("answers = %d, want one defaulted answer per attempt", answers)
	}
}

func eventTypes(events []transcript.Event) string {
	var values []string
	for _, event := range events {
		values = append(values, event.Type)
	}
	return strings.Join(values, ",")
}
