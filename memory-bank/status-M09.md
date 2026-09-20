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
