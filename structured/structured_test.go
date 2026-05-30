package structured

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/transcript"
)

type output struct {
	Name string `json:"name"`
}

type fakeClient struct {
	reply         string
	structuredErr error
	legacyCalls   int
	structCalls   int
	lastSchema    json.RawMessage
	lastTurns     []transcript.Turn
}

func (f *fakeClient) Complete(_ context.Context, turns []transcript.Turn) (transcript.Turn, error) {
	f.legacyCalls++
	f.lastTurns = append([]transcript.Turn(nil), turns...)
	return transcript.Turn{Role: "assistant", Content: f.reply}, nil
}

func (f *fakeClient) CompleteStructured(_ context.Context, turns []transcript.Turn, schema json.RawMessage, out any) error {
	f.structCalls++
	f.lastTurns = append([]transcript.Turn(nil), turns...)
	f.lastSchema = append(json.RawMessage(nil), schema...)
	if f.structuredErr != nil {
		return f.structuredErr
	}
	target := out.(*output)
	target.Name = "structured"
	return nil
}

func TestCompleteJSONStructuredSuccess(t *testing.T) {
	client := &fakeClient{}
	var got output

	result, err := CompleteJSON(context.Background(), client, nil, map[string]any{"type": "object"}, &got, Options{})
	if err != nil {
		t.Fatalf("CompleteJSON returned error: %v", err)
	}
	if result.Mode != ModeStructured || got.Name != "structured" {
		t.Fatalf("result = %#v output = %#v, want structured output", result, got)
	}
	if client.legacyCalls != 0 || client.structCalls != 1 {
		t.Fatalf("calls = structured %d legacy %d", client.structCalls, client.legacyCalls)
	}
	if string(client.lastSchema) != `{"type":"object"}` {
		t.Fatalf("schema = %s, want normalized schema", client.lastSchema)
	}
}

func TestCompleteJSONFallbackOnStructuredError(t *testing.T) {
	client := &fakeClient{
		reply:         "```json\n{\"name\":\"legacy\"}\n```",
		structuredErr: errors.New("structured unavailable"),
	}
	var got output

	result, err := CompleteJSON(context.Background(), client, []transcript.Turn{{Role: "user", Content: "make it"}}, `{"type":"object"}`, &got, Options{FallbackOnStructuredError: true})
	if err != nil {
		t.Fatalf("CompleteJSON returned error: %v", err)
	}
	if result.Mode != ModeLegacy || got.Name != "legacy" {
		t.Fatalf("result = %#v output = %#v, want legacy output", result, got)
	}
	if !strings.Contains(client.lastTurns[0].Content, "Return only JSON") {
		t.Fatalf("legacy instruction was not appended: %#v", client.lastTurns)
	}
}

func TestCompleteJSONLegacyExtraction(t *testing.T) {
	client := ClientFunc(func(context.Context, []transcript.Turn) (transcript.Turn, error) {
		return transcript.Turn{Role: "assistant", Content: "answer: {\"name\":\"legacy\"}"}, nil
	})
	var got output

	result, err := CompleteJSON(context.Background(), client, nil, nil, &got, Options{})
	if err != nil {
		t.Fatalf("CompleteJSON returned error: %v", err)
	}
	if result.Mode != ModeLegacy || result.Raw != `{"name":"legacy"}` || got.Name != "legacy" {
		t.Fatalf("result = %#v output = %#v", result, got)
	}
}

func TestInvalidJSONAndBadTargets(t *testing.T) {
	client := ClientFunc(func(context.Context, []transcript.Turn) (transcript.Turn, error) {
		return transcript.Turn{Role: "assistant", Content: "not json"}, nil
	})
	var got output
	if _, err := CompleteJSON(context.Background(), nil, nil, nil, &got, Options{}); err == nil {
		t.Fatalf("CompleteJSON accepted nil client")
	}
	if _, err := CompleteJSON(context.Background(), client, nil, nil, nil, Options{}); err == nil {
		t.Fatalf("CompleteJSON accepted nil output target")
	}
	if _, err := CompleteJSON(context.Background(), client, nil, nil, &got, Options{}); err == nil {
		t.Fatalf("CompleteJSON accepted invalid JSON response")
	}
	if err := DecodeJSONBlock(`{"name":"ok"}`, nil); err == nil {
		t.Fatalf("DecodeJSONBlock accepted nil target")
	}
}

func TestSchemaNormalizationAndDecodeHelpers(t *testing.T) {
	schema, err := NormalizeSchema(" { \"type\" : \"object\", \"properties\" : { \"name\" : { \"type\" : \"string\" } } } ")
	if err != nil {
		t.Fatalf("NormalizeSchema returned error: %v", err)
	}
	if string(schema) != `{"properties":{"name":{"type":"string"}},"type":"object"}` {
		t.Fatalf("schema = %s, want compact deterministic JSON", schema)
	}
	if _, err := NormalizeSchema("{"); err == nil {
		t.Fatalf("NormalizeSchema accepted invalid JSON")
	}
	if jsonText, err := ExtractJSONBlock("prefix [1,{\"x\":\"}\"}] suffix"); err != nil || jsonText != `[1,{"x":"}"}]` {
		t.Fatalf("ExtractJSONBlock = %q, %v", jsonText, err)
	}
	var got output
	if err := DecodeJSONBlock("```json\n{\"name\":\"decoded\"}\n```", &got); err != nil {
		t.Fatalf("DecodeJSONBlock returned error: %v", err)
	}
	if got.Name != "decoded" {
		t.Fatalf("decoded output = %#v", got)
	}
}
