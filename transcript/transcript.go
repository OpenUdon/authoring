package transcript

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/OpenUdon/authoring/decision"
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
	record.Version = firstNonEmpty(trim(record.Version), Version)
	record.SessionID = trim(record.SessionID)
	record.TimeUTC = trim(record.TimeUTC)
	record.Provider = normalizeProvider(record.Provider)
	record.Turns = normalizeTurns(record.Turns)
	record.Events = normalizeEvents(record.Events)
	record.Diagnostics = trust.NormalizeDiagnostics(record.Diagnostics)
	record.Artifacts = normalizeArtifacts(record.Artifacts)
	record.Metadata = normalizeMetadata(record.Metadata)
	return record
}

// CanonicalJSON returns deterministic indented JSON for record.
func CanonicalJSON(record Record) ([]byte, error) {
	return json.MarshalIndent(Normalize(record), "", "  ")
}

func normalizeTurns(turns []Turn) []Turn {
	out := make([]Turn, 0, len(turns))
	for _, turn := range turns {
		turn.ID = trim(turn.ID)
		turn.Role = normalizeToken(turn.Role)
		turn.Label = trim(turn.Label)
		turn.Content = trim(turn.Content)
		turn.TimeUTC = trim(turn.TimeUTC)
		turn.Provider = normalizeProvider(turn.Provider)
		turn.Decisions = session.Normalize(session.State{Decisions: turn.Decisions}).Decisions
		turn.Diagnostics = trust.NormalizeDiagnostics(turn.Diagnostics)
		if turn.ID == "" && turn.Role == "" && turn.Content == "" && len(turn.Diagnostics) == 0 {
			continue
		}
		out = append(out, turn)
	}
	slices.SortStableFunc(out, func(a, b Turn) int {
		return compareStrings(a.TimeUTC, b.TimeUTC, a.ID, b.ID, a.Role, b.Role, a.Label, b.Label)
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
		event.ID = trim(event.ID)
		event.Type = normalizeToken(event.Type)
		event.Stage = normalizeToken(event.Stage)
		event.Severity = normalizeToken(event.Severity)
		event.Message = trim(event.Message)
		event.TimeUTC = trim(event.TimeUTC)
		event.Fields = normalizeMetadata(event.Fields)
		event.Diagnostics = trust.NormalizeDiagnostics(event.Diagnostics)
		if event.Type == "" && event.Message == "" && len(event.Diagnostics) == 0 {
			continue
		}
		out = append(out, event)
	}
	slices.SortStableFunc(out, func(a, b Event) int {
		if diff := compareSeverity(a.Severity, b.Severity); diff != 0 {
			return diff
		}
		return compareStrings(a.TimeUTC, b.TimeUTC, a.Type, b.Type, a.Stage, b.Stage, a.ID, b.ID)
	})
	return out
}

func normalizeProvider(provider *ModelProvenance) *ModelProvenance {
	if provider == nil {
		return nil
	}
	out := *provider
	out.Provider = normalizeToken(out.Provider)
	out.Model = trim(out.Model)
	out.Endpoint = trim(out.Endpoint)
	out.RequestID = trim(out.RequestID)
	out.ResponseID = trim(out.ResponseID)
	out.Seed = trim(out.Seed)
	if out.Provider == "" && out.Model == "" && out.Endpoint == "" && out.RequestID == "" && out.ResponseID == "" && out.Seed == "" {
		return nil
	}
	return &out
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
