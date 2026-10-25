# Retired milestone M05 - Draft Lifecycle And Artifact Safety

**Milestone.** M05
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M05.md
**Source specification.** tabilet/memory-bank/milestone.md#m05-draft-lifecycle-and-artifact-safety
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
## M05 Draft Lifecycle And Artifact Safety

### M05 Draft Lifecycle And Artifact Safety

**Goal.** Add generic draft lifecycle persistence and safe generated-artifact
records.

Acceptance:

- Authoring supports draft load/save/delete, autosave lifecycle, transcript
  persistence, atomic writes, safe artifact path records, and digest summaries.
- Artifact and digest summaries use `evidence/artifact` and `evidence/digest`.
- Tests cover autosave, atomic write failure handling, path validation, and
  deterministic transcript generation.

````

## Status record

````markdown
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
````
