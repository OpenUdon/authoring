# Retired milestone M13 - Report And Retention Metadata

**Milestone.** M13
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M13.md
**Source specification.** tabilet/memory-bank/milestone.md#m13-report-and-retention-metadata
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
## M13 Report And Retention Metadata

### M13 Report And Retention Metadata

**Goal.** Add shared report metadata and digest sidecar helpers.

Acceptance:

- Report metadata includes run ID, command, commit, generation time, retention
  class, provider-output flag, archive safety, redaction-required flag, and
  digest sidecars.
- The implementation uses Evidence digest and redaction helpers.
- Tests cover retention classes, digest sidecars, provider-output flags, and
  redaction-required propagation.

Delivered:

- `report` now includes `authoring.report-metadata.v1` report metadata with
  run/command/commit/generation fields, generic retention classes,
  provider-output/archive/redaction flags, and digest sidecars.
- Digest sidecars use M02 `trust` digest records and safe artifact path
  cleaning; redaction-required propagation delegates to M02 `trust` redaction
  helpers.

````

## Status record

````markdown
# Status M13

## Goal

Add shared report and retention metadata.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define report metadata | `[+]` | Added `ReportMetadata` with run ID, command, commit, generation time, retention class, provider-output flag, archive safety, and redaction-required flag. |
| Add digest sidecars | `[+]` | Added digest sidecars for report bytes and M02 artifact records using `trust` digest aliases and safe path helpers. |
| Propagate redaction requirements | `[+]` | Added redaction checks through M02 `trust.RedactDocument`; redaction-required metadata clears archive-safe status. |
| Add report tests | `[+]` | Covered retention normalization, digest sidecars, provider-output/archive flags, unsafe sidecar paths, and redaction-required propagation. |

## Notes

- Product-specific report bodies and review artifacts remain downstream.
- Deep review finding fixed before commit: digest sidecar path normalization now
  uses the shared safe artifact path cleaner and drops unsafe paths.

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
````
