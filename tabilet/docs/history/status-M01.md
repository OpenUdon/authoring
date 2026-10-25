# Retired milestone M01 - Harness And Boundary Setup

**Milestone.** M01
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M01.md
**Source specification.** tabilet/memory-bank/milestone.md#m01-harness-and-boundary-setup
**Evidence.** be023f06d0837b58a651449b6d89a00c12eb04cd
**Worktree.** includes uncommitted changes
**Review.** legacy
**Review iterations.** not recorded
**Verification.** Historical verification and its limitations are retained verbatim in the status record; original commands are not re-run or relabelled by this retirement. At the observed evidence revision, go test ./... and go vet ./... passed on 2026-10-25; literal source preservation and retirement-envelope validation passed before removal. M01 originally listed planned checks only; no original executed-check result is inferred.
**Consolidated into.** Current product, architecture, stack, and lessons remain in tabilet/memory-bank/; no implementation or current product contract changes. Superseded roadmap and retention wording is preserved in knowledge.md.
**Specification encoding.** The added level-two identity heading is a retirement-parser wrapper; the following level-three section is the complete original source, retained verbatim.
**Legacy closure.** This milestone closed before adoption of the persisted review gate. Its original closure evidence is the complete literal status record; no review iteration or retrospective implementation review is invented.

## Milestone specification

````markdown
## M01 Harness And Boundary Setup

### M01 Harness And Boundary Setup

**Goal.** Establish Authoring as a public Go module with Ramen-style memory
bank harness docs, explicit ownership boundaries, the bound-runtime execution
model, and baseline commands.

Acceptance:

- `AGENTS.md`, `tabilet/memory-bank/product.md`, `tabilet/memory-bank/architecture.md`,
  `tabilet/memory-bank/tech-stack.md`, `tabilet/memory-bank/milestone.md`, and
  `tabilet/memory-bank/status-M01.md` exist through the tracked snapshot.
- The root checkout has the same symlink-facing harness pattern used by Ramen.
- `go.mod` exists with module path `github.com/OpenUdon/authoring`.
- Documentation states that Authoring owns shared orchestration while OpenUdon
  and Ramen own runtime adapters.

````

## Status record

````markdown
# Status M01

## Goal

Establish the Authoring public module harness and boundary record.

## Tasks

| Item | State | Notes |
|---|---|---|
| Create tracked harness snapshot | `[+]` | Add AGENTS, memory-bank docs, milestone index, and first status file under `../tofu/authoring`. |
| Wire checkout harness paths | `[+]` | Use the same symlink-facing `AGENTS.md`, `memory-bank/`, and `evolution/` pattern as Ramen. |
| Record execution model | `[+]` | AGENTS and architecture docs describe shared upstream orchestration with downstream bound runtimes. |
| Add module scaffold | `[+]` | Add `go.mod` for `github.com/OpenUdon/authoring`. |

## Notes

- Initial design input is preserved in `../authoring/x.md`.
- Authoring starts as a documentation/module harness. Implementation begins in
  M02 with session and transcript contracts.

## Verification

Planned checks:

```bash
go test ./...
go vet ./...
git diff --check
git -C ../tofu diff --check -- authoring
```
````
