package report

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"github.com/OpenUdon/authoring/decision"
	"github.com/OpenUdon/authoring/readiness"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
	"github.com/OpenUdon/authoring/trust"
)

const (
	Version = "authoring.agent-result.v1"

	StatusComplete   = "complete"
	StatusNeedsInput = "needs_input"
	StatusFailed     = "failed"
	StatusCanceled   = "canceled"
)

// Result is the durable product-neutral noninteractive authoring outcome.
type Result struct {
	Version        string                   `json:"version"`
	Status         string                   `json:"status"`
	Summary        string                   `json:"summary,omitempty"`
	Readiness      *readiness.Result        `json:"readiness,omitempty"`
	TopIssue       *readiness.Issue         `json:"top_issue,omitempty"`
	RepairStatus   string                   `json:"repair_status,omitempty"`
	RepairAttempts int                      `json:"repair_attempts,omitempty"`
	Decisions      []DecisionBehavior       `json:"decisions,omitempty"`
	Diagnostics    []trust.DiagnosticRecord `json:"diagnostics,omitempty"`
	Artifacts      []trust.ArtifactRecord   `json:"artifacts,omitempty"`
	Digests        []trust.DigestRecord     `json:"digests,omitempty"`
	Session        *SessionMetadata         `json:"session,omitempty"`
	Transcript     *TranscriptMetadata      `json:"transcript,omitempty"`
	Metadata       map[string]string        `json:"metadata,omitempty"`
}

// DecisionBehavior is the result-safe summary of one decision record.
type DecisionBehavior struct {
	Stage                string `json:"stage,omitempty"`
	Slot                 string `json:"slot,omitempty"`
	Confidence           string `json:"confidence,omitempty"`
	Behavior             string `json:"behavior,omitempty"`
	RequiresConfirmation bool   `json:"requires_confirmation,omitempty"`
}

// SessionMetadata summarizes an M03 session without embedding prompt answers.
type SessionMetadata struct {
	Version             string `json:"version,omitempty"`
	ID                  string `json:"id,omitempty"`
	Goal                string `json:"goal,omitempty"`
	Mode                string `json:"mode,omitempty"`
	CreatedUTC          string `json:"created_utc,omitempty"`
	UpdatedUTC          string `json:"updated_utc,omitempty"`
	TurnCount           int    `json:"turn_count,omitempty"`
	AnswerCount         int    `json:"answer_count,omitempty"`
	DecisionCount       int    `json:"decision_count,omitempty"`
	ReadinessIssueCount int    `json:"readiness_issue_count,omitempty"`
	ArtifactCount       int    `json:"artifact_count,omitempty"`
	DiagnosticCount     int    `json:"diagnostic_count,omitempty"`
}

// TranscriptMetadata summarizes an M03 transcript without embedding turns.
type TranscriptMetadata struct {
	Version         string                      `json:"version,omitempty"`
	SessionID       string                      `json:"session_id,omitempty"`
	TimeUTC         string                      `json:"time_utc,omitempty"`
	Provider        *transcript.ModelProvenance `json:"provider,omitempty"`
	TurnCount       int                         `json:"turn_count,omitempty"`
	EventCount      int                         `json:"event_count,omitempty"`
	DiagnosticCount int                         `json:"diagnostic_count,omitempty"`
	ArtifactCount   int                         `json:"artifact_count,omitempty"`
}

// Normalize returns a deterministic result.
func Normalize(result Result) Result {
	result.Version = firstNonEmpty(strings.TrimSpace(result.Version), Version)
	result.Status = normalizeStatus(result.Status)
	result.Summary = strings.TrimSpace(result.Summary)
	result.Readiness = normalizeReadinessResult(result.Readiness)
	result.TopIssue = normalizeTopIssue(result.TopIssue, result.Readiness)
	result.RepairStatus = normalizeToken(result.RepairStatus)
	if result.RepairAttempts < 0 {
		result.RepairAttempts = 0
	}
	result.Decisions = NormalizeDecisionBehaviors(result.Decisions)
	result.Diagnostics = trust.NormalizeDiagnostics(result.Diagnostics)
	result.Artifacts = normalizeArtifacts(result.Artifacts)
	result.Digests = normalizeDigests(result.Digests)
	result.Session = NormalizeSessionMetadata(result.Session)
	result.Transcript = NormalizeTranscriptMetadata(result.Transcript)
	result.Metadata = normalizeMetadata(result.Metadata)
	return result
}

// CanonicalJSON returns deterministic indented JSON for result.
func CanonicalJSON(result Result) ([]byte, error) {
	return json.MarshalIndent(Normalize(result), "", "  ")
}

// StatusForError returns the generic result status implied by err.
func StatusForError(err error) string {
	if err == nil {
		return StatusComplete
	}
	if isCanceled(err) {
		return StatusCanceled
	}
	return StatusFailed
}

