# Retired milestone M17 - OpenUdon Adoption Phase 2

**Milestone.** M17
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M17.md
**Source specification.** tabilet/memory-bank/milestone.md#m17-openudon-adoption-phase-2
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
## M17 OpenUdon Adoption Phase 2

### M17 OpenUdon Adoption Phase 2

**Goal.** Move OpenUdon iCoT agent/report/scorecard shared mechanics onto
Authoring.

Acceptance:

- OpenUdon uses Authoring for shared agent result, report, retention,
  scorecard, and variant mechanics.
- OpenUdon-specific session schema, prompts, workflow intent, catalog planning,
  repair rules, and package artifacts remain in OpenUdon.
- OpenUdon fixture outcomes and wire versions remain unchanged unless a
  deliberate migration note records the change.

````

## Status record

````markdown
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
````
