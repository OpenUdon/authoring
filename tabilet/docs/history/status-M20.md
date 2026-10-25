# Retired milestone M20 - Cross-Repo Compatibility Gate

**Milestone.** M20
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M20.md
**Source specification.** tabilet/memory-bank/milestone.md#m20-cross-repo-compatibility-gate
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
## M20 Cross-Repo Compatibility Gate

### M20 Cross-Repo Compatibility Gate

**Goal.** Prove Authoring's import boundaries and downstream compatibility.

Acceptance:

- Import-boundary tests prove Authoring imports neither OpenUdon nor Ramen.
- OpenUdon and Ramen can both import Authoring.
- Default tests are provider-free, model-free, executor-free, and
  credential-free.
- The reusable compatibility gate is documented for future exported API
  changes.

````

## Status record

````markdown
# Status M20

## Goal

Add cross-repo compatibility and import-boundary gates.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add Authoring import-boundary tests | `[+]` | Added root `boundary_test.go` to fail if `go list -deps ./...` includes OpenUdon, Ramen, UWS, or apitools. |
| Add dependency checks | `[+]` | Added `scripts/check-compat.sh` to run Authoring checks and sibling OpenUdon/Ramen tests from the parent workspace when those repos exist. |
| Keep default tests clean | `[+]` | Gate includes default and `GOWORK=off` Authoring `go test`/`go vet`, with no model, provider, executor, credential, or live API requirement. |
| Document compatibility gate | `[+]` | Added `COMPATIBILITY.md` and referenced the reusable gate for future exported API changes. |

## Notes

- This milestone protects Authoring as a public shared module rather than a
  product-specific helper package.
- OpenUdon/Ramen downstream checks are workspace checks until Authoring has a
  consumable public module version.

## Verification

Completed checks:

```bash
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../authoring && ./scripts/check-compat.sh)
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
git diff --check
git -C ../tofu diff --check -- authoring
```
````
