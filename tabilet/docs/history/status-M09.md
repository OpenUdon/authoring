# Retired milestone M09 - Decision Evidence And Confidence Policy

**Milestone.** M09
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M09.md
**Source specification.** tabilet/memory-bank/milestone.md#m09-decision-evidence-and-confidence-policy
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
## M09 Decision Evidence And Confidence Policy

### M09 Decision Evidence And Confidence Policy

**Goal.** Add generic decision evidence and confidence handling.

Acceptance:

- Decision evidence records stage, slot, value, source, confidence, rationale,
  alternatives, and confirmation requirement.
- Confidence behavior distinguishes auto-accept, review, low-confidence, and
  conflict states.
- Tests cover sorting, conflict handling, redaction posture, and transcript
  event binding.

````

## Status record

````markdown
# Status M09

## Goal

Add generic decision evidence and confidence behavior.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define decision evidence | `[+]` | Added `decision.Record` and session aliases with stage, slot, value, source, confidence, rationale, evidence, alternatives, and confirmation requirement. |
| Define confidence policy | `[+]` | Added generic auto-accept, review, low-confidence, and conflict behavior plus confirmation checks. |
| Bind decisions to transcript | `[+]` | Added transcript decision-event helpers that emit product-neutral `decision` events. |
| Add sorting and conflict tests | `[+]` | Tests cover deterministic normalization, merge/conflict handling, confidence behavior, transcript event binding, and redaction posture. |

## Notes

- Product-specific confidence thresholds may be downstream policy layered on
  top of generic states.
- Review pass confirmed session integration preserved existing M03 decision
  JSON shape through type aliases.

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
