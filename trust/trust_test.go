package trust

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/evidence/digest"
)

func TestNormalizeDiagnosticsUsesEvidenceOrdering(t *testing.T) {
	records := []DiagnosticRecord{
		{Code: "z", Severity: "info", Message: " later "},
		{Code: "a", Severity: "error", Message: " first "},
		{Code: "m", Severity: "", Message: " warning "},
	}

	got := NormalizeDiagnostics(records)
	if len(got) != 3 {
		t.Fatalf("NormalizeDiagnostics returned %d records, want 3", len(got))
	}
	if got[0].Code != "a" || got[0].Severity != "error" || got[0].Message != "first" {
		t.Fatalf("first diagnostic = %#v, want normalized error diagnostic", got[0])
	}
	if got[1].Severity != "warning" {
		t.Fatalf("second diagnostic severity = %q, want warning", got[1].Severity)
	}
	if !HasBlockingDiagnostics(got) {
		t.Fatalf("HasBlockingDiagnostics returned false for error diagnostic")
	}
}

func TestRedactDocumentUsesEvidenceRedaction(t *testing.T) {
	input := map[string]any{
		"api_token": "secret-token-value",
		"nested": map[string]any{
			"note": "Authorization: Bearer abcdefghijklmnopqrstuvwxyz",
		},
	}

	got, changed := RedactDocument(input)
	if !changed {
		t.Fatalf("RedactDocument changed = false, want true")
	}
	redacted, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("RedactDocument returned %T, want map[string]any", got)
	}
	if redacted["api_token"] != RedactedValue {
		t.Fatalf("api_token = %q, want shared redaction marker", redacted["api_token"])
	}
	nested := redacted["nested"].(map[string]any)
	if !strings.Contains(nested["note"].(string), RedactedValue) {
		t.Fatalf("nested note = %q, want redacted bearer token", nested["note"])
	}
}

func TestArtifactPathAndDigestUseEvidencePrimitives(t *testing.T) {
	path, err := CleanArtifactPath("./reports/result.json")
	if err != nil {
		t.Fatalf("CleanArtifactPath returned error: %v", err)
	}
	if path != "reports/result.json" {
		t.Fatalf("CleanArtifactPath = %q, want reports/result.json", path)
	}
	if _, err := CleanArtifactPath("../escape.json"); err == nil {
		t.Fatalf("CleanArtifactPath accepted escaping path")
	}

	got := SHA256Bytes([]byte("authoring"))
	if got.Algorithm != digest.AlgorithmSHA256 || got.Value == "" {
		t.Fatalf("SHA256Bytes = %#v, want Evidence SHA-256 digest record", got)
	}
}

func TestNewArtifactManifestUsesEvidenceArtifactRecords(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "reports"), 0o755); err != nil {
		t.Fatalf("mkdir reports: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "reports", "result.json"), []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatalf("write result artifact: %v", err)
	}

	manifest, err := NewArtifactManifest(ArtifactManifestOptions{
		Root: root,
		Files: []ArtifactFileOptions{
			{
				Path:      "reports/result.json",
				Kind:      "report",
				MediaType: "application/json",
				Required:  true,
			},
		},
	})
	if err != nil {
		t.Fatalf("NewArtifactManifest returned error: %v", err)
	}
	if len(manifest.Artifacts) != 1 {
		t.Fatalf("manifest has %d artifacts, want 1", len(manifest.Artifacts))
	}
	record := manifest.Artifacts[0]
	if record.Path != "reports/result.json" || record.Kind != "report" || !record.Required {
		t.Fatalf("artifact record = %#v, want report record for reports/result.json", record)
	}
	if record.Digest.Algorithm != digest.AlgorithmSHA256 || record.SizeBytes == 0 {
		t.Fatalf("artifact digest/size = %#v/%d, want Evidence digest and size", record.Digest, record.SizeBytes)
	}
	if _, err := manifest.Digest(); err != nil {
		t.Fatalf("manifest digest returned error: %v", err)
	}
}
