package authoring_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/OpenUdon/authoring/icot"
	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/readiness"
	"github.com/OpenUdon/authoring/report"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/structured"
	"github.com/OpenUdon/authoring/transcript"
	"github.com/OpenUdon/authoring/trust"
)

func Example_manualPrompting() {
	var out bytes.Buffer
	prompts := prompt.NewSession(strings.NewReader("draft name\n"), &out)
	answer, _ := prompts.Ask("name")

	fmt.Printf("%s|%s", answer, out.String())
	// Output: draft name|name:
}

func Example_structuredJSONFallback() {
	type draft struct {
		Name string `json:"name"`
	}
	client := structured.ClientFunc(func(context.Context, []transcript.Turn) (transcript.Turn, error) {
		return transcript.Turn{Role: "assistant", Content: "```json\n{\"name\":\"example\"}\n```"}, nil
	})
	var out draft
	result, _ := structured.CompleteJSON(
		context.Background(),
		client,
		[]transcript.Turn{{Role: "user", Content: "Return a draft."}},
		map[string]any{"type": "object"},
		&out,
		structured.Options{DisableStructuredCompletion: true},
	)

	fmt.Printf("%s %s", result.Mode, out.Name)
	// Output: legacy example
}

func Example_progressiveLoop() {
	run, _ := icot.RunRuntime[fakeState, string, string](
		context.Background(),
		strings.NewReader(""),
		io.Discard,
		fakeRuntime{},
		icot.RuntimeConfig[fakeState, string]{
			Session:     fakeState{Goal: "write a draft"},
			Documents:   []string{"source"},
			MaxAttempts: 2,
		},
	)

	fmt.Printf("%t %s", run.Completed, run.Artifact)
	// Output: true artifact:write a draft
}

func Example_reportAndDigest() {
	digest := trust.SHA256Bytes([]byte("draft"))
	result := report.Normalize(report.Result{
		Status:  report.StatusComplete,
		Summary: "fake run",
		Digests: []trust.DigestRecord{digest},
	})

	fmt.Printf("%s %s %s", result.Version, result.Status, result.Digests[0].Algorithm)
	// Output: authoring.agent-result.v1 complete sha256
}

type fakeState struct {
	Goal    string
	Drafted bool
}

type fakeRuntime struct{}

func (fakeRuntime) Normalize(state *fakeState) {
	state.Goal = strings.TrimSpace(state.Goal)
}

func (fakeRuntime) Draft(_ context.Context, state fakeState, _ []string, _ []session.ReadinessIssue, _ int) (fakeState, error) {
	state.Drafted = true
	return state, nil
}

func (fakeRuntime) Readiness(state fakeState, _ []string) []session.ReadinessIssue {
	if state.Goal == "" {
		return []session.ReadinessIssue{{
			Code:     "goal_required",
			Severity: readiness.SeverityBlocking,
			Slot:     "goal",
			Message:  "Goal is required.",
		}}
	}
	return nil
}

func (fakeRuntime) Ready(state fakeState, issues []session.ReadinessIssue) bool {
	return readiness.Ready(issues) && state.Drafted
}

func (fakeRuntime) ShouldDraft(state fakeState, _ []string, issues []session.ReadinessIssue, _ int) bool {
	return !state.Drafted && readiness.Ready(issues)
}

func (fakeRuntime) PlanQuestion(_ fakeState, _ []string, issues []session.ReadinessIssue) icot.Question {
	if top := readiness.TopIssue(issues); top != nil {
		return icot.Question{ID: top.Code, Prompt: top.Message, Slots: []string{top.Slot}, Required: true}
	}
	return icot.Question{}
}

func (fakeRuntime) ApplyAnswer(state *fakeState, question icot.Question, answer string, _ []string) error {
	for _, slot := range question.Slots {
		if slot == "goal" {
			state.Goal = answer
		}
	}
	return nil
}

func (fakeRuntime) WriteArtifacts(_ context.Context, state *fakeState, _ []string, _ *[]transcript.Event) (string, error) {
	return "artifact:" + state.Goal, nil
}
