# Status M17

## Goal

Move OpenUdon iCoT agent, report, scorecard, and variant shared mechanics onto
Authoring.

## Tasks

| Item | State | Notes |
|---|---|---|
| Adopt agent result contracts | `[+]` | OpenUdon `openudon.icot-author-report.v1` output is unchanged, but author reports are validated through M12 `report.Result` with M10 readiness issue projection. |
| Adopt report metadata helpers | `[+]` | OpenUdon scorecard and authoring-eval retention flags are mapped through M13 `report.ReportMetadata`; OpenUdon-specific report bodies stay downstream. |
| Adopt scorecard primitives | `[+]` | OpenUdon scorecard validation now builds an M14 `report.Scorecard` contract while preserving corpus outcomes and `openudon.icot-scorecard.v1`. |
| Run regression corpus | `[+]` | Focused report/scorecard contract tests pass; full OpenUdon and Authoring gates are recorded below. |

## Notes

- OpenUdon-specific session schema, prompts, workflow intent, catalog
  planning, repair rules, and package artifacts remain in OpenUdon.
- OpenUdon `openudon.icot-*` wire versions are intentionally unchanged. The
  Authoring contracts are used as compatibility/normalization gates, not as
  replacement JSON schemas in this phase.

## Verification

Completed checks:

```bash
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../authoring && GOWORK=off go list -deps ./... | rg 'github.com/OpenUdon/(openudon|ramen|uws|apitools)' || true)
(cd ../openudon && go test ./... && go vet ./...)
(cd ../openudon && make icot-authoring-scorecard)
(cd ../authoring && git diff --check)
(cd ../openudon && git diff --check)
git -C ../tofu diff --check -- authoring openudon
```
