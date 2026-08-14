package interview

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

const (
	// Version is the durable JSON contract version for interview state.
	Version = "authoring.interview.v1"

	StatusOpen         = "open"
	StatusSettled      = "settled"
	StatusDeferred     = "deferred"
	StatusInapplicable = "inapplicable"

	EvidenceObservedFact       = "observed_fact"
	EvidenceUserDecision       = "user_decision"
	EvidenceRecommendation     = "recommendation"
	EvidenceAssumption         = "assumption"
	EvidenceOpenDecision       = "open_decision"
	EvidenceDeferral           = "deferral"
	EvidenceInapplicableBranch = "inapplicable_branch"
)

// State is the product-neutral durable state for one dependency-aware
// interview.
type State struct {
	Version          string            `json:"version"`
	Nodes            []Node            `json:"nodes,omitempty"`
	Evidence         []Evidence        `json:"evidence,omitempty"`
	Answers          []Answer          `json:"answers,omitempty"`
	Deferrals        []Deferral        `json:"deferrals,omitempty"`
	Round            int               `json:"round,omitempty"`
	NoProgressRounds int               `json:"no_progress_rounds,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// Node is one independently answerable decision. Dependencies name decisions
// that must be settled before this node enters the frontier.
type Node struct {
	ID             string   `json:"id"`
	Title          string   `json:"title,omitempty"`
	Prompt         string   `json:"prompt,omitempty"`
	Status         string   `json:"status"`
	Dependencies   []string `json:"dependencies,omitempty"`
	Priority       int      `json:"priority,omitempty"`
	Required       bool     `json:"required,omitempty"`
	Deferrable     bool     `json:"deferrable,omitempty"`
	Rationale      string   `json:"rationale,omitempty"`
	Recommendation string   `json:"recommendation,omitempty"`
	EvidenceRefs   []string `json:"evidence_refs,omitempty"`
}

// Evidence is one concise, reviewable fact or decision in the unified ledger.
// Rationale is public explanation, not hidden chain-of-thought.
type Evidence struct {
	ID         string   `json:"id"`
	Kind       string   `json:"kind"`
	NodeID     string   `json:"node_id,omitempty"`
	Summary    string   `json:"summary"`
	Value      string   `json:"value,omitempty"`
	Source     string   `json:"source,omitempty"`
	References []string `json:"references,omitempty"`
	// Attributes holds concise machine-readable qualifiers needed to apply
	// product-specific safety and readiness policy after resuming. Values are
	// public evidence metadata, never hidden model reasoning.
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Deferral records why a technical leaf is intentionally incomplete and what
// is needed to resume it.
type Deferral struct {
	ID                  string `json:"id"`
	NodeID              string `json:"node_id"`
	Owner               string `json:"owner"`
	Impact              string `json:"impact"`
	UnblockCondition    string `json:"unblock_condition"`
	SuggestedNextAction string `json:"suggested_next_action"`
}

// Answer records one value applied in a frontier round.
type Answer struct {
	ID           string   `json:"id"`
	NodeID       string   `json:"node_id"`
	Value        string   `json:"value"`
	Source       string   `json:"source,omitempty"`
	EvidenceRefs []string `json:"evidence_refs,omitempty"`
}

// Diagnostic describes a graph or state-contract violation.
type Diagnostic struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	NodeID       string `json:"node_id,omitempty"`
	DependencyID string `json:"dependency_id,omitempty"`
}

// ValidationError contains every deterministic validation diagnostic.
type ValidationError struct {
	Diagnostics []Diagnostic
}

func (err ValidationError) Error() string {
	if len(err.Diagnostics) == 0 {
		return "interview state is invalid"
	}
	parts := make([]string, 0, len(err.Diagnostics))
	for _, diagnostic := range err.Diagnostics {
		parts = append(parts, diagnostic.Message)
	}
	return "interview state is invalid: " + strings.Join(parts, "; ")
}

// Normalize returns a deterministic copy of state.
func Normalize(state State) State {
	if strings.TrimSpace(state.Version) == "" {
		state.Version = Version
	} else {
		state.Version = strings.TrimSpace(state.Version)
	}
	if state.Round < 0 {
		state.Round = 0
	}
	if state.NoProgressRounds < 0 {
		state.NoProgressRounds = 0
	}
	state.Nodes = normalizeNodes(state.Nodes)
	state.Evidence = normalizeEvidence(state.Evidence)
	state.Answers = normalizeAnswers(state.Answers)
	state.Deferrals = normalizeDeferrals(state.Deferrals)
	state.Metadata = normalizeMetadata(state.Metadata)
	return state
}

// CanonicalJSON returns deterministic indented JSON for state.
func CanonicalJSON(state State) ([]byte, error) {
	return json.MarshalIndent(Normalize(state), "", "  ")
}

// Validate verifies identifiers, references, cycles, statuses, evidence, and
// deferral completeness.
func Validate(state State) error {
	state = Normalize(state)
	var diagnostics []Diagnostic
	if state.Version != Version {
		diagnostics = append(diagnostics, Diagnostic{Code: "version.unsupported", Message: fmt.Sprintf("unsupported interview version %q; want %q", state.Version, Version)})
	}
	nodes := map[string]Node{}
	for _, node := range state.Nodes {
		if node.ID == "" {
			diagnostics = append(diagnostics, Diagnostic{Code: "node.id.required", Message: "interview node ID is required"})
			continue
		}
		if _, exists := nodes[node.ID]; exists {
			diagnostics = append(diagnostics, Diagnostic{Code: "node.id.duplicate", NodeID: node.ID, Message: fmt.Sprintf("duplicate interview node ID %q", node.ID)})
			continue
		}
		nodes[node.ID] = node
		if !validStatus(node.Status) {
			diagnostics = append(diagnostics, Diagnostic{Code: "node.status.invalid", NodeID: node.ID, Message: fmt.Sprintf("node %q has invalid status %q", node.ID, node.Status)})
		}
	}
	for _, node := range state.Nodes {
		for _, dependency := range node.Dependencies {
			if dependency == node.ID {
				diagnostics = append(diagnostics, Diagnostic{Code: "node.dependency.self", NodeID: node.ID, DependencyID: dependency, Message: fmt.Sprintf("node %q depends on itself", node.ID)})
				continue
			}
			if _, ok := nodes[dependency]; !ok {
				diagnostics = append(diagnostics, Diagnostic{Code: "node.dependency.missing", NodeID: node.ID, DependencyID: dependency, Message: fmt.Sprintf("node %q depends on missing node %q", node.ID, dependency)})
			}
		}
	}
	diagnostics = append(diagnostics, cycleDiagnostics(nodes)...)

	evidenceIDs := map[string]bool{}
	for _, evidence := range state.Evidence {
		if evidence.ID == "" {
			diagnostics = append(diagnostics, Diagnostic{Code: "evidence.id.required", NodeID: evidence.NodeID, Message: "evidence ID is required"})
		} else if evidenceIDs[evidence.ID] {
			diagnostics = append(diagnostics, Diagnostic{Code: "evidence.id.duplicate", NodeID: evidence.NodeID, Message: fmt.Sprintf("duplicate evidence ID %q", evidence.ID)})
		} else {
			evidenceIDs[evidence.ID] = true
		}
		if !validEvidenceKind(evidence.Kind) {
			diagnostics = append(diagnostics, Diagnostic{Code: "evidence.kind.invalid", NodeID: evidence.NodeID, Message: fmt.Sprintf("evidence %q has invalid kind %q", evidence.ID, evidence.Kind)})
		}
		if evidence.NodeID != "" {
			if _, ok := nodes[evidence.NodeID]; !ok {
				diagnostics = append(diagnostics, Diagnostic{Code: "evidence.node.missing", NodeID: evidence.NodeID, Message: fmt.Sprintf("evidence %q references missing node %q", evidence.ID, evidence.NodeID)})
			}
		}
	}
	for _, node := range state.Nodes {
		for _, ref := range node.EvidenceRefs {
			if !evidenceIDs[ref] {
				diagnostics = append(diagnostics, Diagnostic{Code: "node.evidence.missing", NodeID: node.ID, Message: fmt.Sprintf("node %q references missing evidence %q", node.ID, ref)})
			}
		}
	}

	answerIDs := map[string]bool{}
	for _, answer := range state.Answers {
		if answer.ID == "" {
			diagnostics = append(diagnostics, Diagnostic{Code: "answer.id.required", NodeID: answer.NodeID, Message: "answer ID is required"})
		} else if answerIDs[answer.ID] {
			diagnostics = append(diagnostics, Diagnostic{Code: "answer.id.duplicate", NodeID: answer.NodeID, Message: fmt.Sprintf("duplicate answer ID %q", answer.ID)})
		} else {
			answerIDs[answer.ID] = true
		}
		if _, ok := nodes[answer.NodeID]; !ok {
			diagnostics = append(diagnostics, Diagnostic{Code: "answer.node.missing", NodeID: answer.NodeID, Message: fmt.Sprintf("answer %q references missing node %q", answer.ID, answer.NodeID)})
		}
		for _, ref := range answer.EvidenceRefs {
			if !evidenceIDs[ref] {
				diagnostics = append(diagnostics, Diagnostic{Code: "answer.evidence.missing", NodeID: answer.NodeID, Message: fmt.Sprintf("answer %q references missing evidence %q", answer.ID, ref)})
			}
		}
	}

	deferralNodes := map[string]bool{}
	deferralIDs := map[string]bool{}
	for _, deferral := range state.Deferrals {
		if deferral.ID == "" {
			diagnostics = append(diagnostics, Diagnostic{Code: "deferral.id.required", NodeID: deferral.NodeID, Message: "deferral ID is required"})
		} else if deferralIDs[deferral.ID] {
			diagnostics = append(diagnostics, Diagnostic{Code: "deferral.id.duplicate", NodeID: deferral.NodeID, Message: fmt.Sprintf("duplicate deferral ID %q", deferral.ID)})
		} else {
			deferralIDs[deferral.ID] = true
		}
		node, ok := nodes[deferral.NodeID]
		if !ok {
			diagnostics = append(diagnostics, Diagnostic{Code: "deferral.node.missing", NodeID: deferral.NodeID, Message: fmt.Sprintf("deferral %q references missing node %q", deferral.ID, deferral.NodeID)})
		} else {
			if node.Status != StatusDeferred {
				diagnostics = append(diagnostics, Diagnostic{Code: "deferral.status.invalid", NodeID: deferral.NodeID, Message: fmt.Sprintf("deferral %q belongs to node %q with status %q", deferral.ID, deferral.NodeID, node.Status)})
			}
			if !node.Deferrable {
				diagnostics = append(diagnostics, Diagnostic{Code: "deferral.node.not_deferrable", NodeID: deferral.NodeID, Message: fmt.Sprintf("node %q is not deferrable", deferral.NodeID)})
			}
		}
		if deferralNodes[deferral.NodeID] {
			diagnostics = append(diagnostics, Diagnostic{Code: "deferral.node.duplicate", NodeID: deferral.NodeID, Message: fmt.Sprintf("node %q has multiple active deferrals", deferral.NodeID)})
		}
		deferralNodes[deferral.NodeID] = true
		if deferral.Owner == "" || deferral.Impact == "" || deferral.UnblockCondition == "" || deferral.SuggestedNextAction == "" {
			diagnostics = append(diagnostics, Diagnostic{Code: "deferral.incomplete", NodeID: deferral.NodeID, Message: fmt.Sprintf("deferral %q must record owner, impact, unblock condition, and suggested next action", deferral.ID)})
		}
	}
	for _, node := range state.Nodes {
		if node.Status == StatusDeferred && !deferralNodes[node.ID] {
			diagnostics = append(diagnostics, Diagnostic{Code: "deferral.missing", NodeID: node.ID, Message: fmt.Sprintf("deferred node %q has no deferral record", node.ID)})
		}
	}

	slices.SortStableFunc(diagnostics, compareDiagnostic)
	if len(diagnostics) > 0 {
		return ValidationError{Diagnostics: diagnostics}
	}
	return nil
}

// Frontier returns every open node whose prerequisites are settled, ordered by
// descending priority and then stable node ID.
func Frontier(state State) ([]Node, error) {
	state = Normalize(state)
	if err := Validate(state); err != nil {
		return nil, err
	}
	statuses := make(map[string]string, len(state.Nodes))
	for _, node := range state.Nodes {
		statuses[node.ID] = node.Status
	}
	var frontier []Node
	for _, node := range state.Nodes {
		if node.Status != StatusOpen {
			continue
		}
		ready := true
		for _, dependency := range node.Dependencies {
			if statuses[dependency] != StatusSettled {
				ready = false
				break
			}
		}
		if ready {
			frontier = append(frontier, node)
		}
	}
	slices.SortStableFunc(frontier, func(a, b Node) int {
		if a.Priority != b.Priority {
			return b.Priority - a.Priority
		}
		return strings.Compare(a.ID, b.ID)
	})
	return frontier, nil
}

// ValidateTransition rejects status changes that would rewrite settled or
// inapplicable decisions. Deferred nodes may be reopened after their unblock
// condition is met.
func ValidateTransition(from, to string) error {
	from = normalizeStatus(from)
	to = normalizeStatus(to)
	if !validStatus(from) || !validStatus(to) {
		return fmt.Errorf("invalid interview status transition %q -> %q", from, to)
	}
	if from == to {
		return nil
	}
	valid := from == StatusOpen && (to == StatusSettled || to == StatusDeferred || to == StatusInapplicable)
	valid = valid || (from == StatusDeferred && to == StatusOpen)
	if !valid {
		return fmt.Errorf("invalid interview status transition %q -> %q", from, to)
	}
	return nil
}

// Transition returns state with one validated node status change.
func Transition(state State, nodeID, status string) (State, error) {
	state = Normalize(state)
	if err := Validate(state); err != nil {
		return state, err
	}
	index := slices.IndexFunc(state.Nodes, func(node Node) bool { return node.ID == strings.TrimSpace(nodeID) })
	if index < 0 {
		return state, fmt.Errorf("interview node %q does not exist", nodeID)
	}
	status = normalizeStatus(status)
	if err := ValidateTransition(state.Nodes[index].Status, status); err != nil {
		return state, err
	}
	if status == StatusDeferred {
		return state, fmt.Errorf("use ApplyDeferral to defer node %q with complete unblock metadata", state.Nodes[index].ID)
	}
	state.Nodes[index].Status = status
	if status == StatusOpen {
		state.Deferrals = slices.DeleteFunc(state.Deferrals, func(deferral Deferral) bool { return deferral.NodeID == state.Nodes[index].ID })
	}
	return Normalize(state), nil
}

// ApplyAnswer settles a ready frontier node and appends its durable answer.
func ApplyAnswer(state State, answer Answer) (State, error) {
	return ApplyAnswers(state, []Answer{answer})
}

// ApplyAnswers atomically settles a full set of independent answers from the
// current frontier and advances the interview by one round.
func ApplyAnswers(state State, answers []Answer) (State, error) {
	state = Normalize(state)
	if len(answers) == 0 {
		return state, fmt.Errorf("at least one frontier answer is required")
	}
	frontier, err := Frontier(state)
	if err != nil {
		return state, err
	}
	frontierIDs := map[string]bool{}
	for _, node := range frontier {
		frontierIDs[node.ID] = true
	}
	existingIDs := map[string]bool{}
	for _, existing := range state.Answers {
		existingIDs[existing.ID] = true
	}
	seenNodes := map[string]bool{}
	normalized := make([]Answer, 0, len(answers))
	for _, answer := range answers {
		answer = normalizeAnswer(answer)
		if answer.ID == "" || answer.NodeID == "" || answer.Value == "" {
			return state, fmt.Errorf("answer ID, node ID, and value are required")
		}
		if !frontierIDs[answer.NodeID] {
			return state, fmt.Errorf("node %q is not in the current frontier", answer.NodeID)
		}
		if seenNodes[answer.NodeID] {
			return state, fmt.Errorf("frontier round contains multiple answers for node %q", answer.NodeID)
		}
		if existingIDs[answer.ID] {
			return state, fmt.Errorf("duplicate answer ID %q", answer.ID)
		}
		seenNodes[answer.NodeID] = true
		existingIDs[answer.ID] = true
		normalized = append(normalized, answer)
	}
	state.Answers = append(state.Answers, normalized...)
	for i := range state.Nodes {
		if seenNodes[state.Nodes[i].ID] {
			state.Nodes[i].Status = StatusSettled
		}
	}
	state.Round++
	state.NoProgressRounds = 0
	state = Normalize(state)
	if err := Validate(state); err != nil {
		return state, err
	}
	return state, nil
}

// ApplyDeferral defers one ready technical leaf with complete unblock data.
func ApplyDeferral(state State, deferral Deferral) (State, error) {
	state = Normalize(state)
	deferral = normalizeDeferral(deferral)
	frontier, err := Frontier(state)
	if err != nil {
		return state, err
	}
	nodeIndex := slices.IndexFunc(frontier, func(node Node) bool { return node.ID == deferral.NodeID })
	if nodeIndex < 0 {
		return state, fmt.Errorf("node %q is not in the current frontier", deferral.NodeID)
	}
	if !frontier[nodeIndex].Deferrable {
		return state, fmt.Errorf("node %q is not deferrable", deferral.NodeID)
	}
	if deferral.ID == "" || deferral.Owner == "" || deferral.Impact == "" || deferral.UnblockCondition == "" || deferral.SuggestedNextAction == "" {
		return state, fmt.Errorf("deferral ID, owner, impact, unblock condition, and suggested next action are required")
	}
	for i := range state.Nodes {
		if state.Nodes[i].ID == deferral.NodeID {
			state.Nodes[i].Status = StatusDeferred
			break
		}
	}
	state.Deferrals = append(state.Deferrals, deferral)
	state.Round++
	state.NoProgressRounds = 0
	state = Normalize(state)
	if err := Validate(state); err != nil {
		return state, err
	}
	return state, nil
}

func normalizeNodes(nodes []Node) []Node {
	out := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		node.ID = strings.TrimSpace(node.ID)
		node.Title = strings.TrimSpace(node.Title)
		node.Prompt = strings.TrimSpace(node.Prompt)
		node.Status = normalizeStatus(node.Status)
		node.Dependencies = normalizedStrings(node.Dependencies)
		node.Rationale = strings.TrimSpace(node.Rationale)
		node.Recommendation = strings.TrimSpace(node.Recommendation)
		node.EvidenceRefs = normalizedStrings(node.EvidenceRefs)
		out = append(out, node)
	}
	slices.SortStableFunc(out, func(a, b Node) int { return strings.Compare(a.ID, b.ID) })
	return out
}

func normalizeEvidence(records []Evidence) []Evidence {
	out := make([]Evidence, 0, len(records))
	for _, record := range records {
		record.ID = strings.TrimSpace(record.ID)
		record.Kind = strings.ToLower(strings.TrimSpace(record.Kind))
		record.NodeID = strings.TrimSpace(record.NodeID)
		record.Summary = strings.TrimSpace(record.Summary)
		record.Value = strings.TrimSpace(record.Value)
		record.Source = strings.TrimSpace(record.Source)
		record.References = normalizedStrings(record.References)
		record.Attributes = normalizeMetadata(record.Attributes)
		out = append(out, record)
	}
	slices.SortStableFunc(out, func(a, b Evidence) int { return strings.Compare(a.ID, b.ID) })
	return out
}

func normalizeAnswers(answers []Answer) []Answer {
	out := make([]Answer, 0, len(answers))
	for _, answer := range answers {
		out = append(out, normalizeAnswer(answer))
	}
	slices.SortStableFunc(out, func(a, b Answer) int { return strings.Compare(a.ID, b.ID) })
	return out
}

func normalizeAnswer(answer Answer) Answer {
	answer.ID = strings.TrimSpace(answer.ID)
	answer.NodeID = strings.TrimSpace(answer.NodeID)
	answer.Value = strings.TrimSpace(answer.Value)
	answer.Source = strings.TrimSpace(answer.Source)
	answer.EvidenceRefs = normalizedStrings(answer.EvidenceRefs)
	return answer
}

func normalizeDeferrals(deferrals []Deferral) []Deferral {
	out := make([]Deferral, 0, len(deferrals))
	for _, deferral := range deferrals {
		out = append(out, normalizeDeferral(deferral))
	}
	slices.SortStableFunc(out, func(a, b Deferral) int { return strings.Compare(a.ID, b.ID) })
	return out
}

func normalizeDeferral(deferral Deferral) Deferral {
	deferral.ID = strings.TrimSpace(deferral.ID)
	deferral.NodeID = strings.TrimSpace(deferral.NodeID)
	deferral.Owner = strings.TrimSpace(deferral.Owner)
	deferral.Impact = strings.TrimSpace(deferral.Impact)
	deferral.UnblockCondition = strings.TrimSpace(deferral.UnblockCondition)
	deferral.SuggestedNextAction = strings.TrimSpace(deferral.SuggestedNextAction)
	return deferral
}

func normalizeMetadata(metadata map[string]string) map[string]string {
	if len(metadata) == 0 {
		return nil
	}
	out := map[string]string{}
	for key, value := range metadata {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" && value != "" {
			out[key] = value
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func normalizedStrings(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	slices.Sort(out)
	return out
}

func normalizeStatus(status string) string {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		return StatusOpen
	}
	return status
}

func validStatus(status string) bool {
	switch normalizeStatus(status) {
	case StatusOpen, StatusSettled, StatusDeferred, StatusInapplicable:
		return true
	default:
		return false
	}
}

func validEvidenceKind(kind string) bool {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case EvidenceObservedFact, EvidenceUserDecision, EvidenceRecommendation, EvidenceAssumption, EvidenceOpenDecision, EvidenceDeferral, EvidenceInapplicableBranch:
		return true
	default:
		return false
	}
}

func cycleDiagnostics(nodes map[string]Node) []Diagnostic {
	const (
		unvisited = iota
		visiting
		visited
	)
	marks := map[string]int{}
	var diagnostics []Diagnostic
	var visit func(string)
	visit = func(id string) {
		if marks[id] == visited {
			return
		}
		if marks[id] == visiting {
			diagnostics = append(diagnostics, Diagnostic{Code: "node.dependency.cycle", NodeID: id, Message: fmt.Sprintf("dependency cycle includes node %q", id)})
			return
		}
		marks[id] = visiting
		for _, dependency := range nodes[id].Dependencies {
			if _, ok := nodes[dependency]; ok {
				visit(dependency)
			}
		}
		marks[id] = visited
	}
	ids := make([]string, 0, len(nodes))
	for id := range nodes {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		visit(id)
	}
	return diagnostics
}

func compareDiagnostic(a, b Diagnostic) int {
	if value := strings.Compare(a.Code, b.Code); value != 0 {
		return value
	}
	if value := strings.Compare(a.NodeID, b.NodeID); value != 0 {
		return value
	}
	if value := strings.Compare(a.DependencyID, b.DependencyID); value != 0 {
		return value
	}
	return strings.Compare(a.Message, b.Message)
}
