package prompt

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/OpenUdon/authoring/lifecycle"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

const (
	// TranscriptVersion is the prompt transcript envelope version.
	TranscriptVersion = "authoring.prompt-transcript.v1"
)

// DefaultMode controls how prompt defaults are handled.
type DefaultMode int

const (
	// DefaultsAsk prints defaulted prompts and waits for input.
	DefaultsAsk DefaultMode = iota
	// DefaultsShow prints defaulted prompts and accepts defaults.
	DefaultsShow
	// DefaultsSilent accepts defaulted prompts without printing them.
	DefaultsSilent
)

// Session prompts on a reader/writer pair and records prompt turns.
type Session struct {
	reader      *bufio.Reader
	out         io.Writer
	turns       []session.PromptTurn
	defaultMode DefaultMode
}

// ReplayScript is a deterministic prompt replay fixture.
type ReplayScript struct {
	Turns []session.PromptTurn `json:"turns,omitempty"`
	Input string               `json:"input,omitempty"`
}

// PromptTranscript is a persisted local transcript for prompt replay and
// review. It embeds the M03 transcript and session records instead of defining
// product-specific state.
type PromptTranscript struct {
	Version    string            `json:"version"`
	TimeUTC    string            `json:"time_utc"`
	Transcript transcript.Record `json:"transcript"`
	Session    session.State     `json:"session,omitempty"`
}

// NewSession creates a local prompt session.
func NewSession(in io.Reader, out io.Writer) *Session {
	if in == nil {
		in = strings.NewReader("")
	}
	if out == nil {
		out = io.Discard
	}
	if reader, ok := in.(*bufio.Reader); ok {
		return &Session{reader: reader, out: out}
	}
	return &Session{reader: bufio.NewReader(in), out: out}
}

// SetDefaultMode controls whether defaulted prompts ask, auto-accept visibly,
// or auto-accept silently. Required free-form prompts still ask.
func (s *Session) SetDefaultMode(mode DefaultMode) {
	if s == nil {
		return
	}
	s.defaultMode = mode
}

// Ask prompts for a required free-form value.
func (s *Session) Ask(label string) (string, error) {
	label = strings.TrimSpace(label)
	fmt.Fprintf(s.out, "%s: ", label)
	value, err := s.next()
	s.record(session.PromptTurn{Label: label, Answer: value, Source: "user"}, err)
	return value, err
}

// AskDefault prompts for a value, returning current when the answer is blank.
func (s *Session) AskDefault(label, current string) (string, error) {
	return s.askDefault(label, current, false, false)
}

// AskDefaultForced prints a defaulted prompt and waits for user input even when
// the default mode would normally auto-accept the default.
func (s *Session) AskDefaultForced(label, current string) (string, error) {
	return s.askDefault(label, current, false, true)
}

// AskOptionalDefault prompts for an optional value, allowing automatic default
// acceptance even when the current value is blank.
func (s *Session) AskOptionalDefault(label, current string) (string, error) {
	return s.askDefault(label, current, true, false)
}

// AskDefaultRequired prompts until a non-empty value is available.
func (s *Session) AskDefaultRequired(label, current string) (string, error) {
	for {
		value, err := s.AskDefault(label, current)
		if err != nil {
			return "", fmt.Errorf("%s: %w", strings.TrimSpace(label), err)
		}
		value = strings.TrimSpace(value)
		if value != "" {
			return value, nil
		}
		fmt.Fprintf(s.out, "%s is required.\n", strings.TrimSpace(label))
	}
}

// AskYesNo prompts for a yes/no answer with a default.
func (s *Session) AskYesNo(label string, defaultYes bool) (bool, error) {
	label = strings.TrimSpace(label)
	suffix := "y/N"
	answer := "no"
	if defaultYes {
		suffix = "Y/n"
		answer = "yes"
	}
	if s.defaultMode != DefaultsAsk {
		if s.defaultMode == DefaultsShow {
			fmt.Fprintf(s.out, "%s [%s]: %s\n", label, suffix, answer)
		}
		s.record(session.PromptTurn{Label: label, Answer: answer, Default: answer, Source: "default"}, nil)
		return defaultYes, nil
	}
	for {
		fmt.Fprintf(s.out, "%s [%s]: ", label, suffix)
		value, err := s.next()
		if err != nil {
			s.record(session.PromptTurn{Label: label, Answer: value, Default: answer, Source: "user"}, err)
			return false, err
		}
		raw := value
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			s.record(session.PromptTurn{Label: label, Answer: raw, Default: answer, Source: "default"}, nil)
			return defaultYes, nil
		}
		switch value {
		case "y", "yes", "true", "allow", "allowed", "approve", "approved":
			s.record(session.PromptTurn{Label: label, Answer: raw, Default: answer, Source: "user"}, nil)
			return true, nil
		case "n", "no", "false", "deny", "denied":
			s.record(session.PromptTurn{Label: label, Answer: raw, Default: answer, Source: "user"}, nil)
			return false, nil
		default:
			s.record(session.PromptTurn{Label: label, Answer: raw, Default: answer, Source: "user"}, nil)
			fmt.Fprintln(s.out, "Please answer yes or no.")
		}
	}
}

// Turns returns a copy of recorded prompt turns.
func (s *Session) Turns() []session.PromptTurn {
	if s == nil {
		return nil
	}
	return append([]session.PromptTurn(nil), s.turns...)
}

