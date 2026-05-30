package icotcli

import (
	"flag"
	"testing"

	"github.com/OpenUdon/authoring/prompt"
)

func TestPromptDefaultMode(t *testing.T) {
	tests := []struct {
		mode string
		want prompt.DefaultMode
	}{
		{mode: "full", want: prompt.DefaultsAsk},
		{mode: "normal", want: prompt.DefaultsShow},
		{mode: "fast", want: prompt.DefaultsSilent},
		{mode: "", want: prompt.DefaultsShow},
	}
	for _, tt := range tests {
		got, err := PromptDefaultMode(tt.mode)
		if err != nil {
			t.Fatalf("PromptDefaultMode(%q) returned error: %v", tt.mode, err)
		}
		if got != tt.want {
			t.Fatalf("PromptDefaultMode(%q) = %v, want %v", tt.mode, got, tt.want)
		}
	}
	if _, err := PromptDefaultMode("verbose"); err == nil {
		t.Fatalf("PromptDefaultMode accepted invalid mode")
	}
}

func TestFlagsParseSharedOptions(t *testing.T) {
	var flags Flags
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	AddFlags(fs, &flags)
	if err := fs.Parse([]string{"--agent", "--json", "--answers", "answers.json", "--prompt-mode", "fast", "--no-llm", "--provider", "openai", "--model", "gpt-test", "--temperature", "0.2", "--report", "report.json"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	if !flags.Agent || !flags.JSON || !flags.NoLLM || flags.Answers != "answers.json" || flags.PromptMode != "fast" || flags.Provider != "openai" || flags.Model != "gpt-test" || flags.Temperature != 0.2 || flags.Report != "report.json" {
		t.Fatalf("flags = %#v", flags)
	}
}

func TestAddFlagsRespectsSeedDefaults(t *testing.T) {
	flags := Flags{NoLLM: true, PromptMode: "fast"}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	AddFlags(fs, &flags)
	if err := fs.Parse(nil); err != nil {
		t.Fatalf("parse defaults: %v", err)
	}
	if !flags.NoLLM || flags.PromptMode != "fast" {
		t.Fatalf("flags = %#v, want seeded defaults", flags)
	}
}
