# Retired milestone M16 - OpenUdon Adoption Phase 1

**Milestone.** M16
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M16.md
**Source specification.** tabilet/memory-bank/milestone.md#m16-openudon-adoption-phase-1
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
## M16 OpenUdon Adoption Phase 1

### M16 OpenUdon Adoption Phase 1

**Goal.** Replace OpenUdon's generic prompt, transcript, structured-output,
lifecycle, and progressive-loop helpers with Authoring packages.

Acceptance:

- OpenUdon CLI output, JSON versions, fixtures, and package behavior are
  preserved.
- OpenUdon prompts, workflow intent, catalog planning, repair rules, and
  review packages remain in OpenUdon.
- `go test ./...`, focused iCoT tests, and `make icot-authoring-scorecard`
  pass in `../openudon`.

````

## Status record

````markdown
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
````
