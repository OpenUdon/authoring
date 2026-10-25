# Retired milestone M12 - Agent Mode Result Contracts

**Milestone.** M12
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M12.md
**Source specification.** tabilet/memory-bank/milestone.md#m12-agent-mode-result-contracts
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
## M12 Agent Mode Result Contracts

### M12 Agent Mode Result Contracts

**Goal.** Add generic noninteractive result shapes.

Acceptance:

- Result contracts cover `complete`, `needs_input`, `failed`, and `canceled`.
- Results include top readiness issue, diagnostics, generated artifact
  descriptors, transcript metadata, and digest metadata.
- Tests cover deterministic JSON, missing input summaries, cancellation, and
  failure diagnostics.

Delivered:

- `report` defines `authoring.agent-result.v1` with generic statuses,
  readiness/top issue metadata, M11 repair status fields, M09 decision behavior
  summaries, diagnostics, artifact descriptors, digest records, and
  session/transcript metadata.
- Canonical JSON and normalization keep result output deterministic for
  downstream agent-mode tests.

````

## Status record

````markdown
# Status M12

## Goal

Add generic noninteractive agent-mode result contracts.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define result statuses | `[+]` | Added `complete`, `needs_input`, `failed`, and `canceled` statuses in `report.Result`. |
| Add result metadata | `[+]` | Added M10 readiness/top issue metadata, M11 repair status fields, M09 decision behavior summaries, diagnostics, artifact descriptors, M03 session/transcript metadata, and M02 `trust` digest records. |
| Stabilize JSON output | `[+]` | Added `Normalize` and `CanonicalJSON` for deterministic `authoring.agent-result.v1` output. |
| Add result tests | `[+]` | Covered missing input summaries, cancellation/deadline status, failure diagnostics, artifact/digest metadata, transcript metadata, deterministic JSON, and decision-summary preservation. |

## Notes

- Result contracts should not imply product-specific remediation or execution.
- Deep review finding fixed before commit: confirmation-only decision summaries
  are preserved during normalization.

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
