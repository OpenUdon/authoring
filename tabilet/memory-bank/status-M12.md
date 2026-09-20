# Status M12

## Goal

Add generic noninteractive agent-mode result contracts.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define result statuses | `[+]` | Added `complete`, `needs_input`, `failed`, and `canceled` statuses in `report.Result`. |
| Add result metadata | `[+]` | Added M10 readiness/top issue metadata, M11 repair status fields, M09 decision behavior summaries, diagnostics, artifact descriptors, M03 session/transcript metadata, and M02 `trust` digest records. |
| Stabilize JSON output | `[+]` | Added `Normalize` and `CanonicalJSON` for deterministic `authoring.agent-result.v1` output. |
| Add result tests | `[+]` | Covered missing input summaries, cancellation/deadline status, failure diagnostics, artifact/digest metadata, transcript metadata, deterministic JSON, and decision-summary preservation. |

## Notes

- Result contracts should not imply product-specific remediation or execution.
- Deep review finding fixed before commit: confirmation-only decision summaries
  are preserved during normalization.

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
