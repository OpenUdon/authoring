package transcript

import (
	"bytes"
	"testing"

	"github.com/OpenUdon/authoring/decision"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/trust"
)

func TestNormalizeTranscriptDeterministic(t *testing.T) {
	record := Record{
		SessionID: " s1 ",
		Provider: &ModelProvenance{
			Provider: "OpenAI",
			Model:    " gpt-test ",
		},
		Turns: []Turn{
			{ID: "2", Role: "Assistant", Content: " done "},
			{
				ID:      "1",
				Role:    "User",
				Content: " goal ",
				Decisions: []session.DecisionEvidence{
					{Stage: "Operation Selection", Slot: "op", Value: "create", Source: "Model"},
				},
			},
		},
		Events: []Event{
			{Type: "draft", Severity: "info", Message: " later "},
			{Type: "readiness", Severity: "blocking", Message: " first "},
		},
		Diagnostics: []trust.DiagnosticRecord{
			{Code: "z", Severity: "info", Message: " later "},
			{Code: "a", Severity: "error", Message: " first "},
		},
		Metadata: map[string]string{" product ": " openudon "},
	}

	got := Normalize(record)
	if got.Version != Version {
		t.Fatalf("Version = %q, want %q", got.Version, Version)
	}
	if got.Provider == nil || got.Provider.Provider != "openai" || got.Provider.Model != "gpt-test" {
		t.Fatalf("provider = %#v, want normalized provider provenance", got.Provider)
	}
	if got.Turns[0].ID != "1" || got.Turns[0].Role != "user" {
		t.Fatalf("turns = %#v, want deterministic sorted normalized turns", got.Turns)
	}
	if got.Turns[0].Decisions[0].Stage != "operation_selection" {
		t.Fatalf("decision = %#v, want session-normalized decision", got.Turns[0].Decisions[0])
	}
	if got.Events[0].Severity != "blocking" {
		t.Fatalf("events = %#v, want blocking event first", got.Events)
	}
	if got.Diagnostics[0].Code != "a" {
		t.Fatalf("diagnostics = %#v, want Evidence diagnostic ordering", got.Diagnostics)
	}
	if got.Metadata["product"] != "openudon" {
		t.Fatalf("metadata = %#v, want trimmed metadata", got.Metadata)
	}
}

func TestCanonicalJSONStable(t *testing.T) {
	record := Record{
		SessionID: "s1",
		Events: []Event{
			{Type: "z"},
			{Type: "a"},
		},
	}

	first, err := CanonicalJSON(record)
	if err != nil {
		t.Fatalf("CanonicalJSON returned error: %v", err)
	}
	second, err := CanonicalJSON(record)
	if err != nil {
		t.Fatalf("CanonicalJSON returned error: %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("CanonicalJSON changed between calls:\n%s\n%s", first, second)
	}
	if !bytes.Contains(first, []byte(`"version": "authoring.transcript.v1"`)) {
		t.Fatalf("CanonicalJSON missing transcript version: %s", first)
	}
	if bytes.Contains(first, []byte(`"provider"`)) {
		t.Fatalf("CanonicalJSON emitted empty provider: %s", first)
	}
}

func TestProductShapedTranscriptMapsWithoutProductImports(t *testing.T) {
	openUdonShaped := struct {
		TurnLabel string
		Provider  string
		Model     string
	}{
		TurnLabel: "workflow_goal",
		Provider:  "anthropic",
		Model:     "claude-test",
	}
	ramenShaped := struct {
		EventType string
		Artifact  string
	}{
		EventType: "project_drafted",
		Artifact:  "project.uws.yaml",
	}

	mapped := Normalize(Record{
		SessionID: "cross-repo",
		Provider:  &ModelProvenance{Provider: openUdonShaped.Provider, Model: openUdonShaped.Model},
		Turns:     []Turn{{Label: openUdonShaped.TurnLabel, Role: "user", Content: "draft it"}},
		Events:    []Event{{Type: ramenShaped.EventType, Stage: "draft"}},
		Artifacts: []trust.ArtifactRecord{
			{Path: ramenShaped.Artifact, Kind: "project"},
		},
	})

	if mapped.Provider == nil || mapped.Provider.Model != "claude-test" {
		t.Fatalf("mapped provider = %#v", mapped.Provider)
	}
	if mapped.Events[0].Type != "project_drafted" {
		t.Fatalf("mapped event = %#v", mapped.Events[0])
	}
	if mapped.Artifacts[0].Path != "project.uws.yaml" {
		t.Fatalf("mapped artifact = %#v", mapped.Artifacts[0])
	}
}

func TestDecisionEvents(t *testing.T) {
	events := DecisionEvents([]decision.Record{{
		Stage:      "Operation Selection",
		Slot:       "resource.operation",
		Value:      "create",
		Source:     "model",
		Confidence: "low",
		Rationale:  "chosen from goal",
	}})
	if len(events) != 1 {
		t.Fatalf("DecisionEvents returned %d events, want 1", len(events))
	}
	event := events[0]
	if event.Type != "decision" || event.Stage != "operation_selection" || event.Fields["behavior"] != decision.BehaviorLowConfidence {
		t.Fatalf("decision event = %#v", event)
	}
	if event.Message != "chosen from goal" {
		t.Fatalf("decision event message = %q", event.Message)
	}
}
