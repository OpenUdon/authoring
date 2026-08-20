package icot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/authoring/lifecycle"
	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

type fakeState struct {
	Goal  string
	Ready bool
}

type frontierState struct {
	Settled map[string]string
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

func TestLoopEntryPointsRejectNilContextsAndPromptSessions(t *testing.T) {
	if _, err := Run[fakeState, string, string](nil, nil, nil, Options[fakeState, string, string]{}); err == nil || !strings.Contains(err.Error(), "context is required") {
		t.Fatalf("nil Run context error = %v", err)
	}
	if _, err := runWithPromptSession[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{}); err == nil || !strings.Contains(err.Error(), "prompt session is required") {
		t.Fatalf("nil prompt-session error = %v", err)
	}
	var prompts *PromptSession
	if _, err := prompts.Ask("Goal"); err == nil || !strings.Contains(err.Error(), "prompt session is required") {
		t.Fatalf("nil compatibility prompt error = %v", err)
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
	if !strings.Contains(eventTypes(result.Events), "draft_error") || !strings.Contains(eventTypes(result.Events), "round_applied") {
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

func TestRunCancelStopsRoundCollectionImmediately(t *testing.T) {
	applied := false
	result, err := Run[frontierState, string, string](context.Background(), strings.NewReader("cancel\nsecond-answer\n"), nil, Options[frontierState, string, string]{
		Session:     frontierState{Settled: map[string]string{}},
		DefaultMode: prompt.DefaultsAsk,
		CheckReadiness: func(frontierState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "decisions.open", Severity: "blocking"}}
		},
		PlanFrontier: func(frontierState, []string, []session.ReadinessIssue) []Question {
			return []Question{
				{ID: "first", Prompt: "First", Required: true, Forced: true},
				{ID: "second", Prompt: "Second", Required: true, Forced: true},
			}
		},
		ApplyRound: func(*frontierState, []RoundAnswer, []string) error {
			applied = true
			return nil
		},
	})
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("Run error = %v, want ErrCanceled", err)
	}
	if applied || len(result.Answers) != 0 || len(result.Turns) != 1 {
		t.Fatalf("cancel result = %#v, applied=%t", result, applied)
	}
}

func TestRunStopsAfterThreeConsecutiveNoProgressRounds(t *testing.T) {
	var answers int
	_, err := Run[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{
		DefaultMode: prompt.DefaultsSilent,
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
	if !errors.Is(err, ErrNoProgress) {
		t.Fatalf("Run error = %v, want ErrNoProgress", err)
	}
	if answers != 3 {
		t.Fatalf("answers = %d, want one defaulted answer per no-progress round", answers)
	}
}

func TestRunAppliesMoreThanTwentyQuestionsAsOneAtomicRound(t *testing.T) {
	var questions []Question
	for i := 1; i <= 25; i++ {
		questions = append(questions, Question{
			ID:             fmt.Sprintf("q%02d", i),
			Prompt:         fmt.Sprintf("Question %02d", i),
			Recommendation: fmt.Sprintf("answer-%02d", i),
			Priority:       100 - i,
			Rationale:      "independent ready decision",
		})
	}
	var out strings.Builder
	autosaves := 0
	applyCalls := 0
	result, err := Run[frontierState, string, string](context.Background(), nil, &out, Options[frontierState, string, string]{
		Session:     frontierState{Settled: map[string]string{}},
		DefaultMode: prompt.DefaultsShow,
		CheckReadiness: func(state frontierState, _ []string) []session.ReadinessIssue {
			if len(state.Settled) == 25 {
				return nil
			}
			return []session.ReadinessIssue{{Code: "decisions.open", Severity: "blocking"}}
		},
		PlanFrontier: func(frontierState, []string, []session.ReadinessIssue) []Question {
			return questions
		},
		ApplyRound: func(state *frontierState, answers []RoundAnswer, _ []string) error {
			applyCalls++
			for _, answer := range answers {
				state.Settled[answer.QuestionID] = answer.Value
			}
			return nil
		},
		Autosave: func(frontierState) error { autosaves++; return nil },
		FinalConfirm: func(context.Context, *frontierState, []string, *[]transcript.Event) (string, error) {
			return "saved", nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Completed || len(result.Session.Settled) != 25 || applyCalls != 1 || autosaves != 1 {
		t.Fatalf("result=%#v applyCalls=%d autosaves=%d", result, applyCalls, autosaves)
	}
	output := out.String()
	lastPlanned := strings.Index(output, "25. Question 25")
	firstCollection := strings.Index(output, "Question 01 [answer-01]")
	if lastPlanned < 0 || firstCollection < 0 || lastPlanned > firstCollection {
		t.Fatalf("frontier was not fully displayed before collection:\n%s", output)
	}
}

func TestRunRecomputesDependentFrontierAfterEachRound(t *testing.T) {
	var planned []string
	result, err := Run[frontierState, string, string](context.Background(), nil, nil, Options[frontierState, string, string]{
		Session:     frontierState{Settled: map[string]string{}},
		DefaultMode: prompt.DefaultsSilent,
		CheckReadiness: func(state frontierState, _ []string) []session.ReadinessIssue {
			if len(state.Settled) == 3 {
				return nil
			}
			return []session.ReadinessIssue{{Code: "decisions.open", Severity: "blocking"}}
		},
		PlanFrontier: func(state frontierState, _ []string, _ []session.ReadinessIssue) []Question {
			if len(state.Settled) == 0 {
				planned = append(planned, "a,b")
				return []Question{{ID: "a", Prompt: "A", Recommendation: "a"}, {ID: "b", Prompt: "B", Recommendation: "b"}}
			}
			planned = append(planned, "c")
			return []Question{{ID: "c", Prompt: "C", Recommendation: "c"}}
		},
		ApplyRound: func(state *frontierState, answers []RoundAnswer, _ []string) error {
			for _, answer := range answers {
				state.Settled[answer.QuestionID] = answer.Value
			}
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(planned, ";"); got != "a,b;c" {
		t.Fatalf("planned frontiers = %q", got)
	}
	if result.Rounds != 3 || !result.Completed {
		t.Fatalf("result = %#v", result)
	}
}

func TestRunFastModeShowsOnlyQuestionsThatNeedInput(t *testing.T) {
	var out strings.Builder
	result, err := Run[frontierState, string, string](context.Background(), strings.NewReader("manual\n"), &out, Options[frontierState, string, string]{
		Session:     frontierState{Settled: map[string]string{}},
		DefaultMode: prompt.DefaultsSilent,
		CheckReadiness: func(state frontierState, _ []string) []session.ReadinessIssue {
			if len(state.Settled) == 2 {
				return nil
			}
			return []session.ReadinessIssue{{Code: "open", Severity: "blocking"}}
		},
		PlanFrontier: func(frontierState, []string, []session.ReadinessIssue) []Question {
			return []Question{
				{ID: "safe", Prompt: "Safe default", Recommendation: "safe"},
				{ID: "forced", Prompt: "Forced choice", Recommendation: "review", Forced: true},
			}
		},
		ApplyRound: func(state *frontierState, answers []RoundAnswer, _ []string) error {
			for _, answer := range answers {
				state.Settled[answer.QuestionID] = answer.Value
			}
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "Safe default") || !strings.Contains(out.String(), "Forced choice") {
		t.Fatalf("fast output = %q", out.String())
	}
	if result.Session.Settled["safe"] != "safe" || result.Session.Settled["forced"] != "manual" {
		t.Fatalf("settled = %#v", result.Session.Settled)
	}
}

func TestRunSemanticFingerprintIgnoresVolatileState(t *testing.T) {
	type volatileState struct {
		Goal      string
		UpdatedAt time.Time
	}
	apply := func(state *volatileState, _ Question, _ string, _ []string) error {
		state.UpdatedAt = state.UpdatedAt.Add(time.Second)
		return nil
	}
	base := Options[volatileState, string, string]{
		Session:     volatileState{Goal: "unchanged", UpdatedAt: time.Unix(0, 0)},
		DefaultMode: prompt.DefaultsSilent,
		MaxRounds:   3,
		CheckReadiness: func(volatileState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(volatileState, []string, []session.ReadinessIssue) Question {
			return Question{ID: "goal", Prompt: "Goal", Recommendation: "unchanged"}
		},
		ApplyAnswer: apply,
	}
	_, err := Run(context.Background(), nil, nil, base)
	if !errors.Is(err, ErrRoundLimit) || errors.Is(err, ErrNoProgress) {
		t.Fatalf("default volatile fingerprint error = %v, want round limit", err)
	}
	base.ProgressFingerprint = func(state volatileState, _ []string, issues []session.ReadinessIssue) (string, error) {
		data, err := json.Marshal(struct {
			Goal   string
			Issues []session.ReadinessIssue
		}{state.Goal, issues})
		return string(data), err
	}
	_, err = Run(context.Background(), nil, nil, base)
	if !errors.Is(err, ErrNoProgress) || !errors.Is(err, ErrNeedsInput) {
		t.Fatalf("semantic fingerprint error = %v, want no progress and needs input", err)
	}
}

func TestRunRejectsUnmarshalableFingerprintAndInvalidRoundLimit(t *testing.T) {
	type badState struct{ Values chan string }
	_, err := Run(context.Background(), nil, nil, Options[badState, string, string]{
		Session: badState{Values: make(chan string)},
		CheckReadiness: func(badState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
	})
	if err == nil || !strings.Contains(err.Error(), "compute progress fingerprint") || !strings.Contains(err.Error(), "unsupported type") {
		t.Fatalf("unmarshalable fingerprint error = %v", err)
	}
	_, err = Run[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{MaxRounds: -1})
	if err == nil || !strings.Contains(err.Error(), "nonnegative") {
		t.Fatalf("negative max-round error = %v", err)
	}
	if got, err := resolveMaxRounds(0); err != nil || got != DefaultMaxRounds {
		t.Fatalf("resolveMaxRounds(0) = %d, %v", got, err)
	}
	if DefaultMaxRounds != 1000 {
		t.Fatalf("DefaultMaxRounds = %d", DefaultMaxRounds)
	}
}

func TestRunRoundFuseAndSilentCancelRecommendation(t *testing.T) {
	applied := 0
	_, err := Run[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{
		DefaultMode: prompt.DefaultsSilent,
		MaxRounds:   2,
		CheckReadiness: func(fakeState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(state fakeState, _ []string, _ []session.ReadinessIssue) Question {
			return Question{ID: "goal", Prompt: "Goal", Recommendation: fmt.Sprintf("value-%d", applied+1)}
		},
		ApplyAnswer: func(state *fakeState, _ Question, answer string, _ []string) error {
			applied++
			state.Goal = answer
			return nil
		},
	})
	if !errors.Is(err, ErrNeedsInput) || !errors.Is(err, ErrRoundLimit) || applied != 2 {
		t.Fatalf("round fuse error=%v applied=%d", err, applied)
	}

	var got RoundAnswer
	result, err := Run[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{
		DefaultMode: prompt.DefaultsSilent,
		CheckReadiness: func(state fakeState, _ []string) []session.ReadinessIssue {
			if state.Ready {
				return nil
			}
			return []session.ReadinessIssue{{Code: "missing", Severity: "blocking"}}
		},
		PlanQuestion: func(fakeState, []string, []session.ReadinessIssue) Question {
			return Question{ID: "choice", Prompt: "Choice", Recommendation: "cancel"}
		},
		ApplyRound: func(state *fakeState, answers []RoundAnswer, _ []string) error {
			got = answers[0]
			state.Ready = true
			return nil
		},
	})
	if err != nil || !result.Completed || got.Value != "cancel" || got.Source != DefaultRecommendationSource {
		t.Fatalf("silent cancel result=%#v answer=%#v err=%v", result, got, err)
	}
}

func TestRunSequencesEventsThroughNormalizationAndPersistence(t *testing.T) {
	var observed []transcript.Event
	result, err := Run[fakeState, string, string](context.Background(), nil, nil, Options[fakeState, string, string]{
		Session:        fakeState{Ready: true},
		CheckReadiness: func(fakeState, []string) []session.ReadinessIssue { return nil },
		OnEvent:        func(event transcript.Event) { observed = append(observed, event) },
		FinalConfirm: func(_ context.Context, _ *fakeState, _ []string, events *[]transcript.Event) (string, error) {
			*events = append(*events, transcript.Event{Type: "adapter_event"})
			return "done", nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Events) != len(observed) || len(result.Events) < 3 {
		t.Fatalf("result events=%#v observed=%#v", result.Events, observed)
	}
	for index, event := range result.Events {
		want := fmt.Sprintf("%06d", index+1)
		if event.ID != want || observed[index].ID != want {
			t.Fatalf("event %d IDs result=%q observed=%q", index, event.ID, observed[index].ID)
		}
	}
	normalized := transcript.Normalize(transcript.Record{Events: result.Events})
	for index := range result.Events {
		if normalized.Events[index].Type != result.Events[index].Type || normalized.Events[index].ID != result.Events[index].ID {
			t.Fatalf("normalize reordered events: %#v", normalized.Events)
		}
	}
	path := filepath.Join(t.TempDir(), "events.json")
	if err := lifecycle.SaveTranscript(path, normalized); err != nil {
		t.Fatal(err)
	}
	loaded, ok, err := lifecycle.LoadTranscript(path)
	if err != nil || !ok {
		t.Fatalf("load transcript ok=%t err=%v", ok, err)
	}
	for index := range result.Events {
		if loaded.Events[index].Type != result.Events[index].Type || loaded.Events[index].ID != result.Events[index].ID {
			t.Fatalf("persistence reordered events: %#v", loaded.Events)
		}
	}
}

func TestValidateRoundAnswersRejectsUnknownIDs(t *testing.T) {
	questions := []Question{{ID: "shown", Prompt: "Shown"}}
	if err := validateRoundAnswers([]RoundAnswer{{QuestionID: "other", Value: "x"}}, questions); err == nil {
		t.Fatal("unknown answer ID was accepted")
	}
}

func eventTypes(events []transcript.Event) string {
	var values []string
	for _, event := range events {
		values = append(values, event.Type)
	}
	return strings.Join(values, ",")
}
