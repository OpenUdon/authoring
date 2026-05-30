package transcript

import (
	"encoding/json"
	"slices"

	"github.com/OpenUdon/authoring/decision"
	"github.com/OpenUdon/authoring/internal/norm"
	"github.com/OpenUdon/authoring/internal/records"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/trust"
)

const (
	// Version is the durable JSON contract version for Authoring transcripts.
	Version = "authoring.transcript.v1"
)

// Record is the durable product-neutral transcript for one authoring run.
type Record struct {
	Version     string                   `json:"version"`
	SessionID   string                   `json:"session_id,omitempty"`
	TimeUTC     string                   `json:"time_utc,omitempty"`
	Provider    *ModelProvenance         `json:"provider,omitempty"`
	Turns       []Turn                   `json:"turns,omitempty"`
	Events      []Event                  `json:"events,omitempty"`
	Diagnostics []trust.DiagnosticRecord `json:"diagnostics,omitempty"`
	Artifacts   []trust.ArtifactRecord   `json:"artifacts,omitempty"`
	Metadata    map[string]string        `json:"metadata,omitempty"`
}

// ModelProvenance records downstream-supplied model/provider metadata without
// defining a provider client.
type ModelProvenance struct {
	Provider   string `json:"provider,omitempty"`
	Model      string `json:"model,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	ResponseID string `json:"response_id,omitempty"`
	Seed       string `json:"seed,omitempty"`
}

// Turn records one transcript turn.
type Turn struct {
	ID          string                     `json:"id,omitempty"`
	Role        string                     `json:"role,omitempty"`
	Label       string                     `json:"label,omitempty"`
	Content     string                     `json:"content,omitempty"`
	TimeUTC     string                     `json:"time_utc,omitempty"`
	Provider    *ModelProvenance           `json:"provider,omitempty"`
	Decisions   []session.DecisionEvidence `json:"decisions,omitempty"`
	Diagnostics []trust.DiagnosticRecord   `json:"diagnostics,omitempty"`
	Redacted    bool                       `json:"redacted,omitempty"`
}

// Event records a structured transcript event.
type Event struct {
	ID          string                   `json:"id,omitempty"`
	Type        string                   `json:"type"`
	Stage       string                   `json:"stage,omitempty"`
	Severity    string                   `json:"severity,omitempty"`
	Message     string                   `json:"message,omitempty"`
	TimeUTC     string                   `json:"time_utc,omitempty"`
	Fields      map[string]string        `json:"fields,omitempty"`
	Diagnostics []trust.DiagnosticRecord `json:"diagnostics,omitempty"`
}

// Normalize returns a deterministic copy of record.
func Normalize(record Record) Record {
	record.Version = norm.FirstNonEmpty(record.Version, Version)
	record.SessionID = norm.Trim(record.SessionID)
	record.TimeUTC = norm.Trim(record.TimeUTC)
	record.Provider = normalizeProvider(record.Provider)
	record.Turns = normalizeTurns(record.Turns)
	record.Events = normalizeEvents(record.Events)
	record.Diagnostics = trust.NormalizeDiagnostics(record.Diagnostics)
	record.Artifacts = records.Artifacts(record.Artifacts)
	record.Metadata = norm.Metadata(record.Metadata)
	return record
}

// CanonicalJSON returns deterministic indented JSON for record.
func CanonicalJSON(record Record) ([]byte, error) {
	return json.MarshalIndent(Normalize(record), "", "  ")
}

func normalizeTurns(turns []Turn) []Turn {
	out := make([]Turn, 0, len(turns))
	for _, turn := range turns {
		turn.ID = norm.Trim(turn.ID)
		turn.Role = norm.Token(turn.Role)
		turn.Label = norm.Trim(turn.Label)
		turn.Content = norm.Trim(turn.Content)
		turn.TimeUTC = norm.Trim(turn.TimeUTC)
		turn.Provider = normalizeProvider(turn.Provider)
		turn.Decisions = session.Normalize(session.State{Decisions: turn.Decisions}).Decisions
		turn.Diagnostics = trust.NormalizeDiagnostics(turn.Diagnostics)
		if turn.ID == "" && turn.Role == "" && turn.Content == "" && len(turn.Diagnostics) == 0 {
			continue
		}
		out = append(out, turn)
	}
	slices.SortStableFunc(out, func(a, b Turn) int {
		return norm.CompareStrings(a.TimeUTC, b.TimeUTC, a.ID, b.ID, a.Role, b.Role, a.Label, b.Label)
	})
	return out
}

// DecisionEvent returns a generic transcript event for one decision record.
func DecisionEvent(record decision.Record) Event {
	record = decision.Normalize(record)
	return Event{
		Type:    "decision",
		Stage:   record.Stage,
		Message: record.Rationale,
		Fields:  decision.EventFields(record),
	}
}

// DecisionEvents returns generic transcript events for decision records.
func DecisionEvents(records []decision.Record) []Event {
	records = decision.NormalizeAll(records)
	events := make([]Event, 0, len(records))
	for _, record := range records {
		events = append(events, DecisionEvent(record))
	}
	return events
}

func normalizeEvents(events []Event) []Event {
	out := make([]Event, 0, len(events))
	for _, event := range events {
		event.ID = norm.Trim(event.ID)
		event.Type = norm.Token(event.Type)
		event.Stage = norm.Token(event.Stage)
		event.Severity = norm.Token(event.Severity)
		event.Message = norm.Trim(event.Message)
		event.TimeUTC = norm.Trim(event.TimeUTC)
		event.Fields = norm.Metadata(event.Fields)
		event.Diagnostics = trust.NormalizeDiagnostics(event.Diagnostics)
		if event.Type == "" && event.Message == "" && len(event.Diagnostics) == 0 {
			continue
		}
		out = append(out, event)
	}
	slices.SortStableFunc(out, func(a, b Event) int {
		if diff := norm.CompareSeverity(a.Severity, b.Severity); diff != 0 {
			return diff
		}
		return norm.CompareStrings(a.TimeUTC, b.TimeUTC, a.Type, b.Type, a.Stage, b.Stage, a.ID, b.ID)
	})
	return out
}

func normalizeProvider(provider *ModelProvenance) *ModelProvenance {
	if provider == nil {
		return nil
	}
	out := *provider
	out.Provider = norm.Token(out.Provider)
	out.Model = norm.Trim(out.Model)
	out.Endpoint = norm.Trim(out.Endpoint)
	out.RequestID = norm.Trim(out.RequestID)
	out.ResponseID = norm.Trim(out.ResponseID)
	out.Seed = norm.Trim(out.Seed)
	if out.Provider == "" && out.Model == "" && out.Endpoint == "" && out.RequestID == "" && out.ResponseID == "" && out.Seed == "" {
		return nil
	}
	return &out
}
