# Retired milestone M23 - Internal Normalization And Maintenance Hardening

**Milestone.** M23
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M23.md
**Source specification.** tabilet/memory-bank/milestone.md#m23-internal-normalization-and-maintenance-hardening
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
## M23 Internal Normalization And Maintenance Hardening

### M23 Internal Normalization And Maintenance Hardening

**Goal.** Convert the deep Authoring code-quality review into an internal-only
maintenance slice that removes normalization drift without changing public
contracts or moving downstream semantics upstream.

Acceptance:

- `internal/norm` provides the single stdlib-only implementation for token
  normalization, string comparison, first-non-empty selection, trimming,
  metadata normalization, and severity ordering.
- `internal/records` or an equivalent internal package normalizes
  Evidence-backed digest and artifact records without adding exported `trust`
  APIs.
- Unknown and empty readiness severities sort consistently as blocking/highest
  priority across `readiness`, `session`, and `transcript`, with tests proving
  the shared policy.
- Report status and outcome normalization still accepts legacy hyphenated
  spellings such as `Needs-Input` after the unified token rule.
- `structured.ExtractJSONBlock` strips fenced JSON iteratively or with a clear
  bound and has nested-fence regression coverage.
- Exported comments describe behavior rather than historical milestone IDs,
  and `DeleteDraft` documents best-effort parent-directory pruning.
- `go test ./...`, `go vet ./...`, `GOWORK=off go test ./...`,
  `GOWORK=off go vet ./...`, `git diff --check`, and
  `scripts/check-compat.sh` pass.
- After the Authoring change is committed, OpenUdon and Ramen consume the new
  Authoring pseudo-version, refresh module sums, and pass default plus
  `GOWORK=off` checks without behavior or wire-version changes.

````

## Status record

````markdown
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
````
