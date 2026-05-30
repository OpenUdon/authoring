package trust

import (
	"github.com/OpenUdon/evidence/artifact"
	"github.com/OpenUdon/evidence/diagnostic"
	"github.com/OpenUdon/evidence/digest"
	"github.com/OpenUdon/evidence/redact"
)

const (
	// RedactedValue is the shared marker used for redacted provider or model
	// output before durable Authoring persistence.
	RedactedValue = redact.Value
)

type (
	// ArtifactRecord describes a generated artifact without embedding its
	// content.
	ArtifactRecord = artifact.Record
	// ArtifactManifest describes a deterministic generated artifact set.
	ArtifactManifest = artifact.Manifest
	// ArtifactFileOptions supplies caller-owned metadata for one artifact file.
	ArtifactFileOptions = artifact.FileOptions
	// ArtifactManifestOptions configures artifact manifest generation.
	ArtifactManifestOptions = artifact.ManifestOptions
	// DiagnosticRecord describes a product-neutral validation, readiness,
	// review, or runtime issue.
	DiagnosticRecord = diagnostic.Record
	// DiagnosticLocation identifies the artifact, source, resource, or field
	// tied to a diagnostic without requiring product-specific types.
	DiagnosticLocation = diagnostic.Location
	// DigestRecord describes a digest with an explicit algorithm and encoded
	// value.
	DigestRecord = digest.Record
	// RedactionOptions configures shared redaction behavior.
	RedactionOptions = redact.Options
)

// NormalizeDiagnostics returns normalized diagnostics ordered
// deterministically.
func NormalizeDiagnostics(records []DiagnosticRecord) []DiagnosticRecord {
	return diagnostic.Sort(records)
}

// HasBlockingDiagnostics reports whether records contains any error or
// blocking diagnostic.
func HasBlockingDiagnostics(records []DiagnosticRecord) bool {
	return diagnostic.HasErrors(records)
}

// RedactString redacts secret-like content from value.
func RedactString(value string, opts ...RedactionOptions) string {
	return redact.String(value, opts...)
}

// RedactDocument redacts secret-like values from common JSON/YAML-like
// documents and reports whether any replacement was made.
func RedactDocument(value any, opts ...RedactionOptions) (any, bool) {
	return redact.Document(value, opts...)
}

// CleanArtifactPath returns a canonical safe relative artifact path.
func CleanArtifactPath(path string) (string, error) {
	return artifact.CleanRelativePath(path)
}

// NewArtifactManifest returns a deterministic manifest for caller-supplied
// generated artifact files.
func NewArtifactManifest(opts ArtifactManifestOptions) (ArtifactManifest, error) {
	return artifact.NewManifest(opts)
}

// ArtifactRecordForFile returns a safe artifact record for one generated file.
func ArtifactRecordForFile(root string, opts ArtifactFileOptions) (ArtifactRecord, error) {
	return artifact.FileRecord(root, opts)
}

// SHA256Bytes returns the shared SHA-256 digest record for data.
func SHA256Bytes(data []byte) DigestRecord {
	return digest.SHA256Bytes(data)
}
