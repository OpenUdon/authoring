# Status M23

## Goal

Consolidate Authoring's internal normalization and shared record helper logic
after the deep code-quality review, while preserving public contracts and
downstream OpenUdon/Ramen behavior.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add internal normalization helpers | `[+]` | Added stdlib-only `internal/norm` for token normalization, string comparison, trimming, first-non-empty selection, metadata normalization, and severity ordering. |
| Add internal record normalizers | `[+]` | Added internal `records` helpers for `trust.DigestRecord` and `trust.ArtifactRecord` normalization without adding exported `trust.Normalize*` APIs. |
| Migrate package-local copies | `[+]` | Replaced duplicated helpers in `decision`, `readiness`, `session`, `transcript`, `promptcontext`, `lifecycle`, `report`, and `icot`; retained prompt-context metadata redaction as local policy. |
| Align severity and token behavior | `[+]` | Unknown/empty severities now sort consistently as blocking/highest priority, and hyphenated report status/outcome spellings are accepted through unified token normalization. |
| Add regression coverage | `[+]` | Covered `norm.Token`, `FirstNonEmpty`, metadata normalization, severity ordering, record normalization, cross-package readiness ordering, report `Needs-Input`, and session/transcript version fallback. |
| Complete low-risk cleanup | `[+]` | JSON fence extraction is iterative, historical milestone IDs were removed from exported comments, and draft deletion docs now describe best-effort parent pruning. |
| Verify downstream compatibility | `[+]` | Authoring gates, `scripts/check-compat.sh`, and OpenUdon/Ramen default plus `GOWORK=off` gates pass. Downstream pseudo-version bumps remain deferred until the new Authoring commit is published or otherwise resolvable by Go module tooling. |

## Notes

- This milestone is internal maintenance, not a public API expansion.
- Do not change `authoring.*.v1` version constants, JSON tags, exported package
  names, OpenUdon wire versions, or Ramen project/profile semantics.
- Do not add OpenUdon, Ramen, UWS, apitools, model-provider, executor, or live
  API dependencies to Authoring.
- Separate downstream OpenUdon/Ramen milestones are not required unless the
  eventual module bump needs product behavior changes; expected downstream work
  is limited to dependency updates and compatibility verification.
- The implementation intentionally leaves OpenUdon/Ramen `go.mod` unchanged
  until a new Authoring pseudo-version is available to `GOWORK=off` module
  resolution.

## Verification

Completed checks:

```bash
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../authoring && ./scripts/check-compat.sh)
(cd ../openudon && go test ./... && go vet ./...)
(cd ../openudon && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../ramen && go test ./... && go vet ./...)
(cd ../ramen && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../authoring && git diff --check)
git -C ../tofu diff --check -- authoring
```
