# Retired milestone M04 - Local Prompt And Replay Primitives

**Milestone.** M04
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M04.md
**Source specification.** tabilet/memory-bank/milestone.md#m04-local-prompt-and-replay-primitives
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
## M04 Local Prompt And Replay Primitives

### M04 Local Prompt And Replay Primitives

**Goal.** Extract generic prompt modes and deterministic replay behavior from
OpenUdon.

Acceptance:

- The public prompt package supports ask/show/silent defaults, required
  prompts, yes/no prompts, replay scripts, prompt transcript save/load, and
  prompt-label assertions.
- Tests cover forced questions, required answers, deterministic labels, replay
  exhaustion, and transcript persistence.
- Prompt text and product-specific question wording stay downstream.

````

## Status record

````markdown
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
````