// DecisionBehaviors returns deterministic result-safe summaries for decision
// records.
func DecisionBehaviors(records []decision.Record) []DecisionBehavior {
	out := make([]DecisionBehavior, 0, len(records))
	for _, record := range decision.NormalizeAll(records) {
		out = append(out, DecisionBehavior{
			Stage:                record.Stage,
			Slot:                 record.Slot,
			Confidence:           record.Confidence,
			Behavior:             decision.Behavior(record),
			RequiresConfirmation: decision.RequiresConfirmation(record),
		})
	}
	return NormalizeDecisionBehaviors(out)
}

// NormalizeDecisionBehaviors returns deterministic decision summaries.
func NormalizeDecisionBehaviors(records []DecisionBehavior) []DecisionBehavior {
	out := make([]DecisionBehavior, 0, len(records))
	for _, record := range records {
		record.Stage = normalizeToken(record.Stage)
		record.Slot = strings.TrimSpace(record.Slot)
		record.Confidence = decision.NormalizeConfidence(record.Confidence)
		record.Behavior = normalizeDecisionBehavior(record.Behavior, record.Confidence, record.RequiresConfirmation)
		if record.Stage == "" && record.Slot == "" && record.Confidence == decision.ConfidenceReview && !record.RequiresConfirmation {
			continue
		}
		out = append(out, record)
	}
	slices.SortStableFunc(out, func(a, b DecisionBehavior) int {
		return compareStrings(a.Stage, b.Stage, a.Slot, b.Slot, a.Confidence, b.Confidence, a.Behavior, b.Behavior)
	})
	return out
}

// SessionMetadataFromState summarizes a normalized session state.
func SessionMetadataFromState(state session.State) *SessionMetadata {
	state = session.Normalize(state)
	return NormalizeSessionMetadata(&SessionMetadata{
		Version:             state.Version,
		ID:                  state.ID,
		Goal:                state.Goal,
		Mode:                state.Mode,
		CreatedUTC:          state.CreatedUTC,
		UpdatedUTC:          state.UpdatedUTC,
		TurnCount:           len(state.Turns),
		AnswerCount:         len(state.Answers),
		DecisionCount:       len(state.Decisions),
		ReadinessIssueCount: len(state.Readiness),
		ArtifactCount:       len(state.Artifacts),
		DiagnosticCount:     len(state.Diagnostics),
	})
}

// NormalizeSessionMetadata returns deterministic session metadata.
func NormalizeSessionMetadata(metadata *SessionMetadata) *SessionMetadata {
	if metadata == nil {
		return nil
	}
	out := *metadata
	out.Version = strings.TrimSpace(out.Version)
	out.ID = strings.TrimSpace(out.ID)
	out.Goal = strings.TrimSpace(out.Goal)
	out.Mode = normalizeToken(out.Mode)
	out.CreatedUTC = strings.TrimSpace(out.CreatedUTC)
	out.UpdatedUTC = strings.TrimSpace(out.UpdatedUTC)
	out.TurnCount = nonnegative(out.TurnCount)
	out.AnswerCount = nonnegative(out.AnswerCount)
	out.DecisionCount = nonnegative(out.DecisionCount)
	out.ReadinessIssueCount = nonnegative(out.ReadinessIssueCount)
	out.ArtifactCount = nonnegative(out.ArtifactCount)
	out.DiagnosticCount = nonnegative(out.DiagnosticCount)
	if emptySessionMetadata(out) {
		return nil
	}
	return &out
}

// TranscriptMetadataFromRecord summarizes a normalized transcript record.
func TranscriptMetadataFromRecord(record transcript.Record) *TranscriptMetadata {
	record = transcript.Normalize(record)
	return NormalizeTranscriptMetadata(&TranscriptMetadata{
		Version:         record.Version,
		SessionID:       record.SessionID,
		TimeUTC:         record.TimeUTC,
		Provider:        record.Provider,
		TurnCount:       len(record.Turns),
		EventCount:      len(record.Events),
		DiagnosticCount: len(record.Diagnostics),
		ArtifactCount:   len(record.Artifacts),
	})
}

// NormalizeTranscriptMetadata returns deterministic transcript metadata.
func NormalizeTranscriptMetadata(metadata *TranscriptMetadata) *TranscriptMetadata {
	if metadata == nil {
		return nil
	}
	out := *metadata
	out.Version = strings.TrimSpace(out.Version)
	out.SessionID = strings.TrimSpace(out.SessionID)
	out.TimeUTC = strings.TrimSpace(out.TimeUTC)
	out.Provider = normalizeProvider(out.Provider)
	out.TurnCount = nonnegative(out.TurnCount)
	out.EventCount = nonnegative(out.EventCount)
	out.DiagnosticCount = nonnegative(out.DiagnosticCount)
	out.ArtifactCount = nonnegative(out.ArtifactCount)
	if emptyTranscriptMetadata(out) {
		return nil
	}
	return &out
}

