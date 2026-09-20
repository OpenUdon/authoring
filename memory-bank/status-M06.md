# Status M06

## Goal

Extract provider-neutral structured JSON output helpers.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define structured client interface | `[+]` | Added provider-neutral `Client` and `StructuredClient` interfaces plus result/options types. |
| Normalize schemas | `[+]` | Added deterministic schema normalization from raw JSON, strings, bytes, raw messages, and marshalable values. |
| Add JSON fallback parsing | `[+]` | Added structured-first completion, legacy instruction fallback, JSON object/array block extraction, and decode helpers. |
| Add fake-client tests | `[+]` | Tests cover structured success, structured error fallback, legacy extraction, invalid JSON, nil client, nil output target, and schema normalization. |

## Notes

- OpenAI, Anthropic, Gemini, and any other provider clients stay downstream.
- Review pass confirmed no provider-specific client imports and no OpenUdon,
  Ramen, UWS, or apitools imports.

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
