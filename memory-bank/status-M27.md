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
