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
