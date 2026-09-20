# Status M15

## Goal

Define prompt-safe context interfaces for product-neutral metadata.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define source document shapes | `[+]` | Added prompt-safe source document identity, kind, title, URI, media type, summary, digest, and metadata fields. |
| Define operation candidate shapes | `[+]` | Added operation candidates with source/operation IDs, verb/path, schema references, symbolic credential binding names, tags, confidence, and rationale without API-source imports. |
| Define schema and credential binding shapes | `[+]` | Added schema/field hints and symbolic credential binding records with redaction-backed normalization. |
| Add import-boundary tests | `[+]` | Added build-info import-boundary test and retained external `go list -deps` scan for apitools, UWS, OpenUdon, and Ramen. |

## Notes

- Downstream adapters translate real API metadata into these shapes.
- Deep review findings fixed before commit: metadata values now use
  key-sensitive redaction, URI/path text is redacted, and example normalization
  no longer mutates caller-owned slices.

## Verification

Completed checks:

```bash
go test ./...
go vet ./...
GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off go list -deps ./... | rg 'github.com/OpenUdon/(openudon|ramen|uws|apitools)' || true
git diff --check
git -C ../tofu diff --check -- authoring
```
