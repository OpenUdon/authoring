package interview

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestFrontierReturnsEveryReadyNodeDeterministically(t *testing.T) {
	state := State{Nodes: []Node{
		{ID: "mapping", Dependencies: []string{"operation"}, Priority: 5},
		{ID: "trigger", Priority: 10},
		{ID: "output", Dependencies: []string{"operation"}, Priority: 5},
		{ID: "operation", Dependencies: []string{"source"}, Priority: 8},
		{ID: "source", Priority: 10},
	}}
	frontier, err := Frontier(state)
	if err != nil {
		t.Fatal(err)
	}
	if got := nodeIDs(frontier); !reflect.DeepEqual(got, []string{"source", "trigger"}) {
		t.Fatalf("frontier = %v", got)
	}

	state, err = ApplyAnswers(state, []Answer{
		{ID: "a-source", NodeID: "source", Value: "openapi/pets.yaml", Source: "user"},
		{ID: "a-trigger", NodeID: "trigger", Value: "manual", Source: "user"},
	})
	if err != nil {
		t.Fatal(err)
	}
	frontier, err = Frontier(state)
	if err != nil {
		t.Fatal(err)
	}
	if got := nodeIDs(frontier); !reflect.DeepEqual(got, []string{"operation"}) {
		t.Fatalf("second frontier = %v", got)
	}
	if state.Round != 1 {
		t.Fatalf("round = %d, want one atomic round", state.Round)
	}
}

func TestValidateDiagnosesDuplicateMissingDependencyAndCycle(t *testing.T) {
	err := Validate(State{Nodes: []Node{
		{ID: "a", Dependencies: []string{"b"}},
		{ID: "a"},
		{ID: "b", Dependencies: []string{"a", "missing"}},
	}})
	var validation ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("Validate() error = %v", err)
	}
	codes := diagnosticCodes(validation.Diagnostics)
	for _, code := range []string{"node.id.duplicate", "node.dependency.missing", "node.dependency.cycle"} {
		if !strings.Contains(codes, code) {
			t.Fatalf("diagnostics = %#v, missing %s", validation.Diagnostics, code)
		}
	}
}

func TestDeferralRequiresCompleteRecordAndDeferrableNode(t *testing.T) {
	state := State{Nodes: []Node{{ID: "mapping", Deferrable: true}}}
	if _, err := ApplyDeferral(state, Deferral{ID: "d1", NodeID: "mapping"}); err == nil {
		t.Fatal("ApplyDeferral accepted an incomplete deferral")
	}
	got, err := ApplyDeferral(state, Deferral{
		ID: "d1", NodeID: "mapping", Owner: "operator", Impact: "intent remains a draft",
		UnblockCondition: "operation schema is available", SuggestedNextAction: "provide the source document",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Nodes[0].Status != StatusDeferred || len(got.Deferrals) != 1 {
		t.Fatalf("state = %#v", got)
	}
	if err := Validate(State{Nodes: []Node{{ID: "boundary", Status: StatusDeferred}}}); err == nil {
		t.Fatal("Validate accepted deferred non-deferrable node without record")
	}
}

func TestInvalidStatusTransitionsAreRejected(t *testing.T) {
	for _, transition := range [][2]string{{StatusSettled, StatusOpen}, {StatusSettled, StatusDeferred}, {StatusInapplicable, StatusOpen}, {StatusDeferred, StatusSettled}} {
		if err := ValidateTransition(transition[0], transition[1]); err == nil {
			t.Fatalf("transition %s -> %s accepted", transition[0], transition[1])
		}
	}
	for _, transition := range [][2]string{{StatusOpen, StatusSettled}, {StatusOpen, StatusDeferred}, {StatusOpen, StatusInapplicable}, {StatusDeferred, StatusOpen}} {
		if err := ValidateTransition(transition[0], transition[1]); err != nil {
			t.Fatalf("transition %s -> %s rejected: %v", transition[0], transition[1], err)
		}
	}
}

func TestUnifiedEvidenceReferencesAreValidated(t *testing.T) {
	state := State{
		Nodes:    []Node{{ID: "source", EvidenceRefs: []string{"observed"}}},
		Evidence: []Evidence{{ID: "observed", Kind: EvidenceObservedFact, NodeID: "source", Summary: "source file exists"}},
	}
	if err := Validate(state); err != nil {
		t.Fatal(err)
	}
	state.Evidence = append(state.Evidence, Evidence{ID: "observed", Kind: EvidenceAssumption, Summary: "duplicate"})
	if err := Validate(state); err == nil {
		t.Fatal("Validate accepted duplicate evidence ID")
	}
}

func TestUnifiedEvidenceAttributesAreNormalizedAndDurable(t *testing.T) {
	state := Normalize(State{Evidence: []Evidence{{
		ID: " safety ", Kind: EvidenceOpenDecision, Summary: " confirm send ",
		Attributes: map[string]string{" confidence ": " review ", "empty": "  "},
	}}})
	if got := state.Evidence[0].Attributes; len(got) != 1 || got["confidence"] != "review" {
		t.Fatalf("normalized attributes = %#v", got)
	}
	data, err := CanonicalJSON(state)
	if err != nil {
		t.Fatal(err)
	}
	var roundTrip State
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if got := roundTrip.Evidence[0].Attributes["confidence"]; got != "review" {
		t.Fatalf("round-trip confidence = %q", got)
	}
}

func nodeIDs(nodes []Node) []string {
	out := make([]string, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, node.ID)
	}
	return out
}

func diagnosticCodes(diagnostics []Diagnostic) string {
	var out []string
	for _, diagnostic := range diagnostics {
		out = append(out, diagnostic.Code)
	}
	return strings.Join(out, ",")
}
