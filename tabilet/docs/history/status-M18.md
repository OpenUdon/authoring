# Retired milestone M18 - Ramen Adapter Spike

**Milestone.** M18
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M18.md
**Source specification.** tabilet/memory-bank/milestone.md#m18-ramen-adapter-spike
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
## M18 Ramen Adapter Spike

### M18 Ramen Adapter Spike

**Goal.** Add a Ramen downstream adapter spike that uses Authoring to draft
native Ramen/UWS project skeletons from an operator goal.

Acceptance:

- The spike can produce a draft `project.uws.yaml` skeleton.
- The generated project validates when enough metadata is present.
- The adapter stops before graph/plan if required mapping metadata is missing.
- The adapter is downstream in Ramen; Authoring still imports no Ramen, UWS, or
  OpenUdon packages.

````

## Status record

````markdown
# Status M18

## Goal

Add a Ramen adapter spike that uses Authoring to draft native Ramen/UWS project
skeletons from an operator goal.

## Tasks

| Item | State | Notes |
|---|---|---|
| Bind Ramen fake/runtime adapter | `[+]` | Added `github.com/OpenUdon/ramen/authoring` adapter code that binds Ramen state to M08 runtime interfaces without importing OpenUdon internals. |
| Draft project skeletons | `[+]` | Uses Authoring runtime hooks and M15 `promptcontext` records to generate native `project.uws.yaml` skeletons from goals. |
| Validate when metadata exists | `[+]` | Runs Ramen `validate.Run` for complete skeletons and surfaces validation state in the M12 report result. |
| Stop before graph/plan when needed | `[+]` | Returns M12 `needs_input` results when required operation mapping metadata is missing, before graph or plan execution. |

## Notes

- Ramen desired-state/project semantics remain in Ramen.
- The spike lives downstream in `../ramen/authoring`; Authoring imports neither
  Ramen nor UWS.
- Downstream `GOWORK=off` compatibility remains an M20 gate until Authoring has
  a consumable public module version.

## Verification

Completed checks:

```bash
(cd ../ramen && go test ./authoring)
(cd ../ramen && go test ./... && go vet ./...)
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
git diff --check
git -C ../tofu diff --check -- authoring ramen
```
````
