package report

import (
	"context"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/decision"
	"github.com/OpenUdon/authoring/readiness"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
	"github.com/OpenUdon/authoring/trust"
)

func TestNormalizeNeedsInputResult(t *testing.T) {
	result := Normalize(Result{
		Status:  " needs input ",
		Summary: " missing goal ",
		Readiness: &readiness.Result{Issues: []readiness.Issue{
			{Code: "warn", Severity: "warning"},
			{Code: "missing_goal", Severity: "blocking", Slot: "goal"},
		}},
		Decisions: DecisionBehaviors([]decision.Record{
			{Stage: "draft", Slot: "operation", Confidence: "low"},
		}),
		Session: SessionMetadataFromState(session.State{
			ID:    " s1 ",
			Goal:  " make workflow ",
			Turns: []session.PromptTurn{{ID: "1", Question: "Goal?", Answer: "x"}},
		}),
	})
	if result.Version != Version || result.Status != StatusNeedsInput || result.Summary != "missing goal" {
		t.Fatalf("result = %#v, want normalized needs-input result", result)
	}
	if result.TopIssue == nil || result.TopIssue.Code != "missing_goal" {
		t.Fatalf("top issue = %#v, want missing goal", result.TopIssue)
	}
	if len(result.Decisions) != 1 || result.Decisions[0].Behavior != decision.BehaviorLowConfidence || !result.Decisions[0].RequiresConfirmation {
		t.Fatalf("decisions = %#v, want low-confidence confirmation", result.Decisions)
	}
	if result.Session == nil || result.Session.ID != "s1" || result.Session.TurnCount != 1 {
		t.Fatalf("session = %#v, want summary metadata", result.Session)
	}
}

func TestStatusForCancellation(t *testing.T) {
	if StatusForError(nil) != StatusComplete {
		t.Fatalf("nil error status = %q", StatusForError(nil))
	}
	if Normalize(Result{Status: "needs-input"}).Status != StatusNeedsInput {
		t.Fatalf("needs-input spelling did not normalize")
	}
	if Normalize(Result{Status: "Needs-Input"}).Status != StatusNeedsInput {
		t.Fatalf("Needs-Input spelling did not normalize")
	}
	if NormalizeOutcome("Needs-Input") != StatusNeedsInput {
		t.Fatalf("Needs-Input outcome did not normalize")
	}
	if StatusForError(context.Canceled) != StatusCanceled {
		t.Fatalf("canceled status = %q", StatusForError(context.Canceled))
	}
	if StatusForError(context.DeadlineExceeded) != StatusCanceled {
		t.Fatalf("deadline status = %q", StatusForError(context.DeadlineExceeded))
	}
}

func TestNormalizeDecisionBehaviorsKeepsConfirmation(t *testing.T) {
	records := NormalizeDecisionBehaviors([]DecisionBehavior{{RequiresConfirmation: true}})
	if len(records) != 1 || records[0].Behavior != decision.BehaviorReview || !records[0].RequiresConfirmation {
		t.Fatalf("records = %#v, want confirmation-only review behavior", records)
	}
}

func TestFailureDiagnosticsArtifactsAndJSON(t *testing.T) {
	digest := trust.SHA256Bytes([]byte("artifact"))
	result := Normalize(Result{
		Status: "unknown",
		Diagnostics: []trust.DiagnosticRecord{
			{Code: "z", Severity: "warning", Message: "later"},
			{Code: "a", Severity: "error", Message: "first"},
		},
		Artifacts: []trust.ArtifactRecord{
			{Path: "z.json", Kind: " report ", Digest: digest},
			{Path: "a.json", Kind: " report ", Digest: digest},
		},
		Digests: []trust.DigestRecord{
			{Algorithm: " SHA256 ", Value: " z "},
			digest,
		},
		Transcript: TranscriptMetadataFromRecord(transcript.Record{
			SessionID: "s1",
			Provider: &transcript.ModelProvenance{
				Provider: " OpenAI ",
				Model:    "model",
			},
			Events: []transcript.Event{{Type: "failed", Message: "x"}},
		}),
	})
	if result.Status != StatusFailed {
		t.Fatalf("status = %q, want failed", result.Status)
	}
	if len(result.Diagnostics) != 2 || result.Diagnostics[0].Code != "a" {
		t.Fatalf("diagnostics = %#v, want deterministic diagnostics", result.Diagnostics)
	}
	if len(result.Artifacts) != 2 || result.Artifacts[0].Path != "a.json" || result.Artifacts[0].Kind != "report" {
		t.Fatalf("artifacts = %#v, want sorted normalized artifacts", result.Artifacts)
	}
	if result.Transcript == nil || result.Transcript.Provider.Provider != "openai" || result.Transcript.EventCount != 1 {
		t.Fatalf("transcript = %#v, want normalized transcript metadata", result.Transcript)
	}
	data, err := CanonicalJSON(result)
	if err != nil {
		t.Fatalf("CanonicalJSON error = %v", err)
	}
	if !strings.Contains(string(data), `"status": "failed"`) || !strings.Contains(string(data), `"path": "a.json"`) {
		t.Fatalf("json = %s", data)
	}
}
