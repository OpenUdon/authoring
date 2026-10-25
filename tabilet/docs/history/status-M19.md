# Retired milestone M19 - Ramen Desired-State Authoring

**Milestone.** M19
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M19.md
**Source specification.** tabilet/memory-bank/milestone.md#m19-ramen-desired-state-authoring
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
## M19 Ramen Desired-State Authoring

### M19 Ramen Desired-State Authoring

**Goal.** Expand the Ramen adapter for real desired-state authoring metadata.

Acceptance:

- The adapter covers variables, resources, operation roles, identity fields,
  dependencies, redaction hints, and `x-ramen-desired-state`.
- Validation, graph, and plan checks are acceptance gates.
- The implementation remains downstream in Ramen; Authoring still imports no
  Ramen, UWS, OpenUdon, or API-source packages.
- Ramen-specific desired-state and reconciliation semantics remain in Ramen.

````

## Status record

````markdown
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
````
