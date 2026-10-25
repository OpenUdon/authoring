# Retired milestone M21 - Public API Stabilization

**Milestone.** M21
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M21.md
**Source specification.** tabilet/memory-bank/milestone.md#m21-public-api-stabilization
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
## M21 Public API Stabilization

### M21 Public API Stabilization

**Goal.** Review exported package names, version constants, JSON tags, and
compatibility policy before broader downstream reliance.

Acceptance:

- Public APIs are marked as pre-1.0 unstable where appropriate.
- Version constants and JSON tags are reviewed for durable contracts.
- Migration expectations for OpenUdon and Ramen adopters are documented.
- Machine checks cover durable `authoring.*.v1` constants and exported JSON
  tags on durable records.

````

## Status record

````markdown
# Status M21

## Goal

Stabilize the pre-1.0 public API surface.

## Tasks

| Item | State | Notes |
|---|---|---|
| Review exported package names | `[+]` | `README.md` and `COMPATIBILITY.md` document behavior-grouped packages and confirm product adapters stay downstream. |
| Review version constants and JSON tags | `[+]` | Added `api_surface_test.go` coverage for durable `authoring.*.v1` constants and exported JSON tags on durable record types. |
| Document compatibility policy | `[+]` | Added `COMPATIBILITY.md` with pre-1.0 instability, version/tag migration rules, and downstream adapter expectations. |
| Check downstream migrations | `[+]` | Compatibility policy and gate document OpenUdon/Ramen migration expectations and require dependent workspace checks after exported API changes. |

## Notes

- This is an API review milestone, not a feature expansion milestone.
- No product wire versions were changed.

## Verification

Completed checks:

```bash
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
(cd ../authoring && git diff --check)
git -C ../tofu diff --check -- authoring
```
````
