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
