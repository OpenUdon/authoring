package icotcli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/OpenUdon/authoring/prompt"
)

// Flags holds common interactive iCoT CLI options. Product commands own their
// own artifact, provider-client, and domain-specific flags.
type Flags struct {
	Agent        bool
	JSON         bool
	Answers      string
	PromptMode   string
	NoLLM        bool
	NoTranscript bool
	Provider     string
	Model        string
	Temperature  float64
	Report       string
}

// AddFlags registers common iCoT flags on fs.
func AddFlags(fs *flag.FlagSet, flags *Flags) {
	if fs == nil || flags == nil {
		return
	}
	promptModeDefault := strings.TrimSpace(flags.PromptMode)
	if promptModeDefault == "" {
		promptModeDefault = "normal"
	}
	fs.BoolVar(&flags.Agent, "agent", false, "Run in noninteractive agent mode")
	fs.BoolVar(&flags.JSON, "json", false, "Emit JSON")
	fs.StringVar(&flags.Answers, "answers", "", "Replay answers/session JSON path")
	fs.StringVar(&flags.PromptMode, "prompt-mode", promptModeDefault, "Prompt mode: full, normal, or fast")
	fs.BoolVar(&flags.NoLLM, "no-llm", flags.NoLLM, "Disable optional model assistance")
	fs.BoolVar(&flags.NoTranscript, "no-transcript", false, "Do not write a transcript")
	fs.StringVar(&flags.Provider, "provider", "", "Optional model provider label")
	fs.StringVar(&flags.Model, "model", "", "Optional model name")
	fs.Float64Var(&flags.Temperature, "temperature", 0, "Optional model temperature")
	fs.StringVar(&flags.Report, "report", "", "Optional report output path")
}

// PromptDefaultMode converts the shared prompt-mode flag into prompt behavior.
func PromptDefaultMode(mode string) (prompt.DefaultMode, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "full":
		return prompt.DefaultsAsk, nil
	case "", "normal":
		return prompt.DefaultsShow, nil
	case "fast":
		return prompt.DefaultsSilent, nil
	default:
		return prompt.DefaultsAsk, fmt.Errorf("--prompt-mode must be full, normal, or fast")
	}
}
