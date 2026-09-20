# Status M05

## Goal

Add generic draft lifecycle persistence and artifact-safety helpers.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add draft lifecycle helpers | `[+]` | Added generic JSON draft envelopes, load/save/delete, default draft paths, and autosave helpers. |
| Persist transcripts | `[+]` | Added transcript save/load helpers and optional transcript metadata in draft envelopes. |
| Add atomic writes | `[+]` | Added sibling-temp-file atomic writes and moved M04 prompt transcript writes onto the shared helper. |
| Record safe artifacts | `[+]` | Added file artifact and manifest helpers using M02 `trust` artifact and digest aliases. |

## Notes

- Generated artifact contents remain product-owned and untrusted until
  downstream validation succeeds.
- Review finding fixed before completion: optional session/transcript metadata
  in draft envelopes uses pointers so empty records do not serialize as `{}`.

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
