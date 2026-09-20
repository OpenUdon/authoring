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
