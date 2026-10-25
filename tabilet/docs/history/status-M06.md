# Retired milestone M06 - Structured JSON Output

**Milestone.** M06
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M06.md
**Source specification.** tabilet/memory-bank/milestone.md#m06-structured-json-output
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
## M06 Structured JSON Output

### M06 Structured JSON Output

**Goal.** Extract provider-neutral structured-output helpers.

Acceptance:

- The package provides schema normalization, a structured-client interface,
  legacy JSON fallback, JSON block extraction, decode helpers, and fake-client
  tests.
- Tests cover structured success, structured error fallback, legacy JSON
  extraction, invalid JSON, nil client, nil output target, and schema
  normalization.
- Provider-specific OpenAI, Anthropic, and Gemini clients stay downstream.

````

## Status record

````markdown
# Status M06

## Goal

Extract provider-neutral structured JSON output helpers.

## Tasks

| Item | State | Notes |
|---|---|---|
| Define structured client interface | `[+]` | Added provider-neutral `Client` and `StructuredClient` interfaces plus result/options types. |
| Normalize schemas | `[+]` | Added deterministic schema normalization from raw JSON, strings, bytes, raw messages, and marshalable values. |
| Add JSON fallback parsing | `[+]` | Added structured-first completion, legacy instruction fallback, JSON object/array block extraction, and decode helpers. |
| Add fake-client tests | `[+]` | Tests cover structured success, structured error fallback, legacy extraction, invalid JSON, nil client, nil output target, and schema normalization. |

## Notes

- OpenAI, Anthropic, Gemini, and any other provider clients stay downstream.
- Review pass confirmed no provider-specific client imports and no OpenUdon,
  Ramen, UWS, or apitools imports.

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
