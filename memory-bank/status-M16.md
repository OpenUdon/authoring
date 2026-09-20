# Status M16

## Goal

Adopt Authoring in OpenUdon for generic prompt, transcript, structured-output,
lifecycle, and progressive-loop helpers.

## Tasks

| Item | State | Notes |
|---|---|---|
| Replace generic prompt helpers | `[+]` | OpenUdon `internal/authoring.PromptSession` now wraps M04 `prompt` primitives while preserving prompt wording, default modes, replay labels, and CLI output. |
| Replace transcript and lifecycle helpers | `[+]` | OpenUdon keeps its `openudon.*` transcript wire versions, but local transcript and draft writes now route through M05 lifecycle atomic-write behavior. |
| Replace structured-output helpers | `[+]` | OpenUdon structured JSON completion, legacy fallback, JSON block extraction, decode, and schema normalization now delegate to M06 `structured`; provider clients stay downstream. |
| Replace progressive-loop helpers | `[+]` | Phase 1 keeps OpenUdon's product-specific loop sequencing, prompts, catalog planning, and workflow-intent mutation downstream, while the loop uses Authoring prompt/lifecycle/structured mechanics and preserved scorecard behavior. |
| Translate prompt context | `[+]` | Added an OpenUdon adapter from API-source operation/security/request metadata to M15 `promptcontext` records and passes that context into draft prompt payloads. |

## Notes

- OpenUdon prompts, workflow intent, catalog planning, repair rules, package
  quality, review packages, and handoff semantics stay in OpenUdon.
- OpenUdon adoption currently relies on the parent `go.work` sibling checkout
  for `github.com/OpenUdon/authoring`; no committed `replace ../...`
  directive was added. M20 must close public-module compatibility once
  Authoring has a consumable version.
- The full public `icot` runtime loop remains available in Authoring, but
  OpenUdon's richer product loop is still downstream because it owns catalog
  planning, selected-operation drafting, final edit/explain confirmation, and
  workflow-intent-specific mutation.

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
