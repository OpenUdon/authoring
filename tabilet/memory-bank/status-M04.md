# Status M04

## Goal

Extract local prompt and replay primitives that are generic enough for both
OpenUdon and Ramen.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add prompt modes | `[+]` | `prompt.Session` supports ask, show, and silent default handling. |
| Add required prompts | `[+]` | Added required free-form/default prompts and yes/no prompts with validation. |
| Add replay scripts | `[+]` | Added replay scripts, replay readers, stable turn sequence IDs, and prompt transcript save/load using M03 `session` and `transcript` records. |
| Add label assertions | `[+]` | Added deterministic prompt label order assertions for replay tests. |

## Notes

- Product-specific wording, prompt templates, and question policy stay
  downstream.
- Review findings fixed before completion: prompt turns now receive stable
  sequence IDs so M03 normalization preserves replay order, and transcript
  writes enforce `0600` permissions even when overwriting existing files.

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
