# Retired milestone M15 - Prompt-Safe Context Interfaces

**Milestone.** M15
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M15.md
**Source specification.** tabilet/memory-bank/milestone.md#m15-prompt-safe-context-interfaces
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
## M15 Prompt-Safe Context Interfaces

### M15 Prompt-Safe Context Interfaces

**Goal.** Define narrow product-neutral context interfaces for prompt-safe
metadata.

Acceptance:

- Interfaces cover source documents, operation candidates, schemas, and
  security or credential binding names.
- Authoring does not import `apitools`, UWS, OpenUdon, or Ramen.
- OpenUdon and Ramen adapters translate product metadata into the prompt-safe
  shapes.

Delivered:

- `promptcontext` defines `authoring.prompt-context.v1` source documents,
  operation candidates, schema hints, field hints, symbolic credential
  bindings, normalization, redaction-backed prompt-safety guards, and
  canonical JSON.
- Tests include normalization/redaction coverage plus a build-info
  import-boundary check for apitools, UWS, OpenUdon, and Ramen.

````

## Status record

````markdown
# Status M15

## Goal

Define prompt-safe context interfaces for product-neutral metadata.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define source document shapes | `[+]` | Added prompt-safe source document identity, kind, title, URI, media type, summary, digest, and metadata fields. |
| Define operation candidate shapes | `[+]` | Added operation candidates with source/operation IDs, verb/path, schema references, symbolic credential binding names, tags, confidence, and rationale without API-source imports. |
| Define schema and credential binding shapes | `[+]` | Added schema/field hints and symbolic credential binding records with redaction-backed normalization. |
| Add import-boundary tests | `[+]` | Added build-info import-boundary test and retained external `go list -deps` scan for apitools, UWS, OpenUdon, and Ramen. |

## Notes

- Downstream adapters translate real API metadata into these shapes.
- Deep review findings fixed before commit: metadata values now use
  key-sensitive redaction, URI/path text is redacted, and example normalization
  no longer mutates caller-owned slices.

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
