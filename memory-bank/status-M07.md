# Status M07

## Goal

Move the generic progressive iCoT loop shape into Authoring.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add progressive loop | `[+]` | Added `icot.Run` with draft attempts, readiness decisions, question planning/application, autosave hooks, final confirmation, and bounded attempts. |
| Record events | `[+]` | Emits ordered M03 transcript events for readiness, draft attempt/success/error, questions, needs-input, final confirmation, and final result summaries. |
| Support defaulted answers | `[+]` | Question plans can allow defaulted answers with source evidence while forced questions still require operator input. |
| Add fake-runtime tests | `[+]` | Tests cover success, `needs_input`, cancellation, draft errors, defaulted answers, forced questions, bounded attempts, and event order. |

## Notes

- Draft schemas and product issue codes remain downstream.
- Review findings fixed before completion: result events preserve loop order
  instead of using durable transcript sorting, and cancellation wraps both
  Authoring's `ErrCanceled` and the underlying context error.

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
