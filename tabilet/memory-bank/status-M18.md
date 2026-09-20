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
