package decision

import (
	"slices"
	"strings"

	"github.com/OpenUdon/authoring/trust"
)

const (
	ConfidenceAutoAccept = "auto_accept"
	ConfidenceReview     = "review"
	ConfidenceLow        = "low_confidence"
	ConfidenceConflict   = "conflict"

	BehaviorAutoAccept    = "auto_accept"
	BehaviorReview        = "review"
	BehaviorLowConfidence = "low_confidence"
	BehaviorConflict      = "conflict"
)

// Record describes why a slot value was chosen or still requires
// confirmation.
type Record struct {
	Stage                string        `json:"stage,omitempty"`
	Slot                 string        `json:"slot,omitempty"`
	Value                string        `json:"value,omitempty"`
	Source               string        `json:"source,omitempty"`
	Confidence           string        `json:"confidence,omitempty"`
	Rationale            string        `json:"rationale,omitempty"`
	Evidence             string        `json:"evidence,omitempty"`
	Alternatives         []Alternative `json:"alternatives,omitempty"`
	RequiresConfirmation bool          `json:"requires_confirmation,omitempty"`
}

// Alternative records a rejected or lower-confidence value.
type Alternative struct {
	Value     string `json:"value,omitempty"`
	Rationale string `json:"rationale,omitempty"`
	Source    string `json:"source,omitempty"`
}

// Normalize returns a trimmed, deterministic record.
func Normalize(record Record) Record {
	record.Stage = normalizeToken(record.Stage)
	record.Slot = strings.TrimSpace(record.Slot)
	record.Value = strings.TrimSpace(record.Value)
	record.Source = normalizeToken(record.Source)
	record.Confidence = NormalizeConfidence(record.Confidence)
	record.Rationale = strings.TrimSpace(record.Rationale)
	record.Evidence = strings.TrimSpace(record.Evidence)
	record.Alternatives = NormalizeAlternatives(record.Alternatives)
	return record
}

// NormalizeAll normalizes and sorts records without merging conflicts.
func NormalizeAll(records []Record) []Record {
	out := make([]Record, 0, len(records))
	for _, record := range records {
		record = Normalize(record)
		if record.Stage == "" && record.Slot == "" && record.Value == "" {
			continue
		}
		out = append(out, record)
	}
	slices.SortStableFunc(out, Compare)
	return out
}

// Merge normalizes, merges duplicate records, and marks conflicting values for
// the same stage/slot as requiring confirmation.
func Merge(records []Record) []Record {
	seen := map[string]int{}
	out := make([]Record, 0, len(records))
	for _, record := range NormalizeAll(records) {
		key := record.Stage + "\x00" + record.Slot + "\x00" + record.Value + "\x00" + record.Source
		if existing, ok := seen[key]; ok {
			out[existing] = mergeRecord(out[existing], record)
			continue
		}
		seen[key] = len(out)
		out = append(out, record)
	}
	markConflicts(out)
	slices.SortStableFunc(out, Compare)
	return out
}

// NormalizeAlternatives normalizes and sorts alternatives.
func NormalizeAlternatives(alternatives []Alternative) []Alternative {
	seen := map[string]bool{}
	out := make([]Alternative, 0, len(alternatives))
	for _, alternative := range alternatives {
		alternative.Value = strings.TrimSpace(alternative.Value)
		alternative.Rationale = strings.TrimSpace(alternative.Rationale)
		alternative.Source = normalizeToken(alternative.Source)
		if alternative.Value == "" {
			continue
		}
		key := alternative.Value + "\x00" + alternative.Source
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, alternative)
	}
	slices.SortStableFunc(out, func(a, b Alternative) int {
		return compareStrings(a.Value, b.Value, a.Source, b.Source, a.Rationale, b.Rationale)
	})
	return out
}

// NormalizeConfidence maps caller-specific confidence strings to generic
// confidence states.
func NormalizeConfidence(confidence string) string {
	switch normalizeToken(confidence) {
	case "", "review", "medium", "needs_review":
		return ConfidenceReview
	case "auto", "auto_accept", "auto_accepted", "high", "accepted":
		return ConfidenceAutoAccept
	case "low", "low_confidence", "uncertain":
		return ConfidenceLow
	case "conflict", "conflicting":
		return ConfidenceConflict
	default:
		return ConfidenceReview
	}
}

