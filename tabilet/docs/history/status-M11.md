# Retired milestone M11 - Draft Review And Repair Orchestration

**Milestone.** M11
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M11.md
**Source specification.** tabilet/memory-bank/milestone.md#m11-draft-review-and-repair-orchestration
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
## M11 Draft Review And Repair Orchestration

### M11 Draft Review And Repair Orchestration

**Goal.** Add product-neutral review issue/remediation containers and bounded
repair orchestration.

Acceptance:

- Authoring owns attempt limits, event recording, repair summaries, and final
  failure classification.
- Downstream runtimes own repair rules and draft mutation.
- Tests cover repair success, repair exhaustion, review errors, no-op repairs,
  and transcript events.

Delivered:

- `icot` defines review/remediation containers, repair statuses, bounded
  `RunRepair`, and `RunRuntimeRepair` over the M08 review/repair runtime hook
  shapes.
- Repair orchestration records M03 transcript events for review, repair
  attempt/success/error, no-op, and exhaustion while using M10 readiness
  summaries for issue counts and top issue fields.

````

## Status record

````markdown
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
````
