# Status M03

## Goal

Define public session and transcript contracts for durable generic authoring
records.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add `session` package | `[+]` | Added versioned session state, prompt turns, answers, readiness issues, decisions, artifact summaries, diagnostics, metadata, normalization, and canonical JSON. |
| Add `transcript` package | `[+]` | Added versioned transcript records, turns, events, model/provider provenance, diagnostics, artifact summaries, metadata, normalization, and canonical JSON. |
| Normalize deterministically | `[+]` | Session and transcript records trim fields, normalize tokens, sort records, normalize M02 `trust` diagnostics/artifacts, and produce stable JSON. |
| Add mapping tests | `[+]` | Tests map OpenUdon-shaped and Ramen-shaped structs into generic records without importing either product. |

## Notes

- Keep provider clients and prompt text downstream.
- Review finding fixed before completion: optional transcript provider
  provenance uses pointers so empty provider records do not serialize as `{}`.

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
