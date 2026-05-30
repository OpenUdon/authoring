package session

import (
	"bytes"
	"testing"

	"github.com/OpenUdon/authoring/trust"
	"github.com/OpenUdon/evidence/digest"
)

func TestNormalizeSessionDeterministic(t *testing.T) {
	state := State{
		Version: " ",
		ID:      " demo ",
		Mode:    "Ask User",
		Turns: []PromptTurn{
			{ID: "2", Label: "z", Question: " later "},
			{ID: "1", Label: "a", Question: " first "},
		},
		Answers: []Answer{
			{Slot: "resource.name", Value: " demo "},
			{Slot: "account", Value: " prod ", Source: "User Input"},
		},
		Readiness: []ReadinessIssue{
			{Code: "later", Severity: "info"},
			{Code: "missing_name", Severity: "blocking", Slot: "resource.name"},
		},
		Decisions: []DecisionEvidence{
			{Stage: "Operation Selection", Slot: "op", Value: "create", Source: "Model"},
			{Stage: "Catalog Plan", Slot: "source", Value: "openapi", Source: "User"},
		},
		Artifacts: []trust.ArtifactRecord{
			{Path: "z.json"},
			{Path: "a.json", Digest: digest.Record{Algorithm: "SHA256", Value: " abc "}},
		},
		Diagnostics: []trust.DiagnosticRecord{
			{Code: "z", Severity: "info", Message: " later "},
			{Code: "a", Severity: "error", Message: " first "},
		},
		Metadata: map[string]string{" product ": " ramen ", "empty": " "},
	}

	got := Normalize(state)
	if got.Version != Version {
		t.Fatalf("Version = %q, want %q", got.Version, Version)
	}
	if got.Mode != "ask_user" {
		t.Fatalf("Mode = %q, want ask_user", got.Mode)
	}
	if got.Turns[0].ID != "1" || got.Answers[0].Slot != "account" {
		t.Fatalf("records were not sorted deterministically: %#v %#v", got.Turns, got.Answers)
	}
	if got.Readiness[0].Code != "missing_name" {
		t.Fatalf("first readiness = %#v, want blocking issue first", got.Readiness[0])
	}
	if got.Decisions[0].Stage != "catalog_plan" {
		t.Fatalf("first decision stage = %q, want catalog_plan", got.Decisions[0].Stage)
	}
	if got.Artifacts[0].Path != "a.json" || got.Artifacts[0].Digest.Algorithm != "sha256" {
		t.Fatalf("artifact normalization = %#v", got.Artifacts[0])
	}
	if got.Diagnostics[0].Code != "a" || got.Diagnostics[0].Message != "first" {
		t.Fatalf("diagnostics = %#v, want normalized Evidence diagnostics", got.Diagnostics)
	}
	if got.Metadata["product"] != "ramen" || len(got.Metadata) != 1 {
		t.Fatalf("metadata = %#v, want trimmed non-empty metadata", got.Metadata)
	}
}

func TestCanonicalJSONStable(t *testing.T) {
	state := State{
		ID:      "s1",
		Answers: []Answer{{Slot: "b", Value: "2"}, {Slot: "a", Value: "1"}},
	}

	first, err := CanonicalJSON(state)
	if err != nil {
		t.Fatalf("CanonicalJSON returned error: %v", err)
	}
	second, err := CanonicalJSON(state)
	if err != nil {
		t.Fatalf("CanonicalJSON returned error: %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("CanonicalJSON changed between calls:\n%s\n%s", first, second)
	}
	if !bytes.Contains(first, []byte(`"version": "authoring.session.v1"`)) {
		t.Fatalf("CanonicalJSON missing session version: %s", first)
	}
}

func TestProductShapedStateMapsWithoutProductImports(t *testing.T) {
	openUdonShaped := struct {
		WorkflowName string
		Goal         string
		Credentials  []string
	}{
		WorkflowName: "daily-report",
		Goal:         "Generate the daily report",
		Credentials:  []string{"reporting-api"},
	}
	ramenShaped := struct {
		ResourceAddress string
		OperationID     string
		ProjectPath     string
	}{
		ResourceAddress: "api_report.daily",
		OperationID:     "createReport",
		ProjectPath:     "project.uws.yaml",
	}

	mapped := Normalize(State{
		Goal: openUdonShaped.Goal,
		Answers: []Answer{
			{Slot: "openudon.workflow.name", Value: openUdonShaped.WorkflowName, Source: "adapter"},
			{Slot: "openudon.credential_binding", Value: openUdonShaped.Credentials[0], Source: "adapter"},
			{Slot: "ramen.resource.address", Value: ramenShaped.ResourceAddress, Source: "adapter"},
			{Slot: "ramen.operation_id", Value: ramenShaped.OperationID, Source: "adapter"},
		},
		Artifacts: []trust.ArtifactRecord{{Path: ramenShaped.ProjectPath, Kind: "project"}},
	})

	if len(mapped.Answers) != 4 {
		t.Fatalf("mapped answers = %d, want 4", len(mapped.Answers))
	}
	if mapped.Artifacts[0].Path != "project.uws.yaml" {
		t.Fatalf("mapped artifact = %#v, want Ramen project artifact summary", mapped.Artifacts[0])
	}
}
