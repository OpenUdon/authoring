# Status M02

## Goal

Add Evidence-backed neutral trust primitives before Authoring defines durable
session, transcript, report, or artifact contracts.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add Evidence dependency | `[+]` | `go.mod` requires `github.com/OpenUdon/evidence` with a sibling checkout replace while the module has no public tag. |
| Use neutral diagnostic records | `[+]` | `trust.DiagnosticRecord` and helpers alias/delegate to `evidence/diagnostic`. |
| Use redaction helpers | `[+]` | `trust.RedactString` and `trust.RedactDocument` delegate to `evidence/redact`. |
| Use artifact and digest records | `[+]` | `trust` aliases artifact/digest records and wraps safe paths, manifests, and SHA-256 byte digests. |
| Document trust boundary | `[+]` | Product, architecture, tech-stack, and this status file state that Authoring does not redefine durable trust evidence. |

## Notes

- Product-specific approval, governance, executor, and package evidence stays
  in OpenUdon or Ramen.
- M02 does not require OpenUdon or Ramen code changes; those begin in later
  adoption milestones.
- Review finding fixed before completion: an invalid local Evidence
  pseudo-version was replaced with `v0.0.0` plus a sibling checkout replace.

## Verification

Completed checks:

```bash
go test ./...
go vet ./...
git diff --check
git -C ../tofu diff --check -- authoring
```
