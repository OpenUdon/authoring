package promptcontext

import (
	"runtime/debug"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/trust"
)

func TestNormalizeContext(t *testing.T) {
	ctx := Normalize(Context{
		Sources: []SourceDocument{
			{ID: " openapi ", Kind: " OpenAPI ", Title: " API ", Digest: trust.DigestRecord{Algorithm: " SHA256 ", Value: " abc "}},
			{},
		},
		Operations: []OperationCandidate{
			{
				ID:                 " create ",
				SourceID:           "openapi",
				Verb:               "post",
				Path:               "/v1/items",
				CredentialBindings: []string{" prod ", "prod"},
				Tags:               []string{" Items ", "items"},
				SelectionRationale: "uses token=sk-proj-abcdefghijklmnopqrstuvwxyz",
			},
		},
		Schemas: []SchemaHint{
			{
				ID:       " create-request ",
				Purpose:  " Request Body ",
				Required: []string{"name", "name"},
				Fields: []FieldHint{
					{Name: " secret ", Type: " String ", Sensitive: true, Summary: "api_key='x'"},
					{Name: "name", Type: "String"},
				},
			},
		},
		Credentials: []CredentialBinding{
			{Name: " prod ", Kind: " API Key ", Scope: " Runtime ", Summary: "Bearer abcdefghijklmnopqrstuvwxyz", Metadata: map[string]string{"api_token": "plain"}},
		},
	})
	if ctx.Version != Version || len(ctx.Sources) != 1 || ctx.Sources[0].Kind != "openapi" {
		t.Fatalf("context = %#v, want normalized source", ctx)
	}
	if got := ctx.Operations[0].CredentialBindings; len(got) != 1 || got[0] != "prod" {
		t.Fatalf("credential bindings = %#v, want deduped symbolic names", got)
	}
	if ctx.Operations[0].Verb != "POST" || ctx.Operations[0].Tags[0] != "items" {
		t.Fatalf("operation = %#v, want normalized operation", ctx.Operations[0])
	}
	if ctx.Operations[0].SelectionRationale == "uses token=sk-proj-abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("selection rationale was not redacted")
	}
	if ctx.Credentials[0].Summary == "Bearer abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("credential summary was not redacted")
	}
	if ctx.Credentials[0].Metadata["api_token"] != trust.RedactedValue {
		t.Fatalf("credential metadata = %#v, want key-sensitive redaction", ctx.Credentials[0].Metadata)
	}
}

func TestCanonicalJSON(t *testing.T) {
	data, err := CanonicalJSON(Context{Operations: []OperationCandidate{{ID: "op", Verb: "get"}}})
	if err != nil {
		t.Fatalf("CanonicalJSON error = %v", err)
	}
	if !strings.Contains(string(data), `"version": "authoring.prompt-context.v1"`) || !strings.Contains(string(data), `"verb": "GET"`) {
		t.Fatalf("json = %s", data)
	}
}

func TestImportBoundary(t *testing.T) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		t.Fatal("build info unavailable")
	}
	blocked := []string{
		"github.com/OpenUdon/apitools",
		"github.com/OpenUdon/uws",
		"github.com/OpenUdon/openudon",
		"github.com/OpenUdon/ramen",
	}
	for _, dep := range info.Deps {
		for _, prefix := range blocked {
			if dep.Path == prefix || strings.HasPrefix(dep.Path, prefix+"/") {
				t.Fatalf("blocked dependency %s found in build info", dep.Path)
			}
		}
	}
}
