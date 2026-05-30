package session

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/OpenUdon/authoring/decision"
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
	Message         string                   `json:"message,omitempty"`
	SuggestedAnswer string                   `json:"suggested_answer,omitempty"`
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
	state.Version = firstNonEmpty(trim(state.Version), Version)
	state.ID = trim(state.ID)
	state.Goal = trim(state.Goal)
	state.Mode = normalizeToken(state.Mode)
	state.CreatedUTC = trim(state.CreatedUTC)
	state.UpdatedUTC = trim(state.UpdatedUTC)
	state.Turns = normalizeTurns(state.Turns)
	state.Answers = normalizeAnswers(state.Answers)
	state.Readiness = normalizeReadiness(state.Readiness)
	state.Decisions = normalizeDecisions(state.Decisions)
	state.Artifacts = normalizeArtifacts(state.Artifacts)
	state.Diagnostics = trust.NormalizeDiagnostics(state.Diagnostics)
	state.Metadata = normalizeMetadata(state.Metadata)
	return state
}

// CanonicalJSON returns deterministic indented JSON for state.
func CanonicalJSON(state State) ([]byte, error) {
	return json.MarshalIndent(Normalize(state), "", "  ")
}

func normalizeTurns(turns []PromptTurn) []PromptTurn {
	out := make([]PromptTurn, 0, len(turns))
	for _, turn := range turns {
		turn.ID = trim(turn.ID)
		turn.Label = trim(turn.Label)
		turn.Question = trim(turn.Question)
		turn.Answer = trim(turn.Answer)
		turn.Default = trim(turn.Default)
		turn.Source = normalizeToken(turn.Source)
		turn.TimeUTC = trim(turn.TimeUTC)
		if emptyTurn(turn) {
			continue
		}
		out = append(out, turn)
	}
	slices.SortStableFunc(out, func(a, b PromptTurn) int {
		return compareStrings(
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
		answer.Slot = trim(answer.Slot)
		answer.Value = trim(answer.Value)
		answer.Source = normalizeToken(answer.Source)
		answer.TimeUTC = trim(answer.TimeUTC)
		if answer.Slot == "" {
			continue
		}
		out = append(out, answer)
	}
	slices.SortStableFunc(out, func(a, b Answer) int {
		return compareStrings(a.Slot, b.Slot, a.Source, b.Source, a.TimeUTC, b.TimeUTC, a.Value, b.Value)
	})
	return out
}

func normalizeReadiness(issues []ReadinessIssue) []ReadinessIssue {
	out := make([]ReadinessIssue, 0, len(issues))
	for _, issue := range issues {
		issue.Code = normalizeToken(issue.Code)
		issue.Severity = normalizeToken(issue.Severity)
		issue.Slot = trim(issue.Slot)
		issue.Message = trim(issue.Message)
		issue.SuggestedAnswer = trim(issue.SuggestedAnswer)
		if issue.Diagnostic != nil {
			diagnostic := trust.NormalizeDiagnostics([]trust.DiagnosticRecord{*issue.Diagnostic})[0]
			if diagnostic.Code == "" && diagnostic.Message == "" {
				issue.Diagnostic = nil
			} else {
				issue.Diagnostic = &diagnostic
			}
		}
		issue.Diagnostics = trust.NormalizeDiagnostics(issue.Diagnostics)
		if issue.Code == "" && issue.Message == "" && issue.Diagnostic == nil && len(issue.Diagnostics) == 0 {
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

func normalizeArtifacts(artifacts []trust.ArtifactRecord) []trust.ArtifactRecord {
	out := make([]trust.ArtifactRecord, 0, len(artifacts))
	for _, record := range artifacts {
		record.Path = trim(record.Path)
		record.Kind = normalizeToken(record.Kind)
		record.MediaType = trim(record.MediaType)
		record.Classification = normalizeToken(record.Classification)
		record.Digest.Algorithm = normalizeToken(record.Digest.Algorithm)
		record.Digest.Value = trim(record.Digest.Value)
		if record.Path == "" {
			continue
		}
		out = append(out, record)
	}
	slices.SortStableFunc(out, func(a, b trust.ArtifactRecord) int {
		return compareStrings(a.Path, b.Path, a.Kind, b.Kind, a.MediaType, b.MediaType)
	})
	return out
}

func normalizeMetadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		key = trim(key)
		value = trim(value)
		if key == "" || value == "" {
			continue
		}
		out[key] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func compareReadiness(a, b ReadinessIssue) int {
	if diff := compareSeverity(a.Severity, b.Severity); diff != 0 {
		return diff
	}
	return compareStrings(a.Code, b.Code, a.Slot, b.Slot, a.Message, b.Message)
}

func compareSeverity(a, b string) int {
	return severityRank(a) - severityRank(b)
}

func severityRank(severity string) int {
	switch normalizeToken(severity) {
	case "blocking", "error":
		return 0
	case "warning":
		return 1
	case "advisory":
		return 2
	case "info":
		return 3
	default:
		return 4
	}
}

func compareStrings(values ...string) int {
	for i := 0; i+1 < len(values); i += 2 {
		if values[i] < values[i+1] {
			return -1
		}
		if values[i] > values[i+1] {
			return 1
		}
	}
	return 0
}

func emptyTurn(turn PromptTurn) bool {
	return turn.ID == "" && turn.Label == "" && turn.Question == "" && turn.Answer == "" && turn.Default == ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), "_"))
}

func trim(value string) string {
	return strings.TrimSpace(value)
}
