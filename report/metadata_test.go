package report

import (
	"testing"

	"github.com/OpenUdon/authoring/trust"
)

func TestNormalizeReportMetadata(t *testing.T) {
	metadata := NormalizeReportMetadata(&ReportMetadata{
		RunID:          " run-1 ",
		Command:        " author ",
		Commit:         " abc123 ",
		GeneratedUTC:   "2026-05-30T00:00:00Z",
		RetentionClass: "Long Term",
		ProviderOutput: true,
		ArchiveSafe:    true,
		Metadata: map[string]string{
			" env ": " test ",
			"":      "drop",
		},
	})
	if metadata.Version != MetadataVersion || metadata.RetentionClass != RetentionArchive {
		t.Fatalf("metadata = %#v, want normalized report metadata", metadata)
	}
	if !metadata.ProviderOutput || !metadata.ArchiveSafe || metadata.Metadata["env"] != "test" {
		t.Fatalf("metadata = %#v, want provider/archive flags and metadata", metadata)
	}
}

func TestDigestSidecars(t *testing.T) {
	digest := trust.SHA256Bytes([]byte("artifact"))
	sidecars := NormalizeDigestSidecars([]DigestSidecar{
		DigestSidecarForBytes("z.json", " report ", []byte("report")),
		{Path: "../unsafe.json", Digest: digest},
		DigestSidecarForArtifact(trust.ArtifactRecord{
			Path:      "a.json",
			Kind:      " report ",
			SizeBytes: 8,
			Digest:    digest,
			Required:  true,
		}),
		{},
	})
	if len(sidecars) != 3 {
		t.Fatalf("sidecars = %#v, want empty sidecar removed", sidecars)
	}
	if sidecars[0].Path != "" {
		t.Fatalf("unsafe sidecar path = %q, want dropped path", sidecars[0].Path)
	}
	if sidecars[1].Path != "a.json" || sidecars[1].Kind != "report" || !sidecars[1].Required {
		t.Fatalf("sidecars = %#v, want artifact sidecar after digest-only sidecar", sidecars)
	}
	if sidecars[2].Digest.Value == "" || sidecars[2].SizeBytes == 0 {
		t.Fatalf("sidecar = %#v, want byte digest sidecar", sidecars[2])
	}
}

func TestRedactionRequirement(t *testing.T) {
	metadata := WithRedactionRequirement(ReportMetadata{
		RunID:          "run-1",
		RetentionClass: RetentionArchive,
		ArchiveSafe:    true,
	}, map[string]string{"api_token": "sk-proj-abcdefghijklmnopqrstuvwxyz"})
	if !metadata.RedactionRequired || metadata.ArchiveSafe {
		t.Fatalf("metadata = %#v, want redaction required and archive unsafe", metadata)
	}
	if RedactionRequired("plain text") {
		t.Fatalf("plain text should not require redaction")
	}
}
