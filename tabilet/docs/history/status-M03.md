# Retired milestone M03 - Session And Transcript Contracts

**Milestone.** M03
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M03.md
**Source specification.** tabilet/memory-bank/milestone.md#m03-session-and-transcript-contracts
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
## M03 Session And Transcript Contracts

### M03 Session And Transcript Contracts

**Goal.** Define public `session` and `transcript` packages for prompt turns,
structured events, model/provider provenance, answers, decision evidence,
readiness issues, artifact summaries, and deterministic normalization.

Acceptance:

- Contracts are provider-free and product-neutral.
- Default tests exercise deterministic serialization, normalization, and stable
  sorting.
- OpenUdon and Ramen can map existing authoring state into the contracts
  without importing each other.

````

## Status record

````markdown
# Status M03

## Goal

Define public session and transcript contracts for durable generic authoring
records.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add `session` package | `[+]` | Added versioned session state, prompt turns, answers, readiness issues, decisions, artifact summaries, diagnostics, metadata, normalization, and canonical JSON. |
| Add `transcript` package | `[+]` | Added versioned transcript records, turns, events, model/provider provenance, diagnostics, artifact summaries, metadata, normalization, and canonical JSON. |
| Normalize deterministically | `[+]` | Session and transcript records trim fields, normalize tokens, sort records, normalize M02 `trust` diagnostics/artifacts, and produce stable JSON. |
| Add mapping tests | `[+]` | Tests map OpenUdon-shaped and Ramen-shaped structs into generic records without importing either product. |

## Notes

- Keep provider clients and prompt text downstream.
- Review finding fixed before completion: optional transcript provider
  provenance uses pointers so empty provider records do not serialize as `{}`.

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
