# Retired milestone M27 - Security-Alternative Prompt-Context V2

**Milestone.** M27
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M27.md
**Source specification.** tabilet/memory-bank/milestone.md#m27-security-alternative-prompt-context-v2
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
## M27 Security-Alternative Prompt-Context V2

### M27 Security-Alternative Prompt-Context V2

**Goal.** Preserve API authentication alternatives without flattening them in
the shared prompt-safe operation contract.

Acceptance:

- `authoring.prompt-context.v2` represents outer OR alternatives, inner AND
  symbolic bindings, and explicit anonymous alternatives.
- v1 prompt-context inputs are rejected rather than silently losing the old
  flattened credential list.
- OpenUdon and Ramen select or explicitly defer one alternative before
  emitting runnable artifacts, and cross-repository compatibility passes.

````

## Status record

````markdown
# Status M27 - Security-Alternative Prompt-Context V2

State: Complete

| Item | State | Notes |
|---|---|---|
| M27 structured security alternatives | `[+]` | Commit `5c529b6` replaces flattened operation credential bindings with ordered OR alternatives containing AND bindings, preserves explicit anonymous alternatives, bumps prompt context to `authoring.prompt-context.v2`, adds normalization/API-surface tests, and coordinates OpenUdon/Ramen adoption with v1 rejection and resume-safe selection coverage. |

M27 remains closed and is not reused by the later lifecycle-ranking boundary
correction; that coordinated change is indexed as M28.

## Verification

- `go test ./...`
- downstream OpenUdon, Ramen, and Udon full tests
- standalone, vet, compatibility, and diff checks are recorded by the
  consuming Apitools S01 review.
````
