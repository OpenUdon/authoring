package structured

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/OpenUdon/authoring/transcript"
)

const (
	// ModeStructured reports that a structured-output client produced the
	// result.
	ModeStructured = "structured"
	// ModeLegacy reports that legacy text completion plus JSON extraction
	// produced the result.
	ModeLegacy = "legacy"
)

// Client is the provider-neutral text completion boundary.
type Client interface {
	Complete(context.Context, []transcript.Turn) (transcript.Turn, error)
}

// ClientFunc adapts a function to Client.
type ClientFunc func(context.Context, []transcript.Turn) (transcript.Turn, error)

// Complete calls fn.
func (fn ClientFunc) Complete(ctx context.Context, turns []transcript.Turn) (transcript.Turn, error) {
	return fn(ctx, turns)
}

// StructuredClient is the optional provider-neutral structured completion
// boundary.
type StructuredClient interface {
	Client
	CompleteStructured(context.Context, []transcript.Turn, json.RawMessage, any) error
}

// Options configures structured completion and fallback behavior.
type Options struct {
	LegacyInstruction           string
	FallbackOnStructuredError   bool
	DisableStructuredCompletion bool
}

// Result reports which completion path was used.
type Result struct {
	Mode string `json:"mode,omitempty"`
	Raw  string `json:"raw,omitempty"`
}

// CompleteJSON tries structured completion first, then optionally falls back
// to legacy text completion plus JSON extraction.
func CompleteJSON(ctx context.Context, client Client, turns []transcript.Turn, schema any, out any, opts Options) (Result, error) {
	if client == nil {
		return Result{}, fmt.Errorf("structured client is required")
	}
	target, err := outputTarget(out)
	if err != nil {
		return Result{}, err
	}
	rawSchema, err := NormalizeSchema(schema)
	if err != nil {
		return Result{}, err
	}
	if !opts.DisableStructuredCompletion {
		if structured, ok := client.(StructuredClient); ok {
			scratch := reflect.New(target.Elem().Type())
			err := structured.CompleteStructured(ctx, cloneTurns(turns), rawSchema, scratch.Interface())
			if err == nil {
				target.Elem().Set(scratch.Elem())
				return Result{Mode: ModeStructured}, nil
			}
			if !opts.FallbackOnStructuredError {
				return Result{Mode: ModeStructured}, err
			}
		}
	}
	legacyTurns := AppendLegacyJSONInstruction(turns, opts.LegacyInstruction)
	reply, err := client.Complete(ctx, legacyTurns)
	if err != nil {
		return Result{Mode: ModeLegacy}, err
	}
	jsonText, err := ExtractJSONBlock(reply.Content)
	if err != nil {
		return Result{Mode: ModeLegacy, Raw: reply.Content}, err
	}
	scratch := reflect.New(target.Elem().Type())
	if err := json.Unmarshal([]byte(jsonText), scratch.Interface()); err != nil {
		return Result{Mode: ModeLegacy, Raw: jsonText}, err
	}
	target.Elem().Set(scratch.Elem())
	return Result{Mode: ModeLegacy, Raw: jsonText}, nil
}

// NormalizeSchema returns compact deterministic JSON schema bytes. Empty nil
// schemas normalize to nil.
func NormalizeSchema(schema any) (json.RawMessage, error) {
	if schema == nil {
		return nil, nil
	}
	var data []byte
	switch typed := schema.(type) {
	case json.RawMessage:
		data = append([]byte(nil), typed...)
	case []byte:
		data = append([]byte(nil), typed...)
	case string:
		data = []byte(strings.TrimSpace(typed))
	default:
		var err error
		data, err = json.Marshal(typed)
		if err != nil {
			return nil, err
		}
	}
	data = []byte(strings.TrimSpace(string(data)))
	if len(data) == 0 {
		return nil, nil
	}
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}
	normalized, err := json.Marshal(decoded)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(normalized), nil
}

// ExtractJSONBlock extracts a JSON object or array from a raw model response.
func ExtractJSONBlock(response string) (string, error) {
	response = strings.TrimSpace(response)
	if response == "" {
		return "", fmt.Errorf("empty model response")
	}
	if strings.HasPrefix(response, "```") {
		lines := strings.Split(response, "\n")
		if len(lines) >= 3 {
			body := strings.Join(lines[1:len(lines)-1], "\n")
			return ExtractJSONBlock(body)
		}
	}
	start, end := jsonBounds(response)
	if start < 0 || end <= start {
		return "", fmt.Errorf("no JSON object or array found in model response")
	}
	return strings.TrimSpace(response[start : end+1]), nil
}

// DecodeJSONBlock extracts and decodes a JSON object or array from response.
func DecodeJSONBlock(response string, target any) error {
	if _, err := outputTarget(target); err != nil {
		return err
	}
	jsonText, err := ExtractJSONBlock(response)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonText), target)
}

// AppendLegacyJSONInstruction appends a JSON-only instruction to the last user
// turn if the instruction is not already present.
func AppendLegacyJSONInstruction(turns []transcript.Turn, instruction string) []transcript.Turn {
	instruction = strings.TrimSpace(instruction)
	if instruction == "" {
		instruction = "Return only JSON. Do not include Markdown."
	}
	out := cloneTurns(turns)
	for _, turn := range out {
		if strings.Contains(turn.Content, instruction) {
			return out
		}
	}
	for i := len(out) - 1; i >= 0; i-- {
		if strings.TrimSpace(strings.ToLower(out[i].Role)) == "user" {
			out[i].Content = strings.TrimSpace(out[i].Content) + "\n\n" + instruction
			return out
		}
	}
	return out
}

func outputTarget(out any) (reflect.Value, error) {
	if out == nil {
		return reflect.Value{}, fmt.Errorf("structured output target is required")
	}
	target := reflect.ValueOf(out)
	if target.Kind() != reflect.Pointer || target.IsNil() {
		return reflect.Value{}, fmt.Errorf("structured output target must be a non-nil pointer")
	}
	return target, nil
}

func jsonBounds(response string) (int, int) {
	start := -1
	open := byte(0)
	close := byte(0)
	for i := 0; i < len(response); i++ {
		switch response[i] {
		case '{':
			start, open, close = i, '{', '}'
		case '[':
			start, open, close = i, '[', ']'
		}
		if start >= 0 {
			break
		}
	}
	if start < 0 {
		return -1, -1
	}
	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(response); i++ {
		ch := response[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}
		if ch == '"' {
			inString = true
			continue
		}
		if ch == open {
			depth++
		}
		if ch == close {
			depth--
			if depth == 0 {
				return start, i
			}
		}
	}
	return start, -1
}

func cloneTurns(turns []transcript.Turn) []transcript.Turn {
	return append([]transcript.Turn(nil), turns...)
}
