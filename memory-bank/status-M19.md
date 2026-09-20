# Status M19

## Goal

Expand the Ramen adapter for native desired-state authoring.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add variable and resource authoring | `[+]` | Ramen adapter options accept Ramen-owned `project.Variable` and `project.Resource` records, preserving symbolic inputs in generated native projects. |
| Add operation and identity metadata | `[+]` | Default drafts translate M15 `promptcontext` schema hints into schema paths, identity attributes, request bindings, and redaction hints; explicit resources preserve operation roles, dependencies, credential bindings, and redaction metadata. |
| Emit desired-state metadata | `[+]` | Generated UWS documents populate `x-ramen-desired-state` with variables, API sources, resources, operation roles, identity fields, dependencies, and redaction metadata. |
| Gate with Ramen checks | `[+]` | Optional adapter gates run Ramen validation, graph, and plan checks for generated drafts and surface diagnostics through the Authoring M12 result. |

## Notes

- Authoring provides orchestration only; Ramen owns mapping, validation,
  graphing, planning, state, and reconciliation.
- The implementation remains downstream in `../ramen/authoring`; Authoring
  imports no Ramen, UWS, OpenUdon, or API-source packages.
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
