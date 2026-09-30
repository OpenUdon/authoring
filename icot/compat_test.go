package icot_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/engine"
	"github.com/OpenUdon/authoring/icot"
	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

// Aliases must cross both API boundaries without conversion, including generic
// hooks whose callbacks contain aliased state, questions and prompt sessions.
var (
	_ engine.Options[int, string, string]                                                                                                                                = icot.Options[int, string, string]{}
	_ icot.Options[int, string, string]                                                                                                                                  = engine.Options[int, string, string]{}
	_ engine.InteractiveHooks[int, string, string]                                                                                                                       = icot.InteractiveHooks[int, string, string]{}
	_ icot.InterviewBinding[int, string]                                                                                                                                 = engine.InterviewBinding[int, string]{}
	_ engine.RuntimeConfig[int, string]                                                                                                                                  = icot.RuntimeConfig[int, string]{}
	_ icot.RepairOptions[int, string, string]                                                                                                                            = engine.RepairOptions[int, string, string]{}
	_ func(context.Context, io.Reader, io.Writer, engine.InteractiveHooks[int, string, string]) (string, error)                                                          = icot.RunInteractive[int, string, string]
	_ func(context.Context, io.Reader, io.Writer, engine.InteractiveHooks[int, string, string], engine.InteractiveLifecycleOptions[int, string, string]) (string, error) = icot.RunInteractiveWithLifecycle[int, string, string]
	_ func(context.Context, engine.RepairRuntime[int, string, string], engine.RepairConfig[int, string, string]) (engine.RepairResult[int], error)                       = icot.RunRuntimeRepair[int, string, string]
)

func TestFacadeSharesTypesAndSentinels(t *testing.T) {
	for _, pair := range [][2]reflect.Type{
		{reflect.TypeOf(icot.PromptSession{}), reflect.TypeOf(engine.PromptSession{})},
		{reflect.TypeOf(icot.Result[int, string]{}), reflect.TypeOf(engine.Result[int, string]{})},
		{reflect.TypeOf(icot.Event{}), reflect.TypeOf(engine.Event{})},
		{reflect.TypeOf(icot.RoundAnswer{}), reflect.TypeOf(engine.RoundAnswer{})},
	} {
		if pair[0] != pair[1] {
			t.Fatalf("separate concrete types: %v / %v", pair[0], pair[1])
		}
	}
	for _, pair := range [][2]error{
		{icot.ErrNeedsInput, engine.ErrNeedsInput}, {icot.ErrCanceled, engine.ErrCanceled},
		{icot.ErrNoProgress, engine.ErrNoProgress}, {icot.ErrRoundLimit, engine.ErrRoundLimit},
		{icot.ErrRepairExhausted, engine.ErrRepairExhausted}, {icot.ErrRepairNoop, engine.ErrRepairNoop},
	} {
		if pair[0] != pair[1] || !errors.Is(pair[0], pair[1]) {
			t.Fatal("sentinel identity changed")
		}
	}
	if icot.DefaultMaxRounds != 1000 || icot.DefaultNoProgressLimit != 3 || icot.DefaultRecommendationSource != "recommendation" {
		t.Fatal("loop defaults changed")
	}
}

func TestFacadeLoopParity(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		cancel      bool
		want        error
	}{
		{name: "answer", input: "report\n"},
		{name: "missing-input", want: engine.ErrNeedsInput},
		{name: "cancelled-context", cancel: true, want: engine.ErrCanceled},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancel {
				cancel()
			}
			var baseline []byte
			var baselineOutput string
			for index, run := range []func(context.Context, io.Reader, io.Writer, engine.Options[string, string, string]) (engine.Result[string, string], error){engine.Run[string, string, string], icot.Run[string, string, string]} {
				var output bytes.Buffer
				result, err := run(ctx, strings.NewReader(tc.input), &output, icot.Options[string, string, string]{
					CheckReadiness: func(state string, _ []string) []session.ReadinessIssue {
						if state != "" {
							return nil
						}
						return []session.ReadinessIssue{{Code: "missing-goal", Severity: "blocking"}}
					},
					PlanQuestion: func(string, []string, []session.ReadinessIssue) icot.Question {
						return icot.Question{ID: "goal", Prompt: "Goal", Required: true}
					},
					ApplyAnswer: func(state *string, _ icot.Question, value string, _ []string) error { *state = value; return nil },
					FinalConfirm: func(_ context.Context, state *string, _ []string, _ *[]transcript.Event) (string, error) {
						return *state, nil
					},
				})
				if tc.want != nil {
					if !errors.Is(err, tc.want) {
						t.Fatalf("error=%v want %v", err, tc.want)
					}
				} else if err != nil || !result.Completed || result.Artifact != "report" {
					t.Fatalf("result=%+v err=%v", result, err)
				}
				// Prompt timestamps describe the two actual invocations; compare every
				// other public field and the exact terminal output.
				for i := range result.Turns {
					result.Turns[i].TimeUTC = ""
				}
				for i := range result.Events {
					result.Events[i].TimeUTC = ""
				}
				data, err := json.Marshal(result)
				if err != nil {
					t.Fatal(err)
				}
				if index == 0 {
					baseline = data
					baselineOutput = output.String()
				} else if !bytes.Equal(data, baseline) || output.String() != baselineOutput {
					t.Fatalf("facade outcome/output diverged:\n%s\n%s", baseline, data)
				}
			}
		})
	}
}

func TestFacadePromptMethodsAndForcedDefaults(t *testing.T) {
	var baseline []byte
	var baselineOutput string
	for i, newSession := range []func(io.Reader, io.Writer) *engine.PromptSession{engine.NewPromptSession, icot.NewPromptSession} {
		var out bytes.Buffer
		p := newSession(strings.NewReader("explicit\n"), &out)
		p.SetDefaultMode(prompt.DefaultsSilent)
		if answer, err := p.AskDefault("Safe default", "ready"); err != nil || answer != "ready" {
			t.Fatalf("default %q %v", answer, err)
		}
		if answer, err := p.AskDefaultForced("Required decision", ""); err != nil || answer != "explicit" {
			t.Fatalf("forced %q %v", answer, err)
		}
		turns := p.Turns()
		for j := range turns {
			turns[j].TimeUTC = ""
		}
		data, err := json.Marshal(turns)
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			baseline = data
			baselineOutput = out.String()
		} else if !bytes.Equal(baseline, data) || out.String() != baselineOutput {
			t.Fatal("prompt output/turn parity changed")
		}
	}
}
