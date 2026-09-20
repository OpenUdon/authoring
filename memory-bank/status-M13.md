# Status M13

## Goal

Add shared report and retention metadata.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define report metadata | `[+]` | Added `ReportMetadata` with run ID, command, commit, generation time, retention class, provider-output flag, archive safety, and redaction-required flag. |
| Add digest sidecars | `[+]` | Added digest sidecars for report bytes and M02 artifact records using `trust` digest aliases and safe path helpers. |
| Propagate redaction requirements | `[+]` | Added redaction checks through M02 `trust.RedactDocument`; redaction-required metadata clears archive-safe status. |
| Add report tests | `[+]` | Covered retention normalization, digest sidecars, provider-output/archive flags, unsafe sidecar paths, and redaction-required propagation. |

## Notes

- Product-specific report bodies and review artifacts remain downstream.
- Deep review finding fixed before commit: digest sidecar path normalization now
  uses the shared safe artifact path cleaner and drops unsafe paths.

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