// Behavior returns the generic action implied by record confidence.
func Behavior(record Record) string {
	record = Normalize(record)
	switch record.Confidence {
	case ConfidenceConflict:
		return BehaviorConflict
	case ConfidenceLow:
		return BehaviorLowConfidence
	case ConfidenceAutoAccept:
		if record.RequiresConfirmation {
			return BehaviorReview
		}
		return BehaviorAutoAccept
	default:
		return BehaviorReview
	}
}

// RequiresConfirmation reports whether record should be confirmed before being
// accepted.
func RequiresConfirmation(record Record) bool {
	return record.RequiresConfirmation || Behavior(record) != BehaviorAutoAccept
}

// Redact returns a copy with secret-like values redacted.
func Redact(record Record) Record {
	record.Value = trust.RedactString(record.Value)
	record.Rationale = trust.RedactString(record.Rationale)
	record.Evidence = trust.RedactString(record.Evidence)
	for i := range record.Alternatives {
		record.Alternatives[i].Value = trust.RedactString(record.Alternatives[i].Value)
		record.Alternatives[i].Rationale = trust.RedactString(record.Alternatives[i].Rationale)
	}
	return record
}

// EventFields returns generic transcript event fields for record.
func EventFields(record Record) map[string]string {
	record = Normalize(record)
	fields := map[string]string{}
	for key, value := range map[string]string{
		"stage":      record.Stage,
		"slot":       record.Slot,
		"source":     record.Source,
		"confidence": record.Confidence,
		"behavior":   Behavior(record),
	} {
		if value != "" {
			fields[key] = value
		}
	}
	return fields
}

// Compare orders decision records deterministically.
func Compare(a, b Record) int {
	a = Normalize(a)
	b = Normalize(b)
	return compareStrings(a.Stage, b.Stage, a.Slot, b.Slot, a.Value, b.Value, a.Source, b.Source)
}

func mergeRecord(base, overlay Record) Record {
	base.Confidence = strongerConfidence(base.Confidence, overlay.Confidence)
	if base.Rationale == "" {
		base.Rationale = overlay.Rationale
	}
	if base.Evidence == "" {
		base.Evidence = overlay.Evidence
	}
	base.Alternatives = NormalizeAlternatives(append(base.Alternatives, overlay.Alternatives...))
	base.RequiresConfirmation = base.RequiresConfirmation || overlay.RequiresConfirmation
	return base
}

func markConflicts(records []Record) {
	bySlot := map[string][]int{}
	values := map[string]map[string]bool{}
	for i, record := range records {
		key := record.Stage + "\x00" + record.Slot
		bySlot[key] = append(bySlot[key], i)
		if values[key] == nil {
			values[key] = map[string]bool{}
		}
		values[key][record.Value] = true
	}
	for key, indexes := range bySlot {
		if len(values[key]) <= 1 {
			continue
		}
		for _, index := range indexes {
			records[index].Confidence = ConfidenceConflict
			records[index].RequiresConfirmation = true
			records[index].Alternatives = conflictAlternatives(records, indexes, index)
		}
	}
}

func conflictAlternatives(records []Record, indexes []int, self int) []Alternative {
	alternatives := append([]Alternative(nil), records[self].Alternatives...)
	for _, index := range indexes {
		if index == self || records[index].Value == "" {
			continue
		}
		alternatives = append(alternatives, Alternative{
			Value:     records[index].Value,
			Rationale: records[index].Rationale,
			Source:    records[index].Source,
		})
	}
	return NormalizeAlternatives(alternatives)
}

func strongerConfidence(a, b string) string {
	a = NormalizeConfidence(a)
	b = NormalizeConfidence(b)
	if confidenceRank(b) < confidenceRank(a) {
		return b
	}
	return a
}

func confidenceRank(confidence string) int {
	switch NormalizeConfidence(confidence) {
	case ConfidenceConflict:
		return 0
	case ConfidenceLow:
		return 1
	case ConfidenceReview:
		return 2
	case ConfidenceAutoAccept:
		return 3
	default:
		return 4
	}
}

func normalizeToken(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "-", "_")
	return strings.Join(strings.Fields(value), "_")
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
