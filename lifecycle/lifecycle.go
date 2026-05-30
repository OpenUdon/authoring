package lifecycle

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
	"github.com/OpenUdon/authoring/trust"
)

const (
	// DraftVersion is the generic draft envelope version.
	DraftVersion = "authoring.draft.v1"
)

// Draft is a generic durable draft envelope.
type Draft[T any] struct {
	Version    string                 `json:"version"`
	SavedUTC   string                 `json:"saved_utc"`
	Draft      T                      `json:"draft"`
	Session    *session.State         `json:"session,omitempty"`
	Transcript *transcript.Record     `json:"transcript,omitempty"`
	Artifacts  []trust.ArtifactRecord `json:"artifacts,omitempty"`
	Metadata   map[string]string      `json:"metadata,omitempty"`
}

// DefaultDraftPath returns Authoring's default local draft path below root.
func DefaultDraftPath(root string) string {
	return filepath.Join(root, ".authoring", "draft.json")
}

// LoadDraft reads a JSON draft envelope. It returns ok=false when the draft
// file does not exist.
func LoadDraft[T any](path string) (Draft[T], bool, error) {
	var zero Draft[T]
	path = strings.TrimSpace(path)
	if path == "" {
		return zero, false, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return zero, false, nil
	}
	if err != nil {
		return zero, false, err
	}
	var draft Draft[T]
	if err := json.Unmarshal(data, &draft); err != nil {
		return zero, false, err
	}
	return NormalizeDraft(draft), true, nil
}

// SaveDraft writes a JSON draft envelope atomically with private-file
// permissions. Empty paths are ignored.
func SaveDraft[T any](path string, draft Draft[T]) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	return WriteJSON(path, NormalizeDraft(draft), 0o600)
}

// DeleteDraft removes a draft file and prunes its parent directory if empty.
// Missing files are not an error.
func DeleteDraft(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	if err != nil {
		return err
	}
	_ = os.Remove(filepath.Dir(path))
	return nil
}

// Autosave returns a function that writes draft values to path through the
// generic draft envelope.
func Autosave[T any](path string, base Draft[T], normalize func(*T)) func(T) error {
	return func(value T) error {
		if normalize != nil {
			normalize(&value)
		}
		base.Draft = value
		return SaveDraft(path, base)
	}
}

// SaveTranscript writes a transcript record atomically with private-file
// permissions. Empty paths are ignored.
func SaveTranscript(path string, record transcript.Record) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	return WriteJSON(path, transcript.Normalize(record), 0o600)
}

// LoadTranscript reads a transcript record.
func LoadTranscript(path string) (transcript.Record, bool, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return transcript.Record{}, false, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return transcript.Record{}, false, nil
	}
	if err != nil {
		return transcript.Record{}, false, err
	}
	var record transcript.Record
	if err := json.Unmarshal(data, &record); err != nil {
		return transcript.Record{}, false, err
	}
	return transcript.Normalize(record), true, nil
}

// WriteJSON writes value as deterministic indented JSON with a trailing
// newline.
func WriteJSON(path string, value any, perm os.FileMode) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return AtomicWrite(path, data, perm)
}

// AtomicWrite writes data to path through a sibling temp file and rename.
func AtomicWrite(path string, data []byte, perm os.FileMode) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp.")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	return os.Chmod(path, perm)
}

// FileArtifact returns a safe artifact record for one generated file.
func FileArtifact(root string, opts trust.ArtifactFileOptions) (trust.ArtifactRecord, error) {
	return trust.ArtifactRecordForFile(root, opts)
}

// ArtifactManifest returns a deterministic manifest for generated files.
func ArtifactManifest(root string, files []trust.ArtifactFileOptions) (trust.ArtifactManifest, error) {
	return trust.NewArtifactManifest(trust.ArtifactManifestOptions{Root: root, Files: files})
}

// NormalizeDraft returns a deterministic draft envelope copy.
func NormalizeDraft[T any](draft Draft[T]) Draft[T] {
	draft.Version = strings.TrimSpace(draft.Version)
	if draft.Version == "" {
		draft.Version = DraftVersion
	}
	draft.SavedUTC = strings.TrimSpace(draft.SavedUTC)
	if draft.SavedUTC == "" {
		draft.SavedUTC = time.Now().UTC().Format(time.RFC3339)
	}
	if draft.Session != nil {
		normalized := session.Normalize(*draft.Session)
		draft.Session = &normalized
	}
	if draft.Transcript != nil {
		normalized := transcript.Normalize(*draft.Transcript)
		draft.Transcript = &normalized
	}
	draft.Artifacts = normalizeArtifacts(draft.Artifacts)
	draft.Metadata = normalizeMetadata(draft.Metadata)
	return draft
}

func normalizeArtifacts(artifacts []trust.ArtifactRecord) []trust.ArtifactRecord {
	normalized := session.Normalize(session.State{Artifacts: artifacts})
	return normalized.Artifacts
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
