# Status M11

## Goal

Add product-neutral draft review and repair orchestration.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define review issue containers | `[+]` | Added `ReviewIssue`, `Remediation`, and deterministic normalization while reusing M10 `readiness.Issue`/`Result`. |
| Orchestrate repair attempts | `[+]` | Added bounded `RunRepair` and `RunRuntimeRepair` over the M08 review/repair runtime hook shapes. |
| Record repair events | `[+]` | Added transcript events for review, repair attempt/success/error, exhaustion, no-op repair, and review failure. |
| Add fake-runtime tests | `[+]` | Covered repair success, exhaustion, review errors, no-op repairs, runtime binding, review issue normalization, and event order. |

## Notes

- Downstream products own repair rules and draft mutation.
- Deep review finding fixed before commit: review error events now include the
  current attempt, and top issue event fields omit empty values.

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
