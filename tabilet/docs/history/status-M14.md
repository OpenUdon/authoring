# Retired milestone M14 - Scorecard And Variant Harness Primitives

**Milestone.** M14
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M14.md
**Source specification.** tabilet/memory-bank/milestone.md#m14-scorecard-and-variant-harness-primitives
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
## M14 Scorecard And Variant Harness Primitives

### M14 Scorecard And Variant Harness Primitives

**Goal.** Add generic fixture, variant, and scorecard result primitives.

Acceptance:

- Public structs cover fixture/variant results, expected outcome comparison,
  grouped summary counters, failure-family fields, and report validation
  helpers.
- OpenUdon keeps its corpus and `openudon.icot-*` wire versions.
- Tests cover grouping, comparison, failure-family classification, and stable
  summaries.

Delivered:

- `report` now includes `authoring.scorecard.v1` scorecards, fixture/variant
  result records, outcome normalization, expected/observed comparison, grouped
  counters, failure-family summaries, and canonical scorecard JSON.
- Generic validation helpers return product-neutral diagnostics for malformed
  scorecard shapes without importing OpenUdon corpora or wire schemas.

````

## Status record

````markdown
# Status M14

## Goal

Add generic scorecard and variant harness primitives.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define fixture and variant results | `[+]` | Added `FixtureResult` and `VariantResult` with fixture/variant identity, expected outcome, observed outcome, diagnostics, and metadata. |
| Add grouped summaries | `[+]` | Added scorecard summaries with outcome counters, group summaries, expected comparison counters, and failure-family fields. |
| Add report validation helpers | `[+]` | Added generic scorecard validation diagnostics without OpenUdon corpus imports. |
| Add comparison tests | `[+]` | Covered grouping, expected outcome comparison, failure families, fixture propagation, validation diagnostics, and stable JSON. |

## Notes

- OpenUdon keeps its corpus and `openudon.icot-*` wire versions.
- Deep review finding fixed before commit: `outcome_matched` is always emitted
  because false mismatches are meaningful scorecard data.

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
