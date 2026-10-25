# Retired milestone M07 - Generic Progressive iCoT Loop

**Milestone.** M07
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M07.md
**Source specification.** tabilet/memory-bank/milestone.md#m07-generic-progressive-icot-loop
**Evidence.** be023f06d0837b58a651449b6d89a00c12eb04cd
**Worktree.** includes uncommitted changes
**Review.** legacy
**Review iterations.** not recorded
**Verification.** Historical verification and its limitations are retained verbatim in the status record; original commands are not re-run or relabelled by this retirement. At the observed evidence revision, go test ./... and go vet ./... passed on 2026-10-25; literal source preservation and retirement-envelope validation passed before removal.
**Consolidated into.** Current product, architecture, stack, and lessons remain in tabilet/memory-bank/; no implementation or current product contract changes. Superseded roadmap and retention wording is preserved in knowledge.md.
**Specification encoding.** The added level-two identity heading is a retirement-parser wrapper; the following level-three section is the complete original source, retained verbatim.
**Legacy closure.** This milestone closed before adoption of the persisted review gate. Its original closure evidence is the complete literal status record; no review iteration or retrospective implementation review is invented.

## Milestone specification

````markdown
## M07 Generic Progressive iCoT Loop

### M07 Generic Progressive iCoT Loop

**Goal.** Move OpenUdon's generic progressive loop shape into Authoring.

Acceptance:

- The loop supports runtime hooks, draft attempts, readiness decisions,
  question planning, defaulted answers, transcript events, final confirmation
  callbacks, and bounded attempts.
- Tests cover success, `needs_input`, cancellation, draft errors, forced
  questions, defaulted answers, and deterministic event order.
- Product-specific draft schemas and issue codes stay downstream.

````

## Status record

````markdown
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
````
