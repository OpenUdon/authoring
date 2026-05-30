package promptcontext

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/OpenUdon/authoring/trust"
)

const Version = "authoring.prompt-context.v1"

// Context is a deterministic prompt-safe bundle of downstream metadata.
type Context struct {
	Version     string               `json:"version"`
	Sources     []SourceDocument     `json:"sources,omitempty"`
	Operations  []OperationCandidate `json:"operations,omitempty"`
	Schemas     []SchemaHint         `json:"schemas,omitempty"`
	Credentials []CredentialBinding  `json:"credentials,omitempty"`
	Metadata    map[string]string    `json:"metadata,omitempty"`
}

// SourceDocument is a prompt-safe source document summary.
type SourceDocument struct {
	ID        string             `json:"id"`
	Kind      string             `json:"kind,omitempty"`
	Title     string             `json:"title,omitempty"`
	Version   string             `json:"version,omitempty"`
	URI       string             `json:"uri,omitempty"`
	MediaType string             `json:"media_type,omitempty"`
	Summary   string             `json:"summary,omitempty"`
	Digest    trust.DigestRecord `json:"digest,omitempty"`
	Metadata  map[string]string  `json:"metadata,omitempty"`
}

// OperationCandidate is a prompt-safe operation candidate summary.
type OperationCandidate struct {
	ID                 string            `json:"id"`
	SourceID           string            `json:"source_id,omitempty"`
	OperationID        string            `json:"operation_id,omitempty"`
	Name               string            `json:"name,omitempty"`
	Verb               string            `json:"verb,omitempty"`
	Path               string            `json:"path,omitempty"`
	Summary            string            `json:"summary,omitempty"`
	RequestSchemaID    string            `json:"request_schema_id,omitempty"`
	ResponseSchemaID   string            `json:"response_schema_id,omitempty"`
	CredentialBindings []string          `json:"credential_bindings,omitempty"`
	Tags               []string          `json:"tags,omitempty"`
	Confidence         string            `json:"confidence,omitempty"`
	SelectionRationale string            `json:"selection_rationale,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

// SchemaHint is a prompt-safe schema summary, not a product schema contract.
type SchemaHint struct {
	ID        string            `json:"id"`
	Name      string            `json:"name,omitempty"`
	Purpose   string            `json:"purpose,omitempty"`
	MediaType string            `json:"media_type,omitempty"`
	Summary   string            `json:"summary,omitempty"`
	Required  []string          `json:"required,omitempty"`
	Fields    []FieldHint       `json:"fields,omitempty"`
	Examples  []string          `json:"examples,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// FieldHint is a prompt-safe field summary.
type FieldHint struct {
	Name      string `json:"name"`
	Type      string `json:"type,omitempty"`
	Required  bool   `json:"required,omitempty"`
	Sensitive bool   `json:"sensitive,omitempty"`
	Summary   string `json:"summary,omitempty"`
}

// CredentialBinding names a symbolic credential binding. It must not contain
// credential values.
type CredentialBinding struct {
	Name     string            `json:"name"`
	Kind     string            `json:"kind,omitempty"`
	Scope    string            `json:"scope,omitempty"`
	Required bool              `json:"required,omitempty"`
	Summary  string            `json:"summary,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Normalize returns a deterministic prompt-safe context.
func Normalize(ctx Context) Context {
	ctx.Version = firstNonEmpty(strings.TrimSpace(ctx.Version), Version)
	ctx.Sources = NormalizeSources(ctx.Sources)
	ctx.Operations = NormalizeOperations(ctx.Operations)
	ctx.Schemas = NormalizeSchemas(ctx.Schemas)
	ctx.Credentials = NormalizeCredentials(ctx.Credentials)
	ctx.Metadata = normalizeMetadata(ctx.Metadata)
	return ctx
}

// CanonicalJSON returns deterministic indented JSON for ctx.
func CanonicalJSON(ctx Context) ([]byte, error) {
	return json.MarshalIndent(Normalize(ctx), "", "  ")
}

// NormalizeSources returns deterministic source document summaries.
func NormalizeSources(sources []SourceDocument) []SourceDocument {
	out := make([]SourceDocument, 0, len(sources))
	for _, source := range sources {
		source.ID = strings.TrimSpace(source.ID)
		source.Kind = normalizeToken(source.Kind)
		source.Title = redact(source.Title)
		source.Version = strings.TrimSpace(source.Version)
		source.URI = redact(source.URI)
		source.MediaType = strings.TrimSpace(source.MediaType)
		source.Summary = redact(source.Summary)
		source.Digest = normalizeDigest(source.Digest)
		source.Metadata = normalizeMetadata(source.Metadata)
		if source.ID == "" {
			continue
		}
		out = append(out, source)
	}
	slices.SortStableFunc(out, func(a, b SourceDocument) int {
		return compareStrings(a.ID, b.ID, a.Kind, b.Kind, a.Title, b.Title)
	})
	return out
}

// NormalizeOperations returns deterministic operation candidate summaries.
func NormalizeOperations(operations []OperationCandidate) []OperationCandidate {
	out := make([]OperationCandidate, 0, len(operations))
	for _, operation := range operations {
		operation.ID = strings.TrimSpace(operation.ID)
		operation.SourceID = strings.TrimSpace(operation.SourceID)
		operation.OperationID = strings.TrimSpace(operation.OperationID)
		operation.Name = redact(operation.Name)
		operation.Verb = normalizeVerb(operation.Verb)
		operation.Path = redact(operation.Path)
		operation.Summary = redact(operation.Summary)
		operation.RequestSchemaID = strings.TrimSpace(operation.RequestSchemaID)
		operation.ResponseSchemaID = strings.TrimSpace(operation.ResponseSchemaID)
		operation.CredentialBindings = normalizeList(operation.CredentialBindings, false)
		operation.Tags = normalizeList(operation.Tags, true)
		operation.Confidence = normalizeToken(operation.Confidence)
		operation.SelectionRationale = redact(operation.SelectionRationale)
		operation.Metadata = normalizeMetadata(operation.Metadata)
		if operation.ID == "" {
			continue
		}
		out = append(out, operation)
	}
	slices.SortStableFunc(out, func(a, b OperationCandidate) int {
		return compareStrings(a.ID, b.ID, a.SourceID, b.SourceID, a.OperationID, b.OperationID)
	})
	return out
}

// NormalizeSchemas returns deterministic schema hints.
func NormalizeSchemas(schemas []SchemaHint) []SchemaHint {
	out := make([]SchemaHint, 0, len(schemas))
	for _, schema := range schemas {
		schema.ID = strings.TrimSpace(schema.ID)
		schema.Name = redact(schema.Name)
		schema.Purpose = normalizeToken(schema.Purpose)
		schema.MediaType = strings.TrimSpace(schema.MediaType)
		schema.Summary = redact(schema.Summary)
		schema.Required = normalizeList(schema.Required, false)
		schema.Fields = NormalizeFields(schema.Fields)
		schema.Examples = normalizeRedactedList(schema.Examples)
		schema.Metadata = normalizeMetadata(schema.Metadata)
		if schema.ID == "" {
			continue
		}
		out = append(out, schema)
	}
	slices.SortStableFunc(out, func(a, b SchemaHint) int {
		return compareStrings(a.ID, b.ID, a.Name, b.Name)
	})
	return out
}

// NormalizeFields returns deterministic field hints.
func NormalizeFields(fields []FieldHint) []FieldHint {
	out := make([]FieldHint, 0, len(fields))
	for _, field := range fields {
		field.Name = strings.TrimSpace(field.Name)
		field.Type = normalizeToken(field.Type)
		field.Summary = redact(field.Summary)
		if field.Name == "" {
			continue
		}
		out = append(out, field)
	}
	slices.SortStableFunc(out, func(a, b FieldHint) int {
		return compareStrings(a.Name, b.Name, a.Type, b.Type)
	})
	return out
}

// NormalizeCredentials returns deterministic symbolic credential bindings.
func NormalizeCredentials(credentials []CredentialBinding) []CredentialBinding {
	out := make([]CredentialBinding, 0, len(credentials))
	for _, credential := range credentials {
		credential.Name = strings.TrimSpace(credential.Name)
		credential.Kind = normalizeToken(credential.Kind)
		credential.Scope = normalizeToken(credential.Scope)
		credential.Summary = redact(credential.Summary)
		credential.Metadata = normalizeMetadata(credential.Metadata)
		if credential.Name == "" {
			continue
		}
		out = append(out, credential)
	}
	slices.SortStableFunc(out, func(a, b CredentialBinding) int {
		return compareStrings(a.Name, b.Name, a.Kind, b.Kind, a.Scope, b.Scope)
	})
	return out
}

func normalizeList(values []string, token bool) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if token {
			value = normalizeToken(value)
		} else {
			value = strings.TrimSpace(value)
		}
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	slices.Sort(out)
	return out
}

func normalizeRedactedList(values []string) []string {
	out := make([]string, len(values))
	for i := range values {
		out[i] = redact(values[i])
	}
	return normalizeList(out, false)
}

func normalizeMetadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		key = strings.TrimSpace(key)
		value = redactMetadataValue(key, value)
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

func redactMetadataValue(key, value string) string {
	redacted, _ := trust.RedactDocument(map[string]string{key: value})
	if values, ok := redacted.(map[string]string); ok {
		return strings.TrimSpace(values[key])
	}
	return redact(value)
}

func normalizeDigest(record trust.DigestRecord) trust.DigestRecord {
	record.Algorithm = normalizeToken(record.Algorithm)
	record.Value = strings.TrimSpace(record.Value)
	return record
}

func normalizeVerb(verb string) string {
	verb = strings.TrimSpace(verb)
	if verb == "" {
		return ""
	}
	return strings.ToUpper(verb)
}

func redact(value string) string {
	return trust.RedactString(strings.TrimSpace(value))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), "_"))
}

func compareStrings(values ...string) int {
	for i := 0; i+1 < len(values); i += 2 {
		if values[i] < values[i+1] {
			return -1
		}
		if values[i] > values[i+1] {
			return 1
		}
	}
	return 0
}
