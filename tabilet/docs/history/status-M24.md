# Retired milestone M24 - Shared Interactive iCoT Extraction

**Milestone.** M24
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M24.md
**Source specification.** tabilet/memory-bank/milestone.md#m24-shared-interactive-icot-extraction
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
## M24 Shared Interactive iCoT Extraction

### M24 Shared Interactive iCoT Extraction

**Goal.** Move reusable interactive iCoT lifecycle and CLI flag plumbing into
Authoring while keeping product prompts, artifact schemas, source parsing,
workflow semantics, and provider clients downstream.

Acceptance:

- `icot` exposes the shared opening, extraction, readiness/question,
  prefill, autosave, transcript, cancellation, and needs-input lifecycle.
- `icotcli` exposes provider-neutral prompt-mode/model flag helpers.
- Durable issue, question, prompt-turn, and event vocabulary converges without
  importing OpenUdon, Ramen, UWS, or apitools.
- OpenUdon and Ramen compatibility checks pass.

````

## Status record

````markdown
# Status M24 - Shared Interactive iCoT Extraction

State of the public interactive iCoT extraction used by OpenUdon and Ramen.

## Goal

Move the reusable interactive iCoT lifecycle from downstream product code into
Authoring while keeping prompts, artifact schemas, provider clients, API source
parsing, UWS semantics, Ramen desired-state behavior, and OpenUdon package
behavior downstream.

## Tasks

| Item | State | Notes |
|---|---|---|
| M24.1 public loop API | `[+]` | Added `icot.DraftRequest`, `Extractor`, `NoopExtractor`, `InteractiveHooks`, `InteractiveLifecycleOptions`, `RunInteractive`, `RunInteractiveWithLifecycle`, and product-payload `Event`. |
| M24.2 loop behavior | `[+]` | Ported opening prompt, optional extractor drafting/disambiguation, readiness and question loop, deterministic prefill, draft-question hooks, autosave, transcript save, cancellation, and needs-input behavior. |
| M24.3 CLI plumbing | `[+]` | Added `icotcli` for shared flags and model/prompt-mode label resolution without provider clients or product imports. |
| M24.4 tests | `[+]` | Added fake-extractor coverage for opening prompt, draft merge, autosave, transcript save, question drafting, cancellation, lifecycle cleanup, and prompt-mode/model flag helpers. |
| M24.5 downstream adoption | `[+]` | OpenUdon now delegates its internal generic progressive loop wrappers to the public Authoring API; Ramen consumes the public CLI flag helpers for `ramen icot`. |
| M24.6 vocabulary convergence follow-up | `[+]` | `icot.ReadinessIssue`, `icot.InteractiveQuestion`, and `icot.PromptTurn` now alias durable session/readiness types; readiness issues carry neutral operation/path/remediation fields, questions carry suggested/force/grouped compatibility fields, and interactive events now record both compatibility `kind` and durable `type` with transcript projection helpers. |
| M24.7 `icotcli` helper cleanup | `[+]` | Removed the private `icotcli.firstNonEmpty` duplicate and computes the prompt-mode default directly in `AddFlags`. |
| M24.8 downstream compatibility verification | `[+]` | Added focused compatibility tests for shared issue/question/turn/event behavior and verified OpenUdon/Ramen consumers against the converged exported API. |

## Acceptance Criteria

- Authoring owns only generic interactive loop mechanics and CLI flag plumbing.
- OpenUdon and Ramen continue to own product-specific prompts, model clients,
  artifact formats, API metadata translation, validation, graphing, planning,
  and execution boundaries.
- Default Authoring tests remain provider-free, model-free, executor-free, and
  free of OpenUdon/Ramen/UWS/apitools imports in generic packages.

## Verification

- `go test ./...`
- Downstream focused checks in OpenUdon and Ramen.
- Full standalone `GOWORK=off` checks remain deferred until sibling module
  pseudo-versions are published/resolved.
````
