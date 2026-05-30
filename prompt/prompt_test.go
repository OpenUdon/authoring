package prompt

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

func TestSessionShowsAutoAcceptedDefaults(t *testing.T) {
	var out strings.Builder
	prompts := NewSession(strings.NewReader("manual\n"), &out)
	prompts.SetDefaultMode(DefaultsShow)

	value, err := prompts.AskDefault("Choose operation", "getWeather")
	if err != nil || value != "getWeather" {
		t.Fatalf("AskDefault = %q, %v; want default", value, err)
	}
	yes, err := prompts.AskYesNo("Use API?", true)
	if err != nil || !yes {
		t.Fatalf("AskYesNo = %v, %v; want default true", yes, err)
	}
	optional, err := prompts.AskOptionalDefault("Optional timeout", "")
	if err != nil || optional != "" {
		t.Fatalf("AskOptionalDefault = %q, %v; want blank default", optional, err)
	}
	required, err := prompts.Ask("Workflow goal")
	if err != nil || required != "manual" {
		t.Fatalf("Ask = %q, %v; want manual input", required, err)
	}

	output := out.String()
	for _, expected := range []string{"Choose operation [getWeather]: getWeather", "Use API? [Y/n]: yes", "Optional timeout:", "Workflow goal: "} {
		if !strings.Contains(output, expected) {
			t.Fatalf("prompt output missing %q:\n%s", expected, output)
		}
	}
	turns := prompts.Turns()
	if len(turns) != 4 || turns[0].Answer != "getWeather" || turns[1].Answer != "yes" || turns[2].Answer != "" || turns[3].Answer != "manual" {
		t.Fatalf("turns = %#v", turns)
	}
	if turns[0].ID != "000001" || turns[3].ID != "000004" {
		t.Fatalf("turn IDs = %#v, want stable sequence IDs", turns)
	}
}

func TestSessionHidesAutoAcceptedDefaults(t *testing.T) {
	var out strings.Builder
	prompts := NewSession(strings.NewReader("manual\n"), &out)
	prompts.SetDefaultMode(DefaultsSilent)

	if value, err := prompts.AskDefault("Choose operation", "getWeather"); err != nil || value != "getWeather" {
		t.Fatalf("AskDefault = %q, %v; want default", value, err)
	}
	if yes, err := prompts.AskYesNo("Use API?", true); err != nil || !yes {
		t.Fatalf("AskYesNo = %v, %v; want default true", yes, err)
	}
	if required, err := prompts.Ask("Workflow goal"); err != nil || required != "manual" {
		t.Fatalf("Ask = %q, %v; want manual input", required, err)
	}
	if strings.Contains(out.String(), "Choose operation") || strings.Contains(out.String(), "Use API?") {
		t.Fatalf("silent defaults were printed:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "Workflow goal: ") {
		t.Fatalf("required prompt was not printed:\n%s", out.String())
	}
}

func TestRequiredAndYesNoPrompts(t *testing.T) {
	var out strings.Builder
	prompts := NewSession(strings.NewReader("\nvalue\nmaybe\nno\n"), &out)

	value, err := prompts.AskDefaultRequired("Name", "")
	if err != nil || value != "value" {
		t.Fatalf("AskDefaultRequired = %q, %v; want value", value, err)
	}
	yes, err := prompts.AskYesNo("Continue?", true)
	if err != nil || yes {
		t.Fatalf("AskYesNo = %v, %v; want false after invalid then no", yes, err)
	}
	if !strings.Contains(out.String(), "Name is required.") {
		t.Fatalf("required prompt did not explain missing value:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "Please answer yes or no.") {
		t.Fatalf("yes/no prompt did not reject invalid input:\n%s", out.String())
	}
}

func TestRequiredPromptReturnsEOF(t *testing.T) {
	var out strings.Builder
	prompts := NewSession(strings.NewReader("\n"), &out)

	_, err := prompts.AskDefaultRequired("Name", "")
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("AskDefaultRequired error = %v, want unexpected EOF", err)
	}
}

func TestReplayScriptAndLabelAssertions(t *testing.T) {
	turns := []session.PromptTurn{
		{Label: "Workflow goal", Answer: "run it"},
		{Label: "Use API?", Answer: "yes"},
	}
	script := ScriptFromTurns(turns)
	if script.Input != "run it\nyes\n" {
		t.Fatalf("script input = %q", script.Input)
	}
	data, err := io.ReadAll(script.Reader())
	if err != nil {
		t.Fatalf("read script input: %v", err)
	}
	if string(data) != script.Input {
		t.Fatalf("reader data = %q, want script input", string(data))
	}
	if err := AssertLabelsInOrder("Workflow goal: run it\nUse API? [Y/n]: yes\n", turns); err != nil {
		t.Fatalf("AssertLabelsInOrder returned error: %v", err)
	}
	if err := AssertLabelsInOrder("Use API?\nWorkflow goal\n", turns); err == nil {
		t.Fatalf("AssertLabelsInOrder accepted out-of-order labels")
	}
}

func TestPromptTranscriptSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prompt", "transcript.json")
	turns := []session.PromptTurn{
		{Label: "Z prompt", Answer: "last"},
		{Label: "A prompt", Answer: "first"},
	}
	record := NewTranscript("s1", turns, []transcript.Event{{Type: "draft", Message: "created"}}, session.State{Goal: "run it"})

	if err := SaveTranscript(path, record); err != nil {
		t.Fatalf("SaveTranscript returned error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat transcript: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("transcript permissions = %o, want 0600", got)
	}
	loaded, err := LoadTranscript(path)
	if err != nil {
		t.Fatalf("LoadTranscript returned error: %v", err)
	}
	if loaded.Version != TranscriptVersion {
		t.Fatalf("loaded version = %q, want %q", loaded.Version, TranscriptVersion)
	}
	if loaded.Session.Turns[0].Answer != "last" || loaded.Session.Turns[1].Answer != "first" {
		t.Fatalf("loaded session turns = %#v", loaded.Session.Turns)
	}
	if loaded.Transcript.Turns[0].Content != "last" || loaded.Transcript.Turns[1].Content != "first" {
		t.Fatalf("loaded transcript turns = %#v", loaded.Transcript.Turns)
	}
}
