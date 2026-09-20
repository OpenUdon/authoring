# Status M22

## Goal

Add Authoring documentation and fake-runtime examples.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add README and package docs | `[+]` | Added `README.md`, expanded root package docs, and linked compatibility policy. |
| Add manual prompting example | `[+]` | Added runnable `Example_manualPrompting` using in-memory input/output only. |
| Add agent and structured-output examples | `[+]` | Added runnable structured JSON fallback example with a fake client and no provider credentials. |
| Add progressive loop and report examples | `[+]` | Added runnable fake-runtime progressive loop and report/digest examples. |

## Notes

- Examples must not require model credentials, API execution, workflow
  execution, Terraform/OpenTofu execution, or trusted runner access.
- Examples are compiled and executed by `go test ./...`.

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
