package lifecycle

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
	"github.com/OpenUdon/authoring/trust"
	"github.com/OpenUdon/evidence/digest"
)

type draftBody struct {
	Name string `json:"name"`
}

func TestDraftLifecycleSaveLoadAutosaveDelete(t *testing.T) {
	path := DefaultDraftPath(t.TempDir())
	state := session.State{ID: "s1", Goal: "draft it"}
	base := Draft[draftBody]{
		Session:  &state,
		Metadata: map[string]string{" product ": " openudon "},
	}

	autosave := Autosave(path, base, func(body *draftBody) {
		body.Name = "normalized-" + body.Name
	})
	if err := autosave(draftBody{Name: "draft"}); err != nil {
		t.Fatalf("autosave returned error: %v", err)
	}
	loaded, ok, err := LoadDraft[draftBody](path)
	if err != nil || !ok {
		t.Fatalf("LoadDraft = ok %v err %v, want saved draft", ok, err)
	}
	if loaded.Version != DraftVersion || loaded.Draft.Name != "normalized-draft" {
		t.Fatalf("loaded draft = %#v", loaded)
	}
	if loaded.Metadata["product"] != "openudon" {
		t.Fatalf("metadata = %#v, want normalized metadata", loaded.Metadata)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat draft: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("draft permissions = %o, want 0600", got)
	}
	if err := DeleteDraft(path); err != nil {
		t.Fatalf("DeleteDraft returned error: %v", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("draft still exists after delete: %v", err)
	}
}

func TestAtomicWriteOverwritesPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "file.json")
	if err := AtomicWrite(path, []byte("first"), 0o600); err != nil {
		t.Fatalf("AtomicWrite first returned error: %v", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatalf("chmod test file: %v", err)
	}
	if err := AtomicWrite(path, []byte("second"), 0o600); err != nil {
		t.Fatalf("AtomicWrite second returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "second" {
		t.Fatalf("file data = %q, want second", string(data))
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("file permissions = %o, want 0600", got)
	}
}

func TestTranscriptPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "transcript.json")
	record := transcript.Record{
		SessionID: "s1",
		Events:    []transcript.Event{{Type: "draft", Message: "created"}},
	}
	if err := SaveTranscript(path, record); err != nil {
		t.Fatalf("SaveTranscript returned error: %v", err)
	}
	loaded, ok, err := LoadTranscript(path)
	if err != nil || !ok {
		t.Fatalf("LoadTranscript = ok %v err %v, want saved transcript", ok, err)
	}
	if loaded.Version != transcript.Version || loaded.Events[0].Type != "draft" {
		t.Fatalf("loaded transcript = %#v", loaded)
	}
}

func TestSafeArtifactRecords(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "artifact.json"), []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	record, err := FileArtifact(root, trust.ArtifactFileOptions{Path: "artifact.json", Kind: "report"})
	if err != nil {
		t.Fatalf("FileArtifact returned error: %v", err)
	}
	if record.Path != "artifact.json" || record.Kind != "report" || record.Digest.Algorithm != digest.AlgorithmSHA256 {
		t.Fatalf("artifact record = %#v", record)
	}
	manifest, err := ArtifactManifest(root, []trust.ArtifactFileOptions{{Path: "artifact.json"}})
	if err != nil {
		t.Fatalf("ArtifactManifest returned error: %v", err)
	}
	if len(manifest.Artifacts) != 1 {
		t.Fatalf("manifest artifacts = %d, want 1", len(manifest.Artifacts))
	}
	if _, err := FileArtifact(root, trust.ArtifactFileOptions{Path: "../escape.json"}); err == nil {
		t.Fatalf("FileArtifact accepted escaping path")
	}
}
