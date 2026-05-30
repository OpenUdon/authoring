package session

import (
	"encoding/json"
	"slices"

	"github.com/OpenUdon/authoring/decision"
	"github.com/OpenUdon/authoring/internal/norm"
	"github.com/OpenUdon/authoring/internal/records"
	"github.com/OpenUdon/authoring/trust"
)

const (
	// Version is the durable JSON contract version for Authoring sessions.
	Version = "authoring.session.v1"
)

// State is the product-neutral state captured during an authoring session.
type State struct {
	Version     string                   `json:"version"`
	ID          string                   `json:"id,omitempty"`
	Goal        string                   `json:"goal,omitempty"`
	Mode        string                   `json:"mode,omitempty"`
	CreatedUTC  string                   `json:"created_utc,omitempty"`
	UpdatedUTC  string                   `json:"updated_utc,omitempty"`
	Turns       []PromptTurn             `json:"turns,omitempty"`
	Answers     []Answer                 `json:"answers,omitempty"`
	Readiness   []ReadinessIssue         `json:"readiness,omitempty"`
	Decisions   []DecisionEvidence       `json:"decisions,omitempty"`
	Artifacts   []trust.ArtifactRecord   `json:"artifacts,omitempty"`
	Diagnostics []trust.DiagnosticRecord `json:"diagnostics,omitempty"`
	Metadata    map[string]string        `json:"metadata,omitempty"`
}

// PromptTurn records one prompt/answer exchange without owning prompt text
// policy.
type PromptTurn struct {
	ID        string `json:"id,omitempty"`
	Label     string `json:"label,omitempty"`
	Question  string `json:"question,omitempty"`
	Answer    string `json:"answer,omitempty"`
	Default   string `json:"default,omitempty"`
	Source    string `json:"source,omitempty"`
	TimeUTC   string `json:"time_utc,omitempty"`
	Required  bool   `json:"required,omitempty"`
	Redacted  bool   `json:"redacted,omitempty"`
	Sensitive bool   `json:"sensitive,omitempty"`
}

// Answer records a slot value collected or defaulted during authoring.
type Answer struct {
	Slot      string `json:"slot"`
	Value     string `json:"value,omitempty"`
	Source    string `json:"source,omitempty"`
	TimeUTC   string `json:"time_utc,omitempty"`
	Redacted  bool   `json:"redacted,omitempty"`
	Sensitive bool   `json:"sensitive,omitempty"`
}

// ReadinessIssue records a product-neutral readiness finding.
type ReadinessIssue struct {
	Code            string                   `json:"code"`
	Severity        string                   `json:"severity,omitempty"`
	Slot            string                   `json:"slot,omitempty"`
	OperationID     string                   `json:"operation_id,omitempty"`
	Path            string                   `json:"path,omitempty"`
	Message         string                   `json:"message,omitempty"`
	SuggestedAnswer string                   `json:"suggested_answer,omitempty"`
	Remediation     string                   `json:"remediation,omitempty"`
	Diagnostic      *trust.DiagnosticRecord  `json:"diagnostic,omitempty"`
	Diagnostics     []trust.DiagnosticRecord `json:"diagnostics,omitempty"`
}

// DecisionEvidence records why a slot value was chosen or still requires
// confirmation.
type DecisionEvidence = decision.Record

// DecisionAlternative records a rejected or lower-confidence value.
type DecisionAlternative = decision.Alternative

// Normalize returns a deterministic copy of state.
func Normalize(state State) State {
	state.Version = norm.FirstNonEmpty(state.Version, Version)
	state.ID = norm.Trim(state.ID)
	state.Goal = norm.Trim(state.Goal)
	state.Mode = norm.Token(state.Mode)
	state.CreatedUTC = norm.Trim(state.CreatedUTC)
	state.UpdatedUTC = norm.Trim(state.UpdatedUTC)
	state.Turns = normalizeTurns(state.Turns)
	state.Answers = normalizeAnswers(state.Answers)
	state.Readiness = normalizeReadiness(state.Readiness)
	state.Decisions = normalizeDecisions(state.Decisions)
	state.Artifacts = records.Artifacts(state.Artifacts)
	state.Diagnostics = trust.NormalizeDiagnostics(state.Diagnostics)
	state.Metadata = norm.Metadata(state.Metadata)
	return state
}

