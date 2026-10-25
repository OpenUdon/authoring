# Retired milestone M28 - API Lifecycle-Ranking Boundary Relocation

**Milestone.** M28
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M28.md
**Source specification.** tabilet/memory-bank/milestone.md#m28-api-lifecycle-ranking-boundary-relocation
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
## M28 API Lifecycle-Ranking Boundary Relocation

### M28 API Lifecycle-Ranking Boundary Relocation

**Goal.** Move API lifecycle sibling ranking to its owning apitools module and
remove the source-breaking pre-1.0 package from Authoring.

This milestone uses M28 because M27 is permanently occupied by the completed
security-alternative prompt-context migration; milestone IDs are never reused.

Acceptance:

- `github.com/OpenUdon/apitools/operationlifecycle` consumes
  `apitools.OperationSummary`, names its scoring constants, and only strips a
  Google Discovery `/upload` prefix when explicit source provenance is present.
- OpenUdon and Ramen translate downstream product context to operation
  summaries and pass their full/focused suites.
- Authoring removes its old package, README/compatibility/API-surface entries,
  and retains no import of apitools.

````

## Status record

````markdown
# Status M28 - API Lifecycle-Ranking Boundary Relocation

State: Complete

M28 is the next unused cross-cutting ID. M27 remains permanently assigned to
the completed prompt-context v2 security-alternative migration.

## Tasks

| Item | State | Notes |
|---|---|---|
| M28.1 apitools ranking package | `[+]` | Added `github.com/OpenUdon/apitools/operationlifecycle` over `apitools.OperationSummary`, documented scoring constants, and scoped `/upload` normalization to explicit Google Discovery provenance. |
| M28.2 OpenUdon migration | `[+]` | The elicitor now passes source-bearing operation summaries to apitools while preserving workflow wording and downstream detail selection. |
| M28.3 Ramen migration | `[+]` | Ramen translates prompt-safe candidates to summaries, maps ranked roles back to its product context, and preserves its Kubernetes update preference. |
| M28.4 Authoring removal | `[+]` | Removed Authoring's package, API-surface types, README package entry, and compatibility surface; documented the intentional pre-1.0 import break. |

## Verification

- Apitools full test/vet checks pass.
- OpenUdon and Ramen focused and full suites pass.
- Authoring standalone/race, apitools test/vet, both workspace consumer suites,
  OpenUdon's scorecard, compatibility, and diff gates pass.
- OpenUdon and Ramen pin Apitools
  `v0.0.0-20260820042238-d51b61ead067` and Authoring
  `v0.0.0-20260820042256-2f73e3526583`; both pass standalone test/vet.

Commit tracking: Apitools `d51b61e`; Authoring interview/persistence/removal
`243121d`, `eee0c0d`, and `2f73e35`; OpenUdon `43294bf`/`065bb84`; Ramen
`330fd89`/`721b540`.
````
