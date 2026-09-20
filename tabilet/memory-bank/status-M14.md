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
