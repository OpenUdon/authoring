package records

import (
	"slices"

	"github.com/OpenUdon/authoring/internal/norm"
	"github.com/OpenUdon/authoring/trust"
)

// Digest returns a normalized digest record.
func Digest(record trust.DigestRecord) trust.DigestRecord {
	record.Algorithm = norm.Token(record.Algorithm)
	record.Value = norm.Trim(record.Value)
	return record
}

// Digests returns normalized digest records in deterministic order.
func Digests(digests []trust.DigestRecord) []trust.DigestRecord {
	out := make([]trust.DigestRecord, 0, len(digests))
	for _, record := range digests {
		record = Digest(record)
		if record.Algorithm == "" && record.Value == "" {
			continue
		}
		out = append(out, record)
	}
	slices.SortStableFunc(out, func(a, b trust.DigestRecord) int {
		return norm.CompareStrings(a.Algorithm, b.Algorithm, a.Value, b.Value)
	})
	return out
}

// Artifacts returns normalized artifact records in deterministic order.
func Artifacts(artifacts []trust.ArtifactRecord) []trust.ArtifactRecord {
	out := make([]trust.ArtifactRecord, 0, len(artifacts))
	for _, record := range artifacts {
		record.Path = norm.Trim(record.Path)
		record.Kind = norm.Token(record.Kind)
		record.MediaType = norm.Trim(record.MediaType)
		record.Classification = norm.Token(record.Classification)
		record.Digest = Digest(record.Digest)
		if record.Path == "" {
			continue
		}
		out = append(out, record)
	}
	slices.SortStableFunc(out, func(a, b trust.ArtifactRecord) int {
		return norm.CompareStrings(a.Path, b.Path, a.Kind, b.Kind, a.MediaType, b.MediaType)
	})
	return out
}
