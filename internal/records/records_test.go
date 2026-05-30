package records

import (
	"testing"

	"github.com/OpenUdon/authoring/trust"
)

func TestDigestAndDigests(t *testing.T) {
	got := Digests([]trust.DigestRecord{
		{Algorithm: " SHA-256 ", Value: " b "},
		{Algorithm: " ", Value: " "},
		{Algorithm: " MD5 ", Value: " a "},
	})
	if len(got) != 2 {
		t.Fatalf("Digests len = %d, want 2", len(got))
	}
	if got[0].Algorithm != "md5" || got[0].Value != "a" || got[1].Algorithm != "sha_256" || got[1].Value != "b" {
		t.Fatalf("Digests = %#v, want normalized sorted records", got)
	}
}

func TestArtifacts(t *testing.T) {
	got := Artifacts([]trust.ArtifactRecord{
		{Path: " z.json ", Kind: " Report ", Digest: trust.DigestRecord{Algorithm: "SHA-256", Value: " abc "}},
		{Path: " "},
		{Path: " a.json ", Kind: " Report "},
	})
	if len(got) != 2 {
		t.Fatalf("Artifacts len = %d, want 2", len(got))
	}
	if got[0].Path != "a.json" || got[1].Path != "z.json" || got[1].Kind != "report" || got[1].Digest.Algorithm != "sha_256" {
		t.Fatalf("Artifacts = %#v, want normalized sorted records", got)
	}
}