func (s *Session) askDefault(label, current string, autoBlank, forced bool) (string, error) {
	label = strings.TrimSpace(label)
	current = strings.TrimSpace(current)
	mode := s.defaultMode
	if forced {
		mode = DefaultsAsk
	}
	if mode != DefaultsAsk && (current != "" || autoBlank) {
		if mode == DefaultsShow {
			if current == "" {
				fmt.Fprintf(s.out, "%s:\n", label)
			} else {
				fmt.Fprintf(s.out, "%s [%s]: %s\n", label, OneLine(current), OneLine(current))
			}
		}
		s.record(session.PromptTurn{Label: label, Answer: current, Default: current, Source: "default"}, nil)
		return current, nil
	}
	if current == "" {
		fmt.Fprintf(s.out, "%s: ", label)
	} else {
		fmt.Fprintf(s.out, "%s [%s]: ", label, OneLine(current))
	}
	value, err := s.next()
	if err != nil {
		s.record(session.PromptTurn{Label: label, Answer: value, Default: current, Source: "user"}, err)
		return "", err
	}
	if strings.TrimSpace(value) == "" {
		s.record(session.PromptTurn{Label: label, Answer: current, Default: current, Source: "default"}, nil)
		return current, nil
	}
	s.record(session.PromptTurn{Label: label, Answer: value, Default: current, Source: "user"}, nil)
	return value, nil
}

func (s *Session) next() (string, error) {
	line, err := s.reader.ReadString('\n')
	if err == io.EOF && line != "" {
		return strings.TrimRight(line, "\r\n"), nil
	}
	if err != nil {
		if err == io.EOF {
			return "", io.ErrUnexpectedEOF
		}
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func (s *Session) record(turn session.PromptTurn, err error) {
	if err != nil && strings.TrimSpace(turn.Answer) == "" {
		return
	}
	if strings.TrimSpace(turn.ID) == "" {
		turn.ID = sequenceID(len(s.turns) + 1)
	}
	turn.Label = strings.TrimSpace(turn.Label)
	turn.Answer = strings.TrimSpace(turn.Answer)
	turn.Default = strings.TrimSpace(turn.Default)
	turn.Source = strings.TrimSpace(turn.Source)
	if turn.Label == "" {
		return
	}
	s.turns = append(s.turns, turn)
}

// OneLine normalizes a prompt default for display.
func OneLine(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

// ScriptFromTurns builds a replay script from prompt turns.
func ScriptFromTurns(turns []session.PromptTurn) ReplayScript {
	turns = sequenceTurns(turns)
	script := ReplayScript{Turns: turns}
	var answers []string
	for _, turn := range turns {
		answers = append(answers, turn.Answer)
	}
	if len(answers) > 0 {
		script.Input = strings.Join(answers, "\n") + "\n"
	}
	return script
}

// Reader returns an input reader for script.
func (script ReplayScript) Reader() io.Reader {
	return strings.NewReader(script.Input)
}

// AssertLabelsInOrder verifies that prompt labels were emitted in replay
// order.
func AssertLabelsInOrder(output string, turns []session.PromptTurn) error {
	offset := 0
	for _, turn := range turns {
		label := strings.TrimSpace(turn.Label)
		if label == "" {
			continue
		}
		index := strings.Index(output[offset:], label)
		if index < 0 {
			return fmt.Errorf("prompt label %q not found after offset %d", label, offset)
		}
		offset += index + len(label)
	}
	return nil
}

// NewTranscript builds a prompt transcript envelope from M03 session and
// transcript records.
func NewTranscript(sessionID string, turns []session.PromptTurn, events []transcript.Event, state session.State) PromptTranscript {
	turns = sequenceTurns(turns)
	record := transcript.Record{
		SessionID: strings.TrimSpace(sessionID),
		Turns:     promptTurnsToTranscriptTurns(turns),
		Events:    append([]transcript.Event(nil), events...),
	}
	state.Turns = append([]session.PromptTurn(nil), turns...)
	return PromptTranscript{
		Version:    TranscriptVersion,
		TimeUTC:    time.Now().UTC().Format(time.RFC3339),
		Transcript: transcript.Normalize(record),
		Session:    session.Normalize(state),
	}
}

// SaveTranscript writes a prompt transcript with private-file permissions.
// Empty paths are ignored.
func SaveTranscript(path string, record PromptTranscript) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	record = normalizePromptTranscript(record)
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return lifecycle.AtomicWrite(path, data, 0o600)
}

// LoadTranscript reads a prompt transcript from path.
func LoadTranscript(path string) (PromptTranscript, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PromptTranscript{}, err
	}
	var record PromptTranscript
	if err := json.Unmarshal(data, &record); err != nil {
		return PromptTranscript{}, err
	}
	return normalizePromptTranscript(record), nil
}

func promptTurnsToTranscriptTurns(turns []session.PromptTurn) []transcript.Turn {
	out := make([]transcript.Turn, 0, len(turns))
	for _, turn := range turns {
		out = append(out, transcript.Turn{
			ID:       turn.ID,
			Role:     "operator",
			Label:    turn.Label,
			Content:  turn.Answer,
			Redacted: turn.Redacted,
		})
	}
	return out
}

func sequenceTurns(turns []session.PromptTurn) []session.PromptTurn {
	out := make([]session.PromptTurn, len(turns))
	for i, turn := range turns {
		turn.ID = strings.TrimSpace(turn.ID)
		if turn.ID == "" {
			turn.ID = sequenceID(i + 1)
		}
		out[i] = turn
	}
	return out
}

func sequenceID(index int) string {
	return fmt.Sprintf("%06d", index)
}

func normalizePromptTranscript(record PromptTranscript) PromptTranscript {
	record.Version = strings.TrimSpace(record.Version)
	if record.Version == "" {
		record.Version = TranscriptVersion
	}
	record.TimeUTC = strings.TrimSpace(record.TimeUTC)
	if record.TimeUTC == "" {
		record.TimeUTC = time.Now().UTC().Format(time.RFC3339)
	}
	record.Transcript = transcript.Normalize(record.Transcript)
	record.Session = session.Normalize(record.Session)
	return record
}
