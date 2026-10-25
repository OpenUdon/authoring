# Retired milestone M08 - Runtime Adapter Interface

**Milestone.** M08
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M08.md
**Source specification.** tabilet/memory-bank/milestone.md#m08-runtime-adapter-interface
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
## M08 Runtime Adapter Interface

### M08 Runtime Adapter Interface

**Goal.** Stabilize the bound-runtime API used by downstream OpenUdon and Ramen
adapters.

Acceptance:

- Downstream runtimes can implement draft, review, repair, readiness, question
  planning, artifact writing, and optional document refresh hooks.
- Authoring drives sequencing only and does not embed OpenUdon, Ramen, UWS, API
  source, credential, or execution semantics.
- Fake runtime tests cover missing hooks, hook errors, cancellation, and
  attempt limits.

````

## Status record

````markdown
# Status M08

## Goal

Stabilize the generic bound-runtime adapter API.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define runtime hooks | `[+]` | Added `icot.Runtime`, optional normalize/ready/refresh/draft-policy/review/repair interfaces, `RuntimeConfig`, `BindRuntime`, and `RunRuntime`. |
| Separate sequencing from domain behavior | `[+]` | Authoring binds runtime hooks into M07 loop sequencing; downstream runtimes own draft, readiness, question, answer, artifact, review, and repair semantics. |
| Add hook error handling | `[+]` | Tests cover nil runtime rejection, artifact-write errors, refresh errors, draft-error fallback, and existing M07 cancellation/attempt handling. |
| Add interface docs | `[+]` | Runtime interface comments document downstream responsibilities and keep product semantics out of Authoring. |

## Notes

- Runtime hooks must not smuggle OpenUdon, Ramen, UWS, API-source, credential,
  or execution semantics into Authoring.
- Review pass confirmed no OpenUdon, Ramen, UWS, or apitools imports.

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
