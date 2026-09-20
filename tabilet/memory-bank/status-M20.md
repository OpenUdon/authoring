# Status M20

## Goal

Add cross-repo compatibility and import-boundary gates.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add Authoring import-boundary tests | `[+]` | Added root `boundary_test.go` to fail if `go list -deps ./...` includes OpenUdon, Ramen, UWS, or apitools. |
| Add dependency checks | `[+]` | Added `scripts/check-compat.sh` to run Authoring checks and sibling OpenUdon/Ramen tests from the parent workspace when those repos exist. |
| Keep default tests clean | `[+]` | Gate includes default and `GOWORK=off` Authoring `go test`/`go vet`, with no model, provider, executor, credential, or live API requirement. |
| Document compatibility gate | `[+]` | Added `COMPATIBILITY.md` and referenced the reusable gate for future exported API changes. |

## Notes

- This milestone protects Authoring as a public shared module rather than a
  product-specific helper package.
- OpenUdon/Ramen downstream checks are workspace checks until Authoring has a
  consumable public module version.

## Verification

Completed checks:

```bash
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../authoring && ./scripts/check-compat.sh)
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
git diff --check
git -C ../tofu diff --check -- authoring
```
