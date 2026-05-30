package decision

import (
	"testing"

	"github.com/OpenUdon/authoring/trust"
)

func TestNormalizeAndBehavior(t *testing.T) {
	record := Normalize(Record{
		Stage:      " Operation Selection ",
		Slot:       "resource.operation",
		Value:      " create ",
		Source:     "Model Draft",
		Confidence: "high",
		Alternatives: []Alternative{
			{Value: " update ", Source: "model"},
			{Value: " update ", Source: "model"},
		},
	})
	if record.Stage != "operation_selection" || record.Source != "model_draft" || record.Confidence != ConfidenceAutoAccept {
		t.Fatalf("record = %#v, want normalized decision", record)
	}
	if Behavior(record) != BehaviorAutoAccept || RequiresConfirmation(record) {
		t.Fatalf("behavior=%q requires=%v, want auto accept", Behavior(record), RequiresConfirmation(record))
	}
	if len(record.Alternatives) != 1 || record.Alternatives[0].Value != "update" {
		t.Fatalf("alternatives = %#v", record.Alternatives)
	}
	record.RequiresConfirmation = true
	if Behavior(record) != BehaviorReview || !RequiresConfirmation(record) {
		t.Fatalf("confirmation override did not force review")
	}
}

func TestMergeMarksConflicts(t *testing.T) {
	records := Merge([]Record{
		{Stage: "draft", Slot: "operation", Value: "create", Source: "model", Confidence: "high"},
		{Stage: "draft", Slot: "operation", Value: "delete", Source: "user", Confidence: "review"},
		{Stage: "draft", Slot: "operation", Value: "create", Source: "model", Confidence: "low", Rationale: "maybe"},
	})
	if len(records) != 2 {
		t.Fatalf("records = %#v, want two merged decisions", records)
	}
	for _, record := range records {
		if record.Confidence != ConfidenceConflict || !record.RequiresConfirmation {
			t.Fatalf("record = %#v, want conflict requiring confirmation", record)
		}
		if len(record.Alternatives) == 0 {
			t.Fatalf("record = %#v, want alternatives", record)
		}
	}
}

func TestRedactAndEventFields(t *testing.T) {
	record := Redact(Record{
		Stage:      "auth",
		Slot:       "api_key",
		Value:      "sk-proj-abcdefghijklmnopqrstuvwxyz",
		Source:     "user",
		Confidence: "low",
		Rationale:  "token=sk-proj-abcdefghijklmnopqrstuvwxyz",
	})
	if record.Value != trust.RedactedValue || record.Rationale == "token=sk-proj-abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("redacted record = %#v", record)
	}
	fields := EventFields(record)
	if fields["stage"] != "auth" || fields["slot"] != "api_key" || fields["behavior"] != BehaviorLowConfidence {
		t.Fatalf("event fields = %#v", fields)
	}
}
