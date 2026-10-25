# Retired milestone M10 - Readiness And Question Planning

**Milestone.** M10
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M10.md
**Source specification.** tabilet/memory-bank/milestone.md#m10-readiness-and-question-planning
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
## M10 Readiness And Question Planning

### M10 Readiness And Question Planning

**Goal.** Add generic readiness/result primitives and question planning
contracts.

Acceptance:

- Readiness and question planning support blocking and warning issues,
  suggested answers, forced-question mechanics, and stable sorting.
- Product-specific issue codes and remediation language remain downstream.
- Tests cover top issue selection, deterministic ordering, forced questions,
  and safe defaults.

Delivered:

- `readiness` defines deterministic issue normalization, readiness summaries,
  top issue selection, decision-confirmation issue projection, and question
  plan/default helpers.
- `icot.Question` aliases `readiness.Question`, and the default iCoT readiness
  policy treats blocking findings as blockers while allowing warning-only
  readiness to proceed unless a downstream `Ready` hook overrides it.

````

## Status record

````markdown
# Status M10

## Goal

Add generic readiness and question planning contracts.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define readiness results | `[+]` | Added `readiness.Result`, `Summary`, deterministic issue normalization, blocking/warning subsets, top issue selection, and M09 decision-confirmation issue projection. |
| Define question planning | `[+]` | Added `readiness.Question`, `Plan`, suggested/forced question builders, and safe default-answer detection. |
| Add deterministic sorting | `[+]` | Added issue and question comparators with stable readiness and question-plan evaluation. |
| Add focused tests | `[+]` | Covered blocking/warning handling, warning-only readiness, top issue selection, decision confirmation issues, forced questions, defaults, and plan ordering. |

## Notes

- Product-specific issue codes and remediation wording stay downstream.
- `icot.Question` now aliases `readiness.Question`, preserving the M07/M08 loop
  surface while moving the public question contract into the M10 package.
- Deep review findings fixed before commit: `SuggestedQuestion` now handles an
  empty issue without panicking, and unknown readiness severities block by
  default.

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