// CanonicalJSON returns deterministic indented JSON for state.
func CanonicalJSON(state State) ([]byte, error) {
	return json.MarshalIndent(Normalize(state), "", "  ")
}

func normalizeTurns(turns []PromptTurn) []PromptTurn {
	out := make([]PromptTurn, 0, len(turns))
	for _, turn := range turns {
		turn.ID = norm.Trim(turn.ID)
		turn.Label = norm.Trim(turn.Label)
		turn.Question = norm.Trim(turn.Question)
		turn.Answer = norm.Trim(turn.Answer)
		turn.Default = norm.Trim(turn.Default)
		turn.Source = norm.Token(turn.Source)
		turn.TimeUTC = norm.Trim(turn.TimeUTC)
		if emptyTurn(turn) {
			continue
		}
		out = append(out, turn)
	}
	slices.SortStableFunc(out, func(a, b PromptTurn) int {
		return norm.CompareStrings(
			a.TimeUTC, b.TimeUTC,
			a.ID, b.ID,
			a.Label, b.Label,
			a.Source, b.Source,
			a.Question, b.Question,
		)
	})
	return out
}

func normalizeAnswers(answers []Answer) []Answer {
	out := make([]Answer, 0, len(answers))
	for _, answer := range answers {
		answer.Slot = norm.Trim(answer.Slot)
		answer.Value = norm.Trim(answer.Value)
		answer.Source = norm.Token(answer.Source)
		answer.TimeUTC = norm.Trim(answer.TimeUTC)
		if answer.Slot == "" {
			continue
		}
		out = append(out, answer)
	}
	slices.SortStableFunc(out, func(a, b Answer) int {
		return norm.CompareStrings(a.Slot, b.Slot, a.Source, b.Source, a.TimeUTC, b.TimeUTC, a.Value, b.Value)
	})
	return out
}

func normalizeReadiness(issues []ReadinessIssue) []ReadinessIssue {
	out := make([]ReadinessIssue, 0, len(issues))
	for _, issue := range issues {
		issue.Code = norm.Token(issue.Code)
		issue.Severity = norm.Token(issue.Severity)
		issue.Slot = norm.Trim(issue.Slot)
		issue.OperationID = norm.Trim(issue.OperationID)
		issue.Path = norm.Trim(issue.Path)
		issue.Message = norm.Trim(issue.Message)
		issue.SuggestedAnswer = norm.Trim(issue.SuggestedAnswer)
		issue.Remediation = norm.Trim(issue.Remediation)
		if issue.Diagnostic != nil {
			diagnostic := trust.NormalizeDiagnostics([]trust.DiagnosticRecord{*issue.Diagnostic})[0]
			if diagnostic.Code == "" && diagnostic.Message == "" {
				issue.Diagnostic = nil
			} else {
				issue.Diagnostic = &diagnostic
			}
		}
		issue.Diagnostics = trust.NormalizeDiagnostics(issue.Diagnostics)
		if issue.Code == "" && issue.Message == "" && issue.OperationID == "" && issue.Path == "" && issue.Remediation == "" && issue.Diagnostic == nil && len(issue.Diagnostics) == 0 {
			continue
		}
		out = append(out, issue)
	}
	slices.SortStableFunc(out, compareReadiness)
	return out
}

func normalizeDecisions(decisions []DecisionEvidence) []DecisionEvidence {
	return decision.NormalizeAll(decisions)
}

func compareReadiness(a, b ReadinessIssue) int {
	if diff := norm.CompareSeverity(a.Severity, b.Severity); diff != 0 {
		return diff
	}
	return norm.CompareStrings(a.Code, b.Code, a.Slot, b.Slot, a.OperationID, b.OperationID, a.Path, b.Path, a.Message, b.Message)
}

func emptyTurn(turn PromptTurn) bool {
	return turn.ID == "" && turn.Label == "" && turn.Question == "" && turn.Answer == "" && turn.Default == ""
}