func normalizeStatus(status string) string {
	switch normalizeToken(status) {
	case StatusComplete:
		return StatusComplete
	case StatusNeedsInput, "needs-input", "need_input", "need-input":
		return StatusNeedsInput
	case StatusCanceled, "cancelled":
		return StatusCanceled
	case StatusFailed, "":
		return StatusFailed
	default:
		return StatusFailed
	}
}

func normalizeReadinessResult(result *readiness.Result) *readiness.Result {
	if result == nil {
		return nil
	}
	out := readiness.NormalizeResult(*result)
	return &out
}

func normalizeTopIssue(issue *readiness.Issue, result *readiness.Result) *readiness.Issue {
	if result != nil && result.TopIssue != nil {
		top := *result.TopIssue
		return &top
	}
	if issue == nil {
		return nil
	}
	issues := readiness.NormalizeIssues([]readiness.Issue{*issue})
	if len(issues) == 0 {
		return nil
	}
	return &issues[0]
}

func normalizeDecisionBehavior(behavior, confidence string, requiresConfirmation bool) string {
	behavior = normalizeToken(behavior)
	switch behavior {
	case decision.BehaviorAutoAccept, decision.BehaviorReview, decision.BehaviorLowConfidence, decision.BehaviorConflict:
		return behavior
	}
	record := decision.Record{Confidence: confidence, RequiresConfirmation: requiresConfirmation}
	return decision.Behavior(record)
}

func normalizeArtifacts(artifacts []trust.ArtifactRecord) []trust.ArtifactRecord {
	out := make([]trust.ArtifactRecord, 0, len(artifacts))
	for _, record := range artifacts {
		record.Path = strings.TrimSpace(record.Path)
		record.Kind = normalizeToken(record.Kind)
		record.MediaType = strings.TrimSpace(record.MediaType)
		record.Classification = normalizeToken(record.Classification)
		record.Digest = normalizeDigest(record.Digest)
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

func normalizeDigests(digests []trust.DigestRecord) []trust.DigestRecord {
	out := make([]trust.DigestRecord, 0, len(digests))
	for _, record := range digests {
		record = normalizeDigest(record)
		if record.Algorithm == "" && record.Value == "" {
			continue
		}
		out = append(out, record)
	}
	slices.SortStableFunc(out, func(a, b trust.DigestRecord) int {
		return compareStrings(a.Algorithm, b.Algorithm, a.Value, b.Value)
	})
	return out
}

func normalizeDigest(record trust.DigestRecord) trust.DigestRecord {
	record.Algorithm = normalizeToken(record.Algorithm)
	record.Value = strings.TrimSpace(record.Value)
	return record
}

func normalizeProvider(provider *transcript.ModelProvenance) *transcript.ModelProvenance {
	if provider == nil {
		return nil
	}
	out := *provider
	out.Provider = normalizeToken(out.Provider)
	out.Model = strings.TrimSpace(out.Model)
	out.Endpoint = strings.TrimSpace(out.Endpoint)
	out.RequestID = strings.TrimSpace(out.RequestID)
	out.ResponseID = strings.TrimSpace(out.ResponseID)
	out.Seed = strings.TrimSpace(out.Seed)
	if out.Provider == "" && out.Model == "" && out.Endpoint == "" && out.RequestID == "" && out.ResponseID == "" && out.Seed == "" {
		return nil
	}
	return &out
}

func normalizeMetadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
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

func emptySessionMetadata(metadata SessionMetadata) bool {
	return metadata.Version == "" && metadata.ID == "" && metadata.Goal == "" && metadata.Mode == "" &&
		metadata.CreatedUTC == "" && metadata.UpdatedUTC == "" && metadata.TurnCount == 0 &&
		metadata.AnswerCount == 0 && metadata.DecisionCount == 0 && metadata.ReadinessIssueCount == 0 &&
		metadata.ArtifactCount == 0 && metadata.DiagnosticCount == 0
}

func emptyTranscriptMetadata(metadata TranscriptMetadata) bool {
	return metadata.Version == "" && metadata.SessionID == "" && metadata.TimeUTC == "" && metadata.Provider == nil &&
		metadata.TurnCount == 0 && metadata.EventCount == 0 && metadata.DiagnosticCount == 0 && metadata.ArtifactCount == 0
}

func isCanceled(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
		strings.Contains(strings.ToLower(err.Error()), "canceled") ||
		strings.Contains(strings.ToLower(err.Error()), "cancelled")
}

func nonnegative(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), "_"))
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
