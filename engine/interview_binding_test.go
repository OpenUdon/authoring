package engine

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/interview"
	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/session"
)

type boundInterviewState struct {
	Interview interview.State
	Values    map[string]string
}

func TestInterviewBindingSynchronizesNoProgressCounter(t *testing.T) {
	binding := InterviewBinding[boundInterviewState, string]{
		State: func(state *boundInterviewState) *interview.State { return &state.Interview },
		Clone: func(state boundInterviewState) (boundInterviewState, error) { return state, nil },
		Resolve: func(*boundInterviewState, []string, interview.Node, RoundAnswer) (interview.Resolution, error) {
			return interview.Resolution{}, errors.New("unexpected resolution")
		},
	}
	result, err := Run(context.Background(), strings.NewReader(""), nil, Options[boundInterviewState, string, string]{
		Interview: &binding,
		CheckReadiness: func(boundInterviewState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "blocked", Severity: "blocking"}}
		},
	})
	if !errors.Is(err, ErrNeedsInput) || !errors.Is(err, ErrNoProgress) {
		t.Fatalf("Run error = %v", err)
	}
	if result.NoProgressRounds != 3 || result.Session.Interview.NoProgressRounds != 3 {
		t.Fatalf("counter result=%d state=%d", result.NoProgressRounds, result.Session.Interview.NoProgressRounds)
	}
}

func TestInterviewBindingAutosavesSynchronizedCounterOncePerRound(t *testing.T) {
	binding := InterviewBinding[boundInterviewState, string]{
		State: func(state *boundInterviewState) *interview.State { return &state.Interview },
		Clone: func(state boundInterviewState) (boundInterviewState, error) {
			state.Interview = interview.Normalize(state.Interview)
			return state, nil
		},
		Prepare: func(state *boundInterviewState, _ []string) error {
			frontier, err := interview.Frontier(state.Interview)
			if err != nil {
				return err
			}
			if len(frontier) == 0 {
				id := fmt.Sprintf("node-%d", state.Interview.Round+1)
				state.Interview.Nodes = append(state.Interview.Nodes, interview.Node{ID: id, Prompt: id, Recommendation: "same"})
			}
			return nil
		},
		Resolve: func(state *boundInterviewState, _ []string, node interview.Node, answer RoundAnswer) (interview.Resolution, error) {
			resolved := interview.Answer{ID: fmt.Sprintf("answer-%d", state.Interview.Round+1), NodeID: node.ID, Value: answer.Value}
			return interview.Resolution{NodeID: node.ID, Answer: &resolved}, nil
		},
	}
	var saved []int
	_, err := Run(context.Background(), nil, nil, Options[boundInterviewState, string, string]{
		Interview:   &binding,
		DefaultMode: prompt.DefaultsSilent,
		CheckReadiness: func(boundInterviewState, []string) []session.ReadinessIssue {
			return []session.ReadinessIssue{{Code: "blocked", Severity: "blocking"}}
		},
		ProgressFingerprint: func(boundInterviewState, []string, []session.ReadinessIssue) (string, error) {
			return "semantic-state", nil
		},
		Autosave: func(state boundInterviewState) error {
			saved = append(saved, state.Interview.NoProgressRounds)
			return nil
		},
	})
	if !errors.Is(err, ErrNoProgress) || !slices.Equal(saved, []int{1, 2, 3}) {
		t.Fatalf("error=%v saved counters=%#v", err, saved)
	}
}

func testInterviewBinding(resolveErr error) InterviewBinding[boundInterviewState, string] {
	return InterviewBinding[boundInterviewState, string]{
		State: func(state *boundInterviewState) *interview.State { return &state.Interview },
		Clone: func(state boundInterviewState) (boundInterviewState, error) {
			clone := state
			clone.Interview = interview.Normalize(state.Interview)
			clone.Values = map[string]string{}
			for key, value := range state.Values {
				clone.Values[key] = value
			}
			return clone, nil
		},
		Resolve: func(state *boundInterviewState, _ []string, node interview.Node, answer RoundAnswer) (interview.Resolution, error) {
			state.Values[node.ID] = answer.Value
			if resolveErr != nil && node.ID == "b" {
				return interview.Resolution{}, resolveErr
			}
			resolved := interview.Answer{ID: "answer-" + node.ID, NodeID: node.ID, Value: answer.Value, Source: answer.Source}
			return interview.Resolution{NodeID: node.ID, Answer: &resolved}, nil
		},
	}
}

func TestInterviewBindingPlansAndAppliesCloneAtomically(t *testing.T) {
	state := boundInterviewState{Interview: interview.State{Nodes: []interview.Node{{ID: "a", Prompt: "A", Priority: 2}, {ID: "b", Prompt: "B", Priority: 1}}}, Values: map[string]string{}}
	binding := testInterviewBinding(nil)
	questions, err := binding.Plan(&state, nil)
	if err != nil || len(questions) != 2 || questions[0].ID != "a" || questions[1].ID != "b" {
		t.Fatalf("questions=%#v err=%v", questions, err)
	}
	err = binding.Apply(&state, []RoundAnswer{{QuestionID: "a", Value: "one"}, {QuestionID: "b", Value: "two"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if state.Values["a"] != "one" || state.Values["b"] != "two" || state.Interview.Round != 1 || state.Interview.NoProgressRounds != 0 {
		t.Fatalf("state = %#v", state)
	}
}

func TestInterviewBindingResolveFailureLeavesProductAndInterviewUnchanged(t *testing.T) {
	state := boundInterviewState{Interview: interview.Normalize(interview.State{NoProgressRounds: 2, Nodes: []interview.Node{{ID: "a"}, {ID: "b"}}}), Values: map[string]string{"existing": "value"}}
	want := state
	want.Values = map[string]string{"existing": "value"}
	failure := errors.New("resolution failed")
	err := testInterviewBinding(failure).Apply(&state, []RoundAnswer{{QuestionID: "a", Value: "one"}, {QuestionID: "b", Value: "two"}}, nil)
	if !errors.Is(err, failure) {
		t.Fatalf("error = %v", err)
	}
	if state.Interview.Round != want.Interview.Round || state.Interview.NoProgressRounds != want.Interview.NoProgressRounds || len(state.Interview.Answers) != 0 || state.Values["existing"] != "value" || len(state.Values) != 1 {
		t.Fatalf("failed binding changed state: %#v", state)
	}
}

func TestInterviewBindingPlanFailureLeavesProductStateUnchanged(t *testing.T) {
	state := boundInterviewState{Interview: interview.State{Nodes: []interview.Node{{ID: "a"}}}, Values: map[string]string{"existing": "value"}}
	failure := errors.New("prepare failed")
	binding := testInterviewBinding(nil)
	binding.Prepare = func(clone *boundInterviewState, _ []string) error {
		clone.Values["mutated"] = "temporary"
		return failure
	}
	if _, err := binding.Plan(&state, nil); !errors.Is(err, failure) {
		t.Fatalf("Plan error = %v", err)
	}
	if len(state.Values) != 1 || state.Values["existing"] != "value" {
		t.Fatalf("failed Plan changed state: %#v", state)
	}
}
