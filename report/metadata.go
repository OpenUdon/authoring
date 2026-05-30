package report

import (
	"slices"
	"strings"

	"github.com/OpenUdon/authoring/trust"
)

const (
	MetadataVersion = "authoring.report-metadata.v1"

	RetentionEphemeral = "ephemeral"
	RetentionRun       = "run"
	RetentionArchive   = "archive"
)

// ReportMetadata describes a generated report without embedding
// product-specific report bodies.
type ReportMetadata struct {
	Version           string            `json:"version"`
	RunID             string            `json:"run_id,omitempty"`
	Command           string            `json:"command,omitempty"`
	Commit            string            `json:"commit,omitempty"`
	GeneratedUTC      string            `json:"generated_utc,omitempty"`
	RetentionClass    string            `json:"retention_class,omitempty"`
	ProviderOutput    bool              `json:"provider_output,omitempty"`
	ArchiveSafe       bool              `json:"archive_safe,omitempty"`
	RedactionRequired bool              `json:"redaction_required,omitempty"`
	DigestSidecars    []DigestSidecar   `json:"digest_sidecars,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// DigestSidecar records digest metadata for a report or generated artifact.
type DigestSidecar struct {
	Path      string             `json:"path,omitempty"`
	Kind      string             `json:"kind,omitempty"`
	SizeBytes int64              `json:"size_bytes,omitempty"`
	Digest    trust.DigestRecord `json:"digest"`
	Required  bool               `json:"required,omitempty"`
}

// NormalizeReportMetadata returns deterministic report metadata.
func NormalizeReportMetadata(metadata *ReportMetadata) *ReportMetadata {
	if metadata == nil {
		return nil
	}
	out := *metadata
	out.Version = firstNonEmpty(strings.TrimSpace(out.Version), MetadataVersion)
	out.RunID = strings.TrimSpace(out.RunID)
	out.Command = strings.TrimSpace(out.Command)
	out.Commit = strings.TrimSpace(out.Commit)
	out.GeneratedUTC = strings.TrimSpace(out.GeneratedUTC)
	out.RetentionClass = NormalizeRetentionClass(out.RetentionClass)
	if out.RedactionRequired {
		out.ArchiveSafe = false
	}
	out.DigestSidecars = NormalizeDigestSidecars(out.DigestSidecars)
	out.Metadata = normalizeMetadata(out.Metadata)
	if emptyReportMetadata(out) {
		return nil
	}
	return &out
}

// NormalizeRetentionClass maps caller-specific retention names to generic
// retention classes.
func NormalizeRetentionClass(retention string) string {
	switch normalizeToken(retention) {
	case RetentionEphemeral, "temporary", "temp":
		return RetentionEphemeral
	case RetentionArchive, "archival", "long_term":
		return RetentionArchive
	case RetentionRun, "":
		return RetentionRun
	default:
		return RetentionRun
	}
}

// DigestSidecarForBytes returns a sidecar for in-memory report bytes.
func DigestSidecarForBytes(path, kind string, data []byte) DigestSidecar {
	return NormalizeDigestSidecar(DigestSidecar{
		Path:      path,
		Kind:      kind,
		SizeBytes: int64(len(data)),
		Digest:    trust.SHA256Bytes(data),
	})
}

// DigestSidecarForArtifact returns a sidecar from an M02 artifact record.
func DigestSidecarForArtifact(record trust.ArtifactRecord) DigestSidecar {
	return NormalizeDigestSidecar(DigestSidecar{
		Path:      record.Path,
		Kind:      record.Kind,
		SizeBytes: record.SizeBytes,
		Digest:    record.Digest,
		Required:  record.Required,
	})
}

// NormalizeDigestSidecars returns deterministic sidecars.
func NormalizeDigestSidecars(sidecars []DigestSidecar) []DigestSidecar {
	out := make([]DigestSidecar, 0, len(sidecars))
	for _, sidecar := range sidecars {
		sidecar = NormalizeDigestSidecar(sidecar)
		if sidecar.Path == "" && sidecar.Digest.Algorithm == "" && sidecar.Digest.Value == "" {
			continue
		}
		out = append(out, sidecar)
	}
	slices.SortStableFunc(out, CompareDigestSidecar)
	return out
}

// NormalizeDigestSidecar returns deterministic digest sidecar metadata.
func NormalizeDigestSidecar(sidecar DigestSidecar) DigestSidecar {
	sidecar.Path = cleanSidecarPath(sidecar.Path)
	sidecar.Kind = normalizeToken(sidecar.Kind)
	if sidecar.SizeBytes < 0 {
		sidecar.SizeBytes = 0
	}
	sidecar.Digest = normalizeDigest(sidecar.Digest)
	return sidecar
}

func cleanSidecarPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	clean, err := trust.CleanArtifactPath(path)
	if err != nil {
		return ""
	}
	return clean
}

// CompareDigestSidecar orders sidecars deterministically.
func CompareDigestSidecar(a, b DigestSidecar) int {
	a = NormalizeDigestSidecar(a)
	b = NormalizeDigestSidecar(b)
	return compareStrings(a.Path, b.Path, a.Kind, b.Kind, a.Digest.Algorithm, b.Digest.Algorithm, a.Digest.Value, b.Digest.Value)
}

// WithRedactionRequirement returns metadata with RedactionRequired propagated
// from Evidence redaction checks over values.
func WithRedactionRequirement(metadata ReportMetadata, values ...any) ReportMetadata {
	if RedactionRequired(values...) {
		metadata.RedactionRequired = true
		metadata.ArchiveSafe = false
	}
	normalized := NormalizeReportMetadata(&metadata)
	if normalized == nil {
		return ReportMetadata{}
	}
	return *normalized
}

// RedactionRequired reports whether Evidence redaction would change any value.
func RedactionRequired(values ...any) bool {
	for _, value := range values {
		_, changed := trust.RedactDocument(value)
		if changed {
			return true
		}
	}
	return false
}

func emptyReportMetadata(metadata ReportMetadata) bool {
	return metadata.Version == MetadataVersion && metadata.RunID == "" && metadata.Command == "" &&
		metadata.Commit == "" && metadata.GeneratedUTC == "" && metadata.RetentionClass == RetentionRun &&
		!metadata.ProviderOutput && !metadata.ArchiveSafe && !metadata.RedactionRequired &&
		len(metadata.DigestSidecars) == 0 && len(metadata.Metadata) == 0
}
