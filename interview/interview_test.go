package interview

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestApplyRoundSettlesMixedAnswersAndDeferralsAtomically(t *testing.T) {
	state := State{NoProgressRounds: 2, Nodes: []Node{{ID: "answer", Priority: 2}, {ID: "defer", Priority: 1, Deferrable: true}}}
	answer := Answer{ID: "a1", NodeID: "answer", Value: "chosen", Source: "user", EvidenceRefs: []string{"e1"}}
	deferral := Deferral{ID: "d1", NodeID: "defer", Owner: "api owner", Impact: "mapping blocked", UnblockCondition: "schema published", SuggestedNextAction: "attach source"}
	got, err := ApplyRound(state, []Resolution{
		{NodeID: "answer", Answer: &answer, Evidence: []Evidence{{ID: "e1", Kind: EvidenceUserDecision, NodeID: "answer", Summary: "operator chose value"}}},
		{NodeID: "defer", Deferral: &deferral, Evidence: []Evidence{{ID: "e2", Kind: EvidenceDeferral, NodeID: "defer", Summary: "schema unavailable"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Round != 1 || got.NoProgressRounds != 0 || len(got.Answers) != 1 || len(got.Deferrals) != 1 || len(got.Evidence) != 2 {
		t.Fatalf("settled state = %#v", got)
	}
	statuses := map[string]string{}
	for _, node := range got.Nodes {
		statuses[node.ID] = node.Status
	}
	if statuses["answer"] != StatusSettled || statuses["defer"] != StatusDeferred {
		t.Fatalf("statuses = %#v", statuses)
	}
}

func TestApplyRoundRejectsIncompleteDuplicateAndNonFrontierWithoutMutation(t *testing.T) {
	base := Normalize(State{Nodes: []Node{{ID: "a"}, {ID: "b"}, {ID: "later", Dependencies: []string{"a"}}}})
	before, err := CanonicalJSON(base)
	if err != nil {
		t.Fatal(err)
	}
	a := Answer{ID: "a1", NodeID: "a", Value: "one"}
	b := Answer{ID: "b1", NodeID: "b", Value: "two"}
	later := Answer{ID: "l1", NodeID: "later", Value: "early"}
	for _, resolutions := range [][]Resolution{
		{{NodeID: "a", Answer: &a}},
		{{NodeID: "a", Answer: &a}, {NodeID: "a", Answer: &b}},
		{{NodeID: "a", Answer: &a}, {NodeID: "later", Answer: &later}},
		{{NodeID: "a", Answer: &a}, {NodeID: "b", Answer: &b, Evidence: []Evidence{{ID: "bad", Kind: "invalid", Summary: "bad"}}}},
	} {
		got, err := ApplyRound(base, resolutions)
		if err == nil {
			t.Fatalf("ApplyRound accepted %#v", resolutions)
		}
		after, marshalErr := CanonicalJSON(got)
		if marshalErr != nil || !reflect.DeepEqual(before, after) {
			t.Fatalf("failed round changed state: %s != %s (marshal %v)", before, after, marshalErr)
		}
	}
}

func TestApplyRoundSettlesWideFrontier(t *testing.T) {
	const count = 20000
	nodes := make([]Node, count)
	resolutions := make([]Resolution, count)
	for index := range nodes {
		nodeID := fmt.Sprintf("node-%05d", index)
		nodes[index] = Node{ID: nodeID}
		answer := Answer{ID: fmt.Sprintf("answer-%05d", index), NodeID: nodeID, Value: "accepted"}
		resolutions[index] = Resolution{NodeID: nodeID, Answer: &answer}
	}
	got, err := ApplyRound(State{Nodes: nodes}, resolutions)
	if err != nil {
		t.Fatal(err)
	}
	if got.Round != 1 || len(got.Answers) != count {
		t.Fatalf("wide settlement round=%d answers=%d", got.Round, len(got.Answers))
	}
	for _, node := range got.Nodes {
		if node.Status != StatusSettled {
			t.Fatalf("node %q status = %q", node.ID, node.Status)
		}
	}
}

func TestValidateDeepDependencyGraphIteratively(t *testing.T) {
	const count = 20000
	nodes := make([]Node, count)
	for index := range nodes {
		nodes[index] = Node{ID: fmt.Sprintf("node-%05d", index)}
		if index > 0 {
			nodes[index].Dependencies = []string{nodes[index-1].ID}
		}
	}
	if err := Validate(State{Nodes: nodes}); err != nil {
		t.Fatalf("deep acyclic graph rejected: %v", err)
	}
	nodes[0].Dependencies = []string{nodes[len(nodes)-1].ID}
	if err := Validate(State{Nodes: nodes}); err == nil {
		t.Fatal("deep cycle was not detected")
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
