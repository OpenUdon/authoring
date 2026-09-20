# Status M21

## Goal

Stabilize the pre-1.0 public API surface.

## Tasks

| Item | State | Notes |
|---|---|---|
| Review exported package names | `[+]` | `README.md` and `COMPATIBILITY.md` document behavior-grouped packages and confirm product adapters stay downstream. |
| Review version constants and JSON tags | `[+]` | Added `api_surface_test.go` coverage for durable `authoring.*.v1` constants and exported JSON tags on durable record types. |
| Document compatibility policy | `[+]` | Added `COMPATIBILITY.md` with pre-1.0 instability, version/tag migration rules, and downstream adapter expectations. |
| Check downstream migrations | `[+]` | Compatibility policy and gate document OpenUdon/Ramen migration expectations and require dependent workspace checks after exported API changes. |

## Notes

- This is an API review milestone, not a feature expansion milestone.
- No product wire versions were changed.

## Verification

Completed checks:

```bash
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
(cd ../authoring && git diff --check)
git -C ../tofu diff --check -- authoring
```
