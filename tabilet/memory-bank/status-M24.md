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
