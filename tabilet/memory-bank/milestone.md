# Milestones

This file owns Authoring milestone sequencing, acceptance criteria, current
state, and the status-file index.

## Status ID Pattern

Authoring status files use one uppercase domain letter and a zero-padded number
from `01` through `99`:

```text
M01, M02, M03, ... M09, M10, M11, ...
I01, I02, I03, ...
D01, D02, D03, ...
```

Task/status files use the lane ID:

```text
tabilet/memory-bank/status-M01.md
tabilet/memory-bank/status-M02.md
tabilet/memory-bank/status-M22.md
```

Lane meanings:

- `M`: completed legacy history and future cross-cutting public contracts.
- `I`: interview graph, frontier, readiness, question, session, transcript,
  prompt, and interactive orchestration mechanics.
- `D`: draft lifecycle, review/repair, artifact safety, report, retention, and
  scorecard mechanics.

Do not reuse IDs after a status file exists, reclassify completed M history,
or create aggregate `status.md`. Keep cancelled files with `[X]` rows. Lane
letters classify ownership rather than execution order.

Independent I and D milestones may proceed together only with explicit
non-overlapping package ownership, resolved prerequisites, and downstream
impact. Shared public runtime contracts stay in M or name all cross-lane
dependencies. Prefer one active implementation milestone per lane.

## Current State

Authoring has completed M01 harness setup, M02 Evidence-backed foundation,
M03 session/transcript contracts, M04 local prompt/replay primitives, M05 draft
lifecycle/artifact-safety helpers, M06 structured JSON output helpers, M07
generic progressive iCoT loop, M08 bound-runtime adapter API, M09 decision
evidence/confidence policy, M10 readiness/question-planning contracts, M11
draft review/repair orchestration, M12 agent-mode result contracts, M13
report/retention metadata, M14 scorecard/variant primitives, M15 prompt-safe
context interfaces, M16 OpenUdon adoption phase 1, M17 OpenUdon adoption phase
2, M18 Ramen adapter spike, M19 Ramen desired-state authoring, M20 cross-repo
compatibility gates, M21 public API stabilization, M22 documentation and
examples, M23 internal normalization, M24 shared interactive iCoT extraction,
M25 adaptive evidence-grounded frontier authoring, M26 parallel-lane harness
migration, and M27 security-alternative prompt context. I01 interview/session
integrity, D01 persistence/report hardening, and M28 lifecycle-ranking boundary
relocation are complete and published. OpenUdon and Ramen pin the coordinated
revisions and pass downstream standalone checks. M29 added the neutral engine
with the compatible iCoT facade. M30, approved for planning on 2026-10-02 and
pending, retires that facade with Ramen frozen at its pinned Authoring.
The repository exists
as a public Go module scaffold, and
`AGENTS.md`, `tabilet/memory-bank/`, and `tabilet/evolution/` are symlink-facing paths to the
tracked snapshot under `../tofu/authoring`.

The key architectural decision remains the UWS-style bound-runtime model:
shared upstream orchestration and deterministic records live in Authoring;
runtime-dependent product methods live in downstream OpenUdon and Ramen
adapters.

The next roadmap expands the original M02-M06 sketch instead of discarding it.
The richer implementation source is OpenUdon's mature `internal/authoring`,
`internal/icot`, and `internal/workflowintent` stack, but Authoring extracts
only generic authoring mechanics. OpenUdon prompts, workflow intent schema,
API-source catalog behavior, review package layout, trusted-runner policy, and
OpenUdon wire versions stay in OpenUdon. Ramen desired-state/project
semantics, resource mappings, validation, graphing, planning, state,
governance, and reconciliation stay in Ramen.

M02 added the `trust` package as an Evidence-backed foundation before durable
session/transcript records. M03 added versioned `session` and `transcript`
contracts with deterministic normalization and product-shaped mapping tests.
M04 added generic prompt/replay primitives on top of those records. M05 added
generic lifecycle persistence, atomic writes, transcript persistence, and safe
artifact helpers. M06 added provider-neutral structured JSON output helpers.
M07 added the generic progressive iCoT loop. M08 added the bound-runtime
adapter API. M09 added decision evidence and confidence policy. M10 added
readiness results and question planning. M11 added bounded draft review and
repair orchestration. M12 added noninteractive agent result contracts. M13
added shared report/retention metadata. M14 added scorecard/variant harness
primitives. M15 added prompt-safe context interfaces. M16-M17 migrate
OpenUdon's generic mechanics onto Authoring while preserving OpenUdon behavior.
M18 adds a downstream Ramen adapter spike in `../ramen/authoring` that uses the
Authoring runtime, prompt context, lifecycle, readiness, and report contracts
to draft native `project.uws.yaml` skeletons while Ramen keeps desired-state,
validation, graph, plan, and project semantics downstream. M19 expands that
downstream adapter with Ramen-owned variables, resources, operation roles,
identity fields, dependencies, redaction hints, `x-ramen-desired-state`
emission, and validation/graph/plan gates.
M20-M22 close cross-repo compatibility, public API stabilization, and examples.
M20 adds root import-boundary tests and `scripts/check-compat.sh` for
Authoring plus sibling OpenUdon/Ramen workspace checks. M21 records the
pre-1.0 compatibility policy and tests durable version constants and JSON tags.
M22 adds README/package docs and runnable fake-client/fake-runtime examples.
M23 consolidates duplicated normalization and durable-record helper logic into
internal packages, aligns severity ordering semantics, hardens structured JSON
fence extraction and lifecycle comments, removes milestone IDs from exported
comments, and requires OpenUdon/Ramen compatibility verification plus module
bumps after the Authoring commit.

## Delivery Strategy

Build Authoring in provider-free, model-free slices:

1. `M01`: establish the public module harness, boundaries, execution model,
   status tracking, and initial docs.
2. `M02`: add Evidence as the first dependency and use Evidence primitives for
   neutral diagnostics, redaction, artifact records, and digests.
3. `M03`: define public session and transcript contracts.
4. `M04`: add local prompt and replay primitives.
5. `M05`: add draft lifecycle and artifact safety helpers.
6. `M06`: add provider-neutral structured JSON output helpers.
7. `M07`: add the generic progressive iCoT loop.
8. `M08`: stabilize runtime adapter interfaces.
9. `M09`: add decision evidence and confidence policy.
10. `M10`: add readiness and question planning contracts.
11. `M11`: add draft review and repair orchestration.
12. `M12`: add noninteractive agent-mode result contracts.
13. `M13`: add report and retention metadata.
14. `M14`: add scorecard and variant harness primitives.
15. `M15`: add prompt-safe context interfaces.
16. `M16`: migrate OpenUdon generic prompt, transcript, structured-output,
    lifecycle, and progressive-loop helpers onto Authoring.
17. `M17`: migrate OpenUdon iCoT report, scorecard, and agent shared
    mechanics while keeping OpenUdon-specific schemas and prompts downstream.
18. `M18`: add a Ramen adapter spike for native Ramen/UWS project skeletons.
19. `M19`: expand Ramen desired-state authoring with validation, graph, and
    plan checks as gates.
20. `M20`: add cross-repo compatibility and import-boundary gates.
21. `M21`: stabilize pre-1.0 public API names, versions, JSON tags, and
    compatibility policy.
22. `M22`: add documentation and fake-runtime examples.
23. `M23`: consolidate internal normalization helpers and run downstream
    compatibility adoption.
24. `M24`: extract shared interactive iCoT loop and CLI flag plumbing.
25. `M25`: add the interview graph and consolidate iCoT on frontier rounds.
26. `M26`: migrate the private harness to permanent parallel lanes.
27. `M27`: version security alternatives as OR-of-AND prompt context and
    coordinate downstream adoption.
28. `I01`: harden interview settlement, loop integrity, prompt/transcript
    chronology, and persistence validation shared across entry points.
29. `D01`: consume I01 persistence validation while hardening atomic writes,
    structured fallback extraction, and report cancellation classification.
30. `M28`: move API operation lifecycle ranking to apitools and migrate both
    downstream adapters before removing Authoring's package.

31. `M29`: additive neutral engine extraction for OpenUdon Stage 5, retaining
    the old iCoT API and Ramen compatibility.
32. `M30`: retire the `icot`/`icotcli` compatibility facade with Ramen frozen
    at its pinned Authoring revision.

## Status Files

| Milestone | Status File | Summary |
|---|---|---|
| M01 | [status-M01.md](status-M01.md) | Harness and boundary setup. |
| M02 | [status-M02.md](status-M02.md) | Evidence-backed foundation. |
| M03 | [status-M03.md](status-M03.md) | Session and transcript contracts. |
| M04 | [status-M04.md](status-M04.md) | Local prompt and replay primitives. |
| M05 | [status-M05.md](status-M05.md) | Draft lifecycle and artifact safety. |
| M06 | [status-M06.md](status-M06.md) | Structured JSON output helpers. |
| M07 | [status-M07.md](status-M07.md) | Generic progressive iCoT loop. |
| M08 | [status-M08.md](status-M08.md) | Runtime adapter interface. |
| M09 | [status-M09.md](status-M09.md) | Decision evidence and confidence policy. |
| M10 | [status-M10.md](status-M10.md) | Readiness and question planning. |
| M11 | [status-M11.md](status-M11.md) | Draft review and repair orchestration. |
| M12 | [status-M12.md](status-M12.md) | Agent-mode result contracts. |
| M13 | [status-M13.md](status-M13.md) | Report and retention metadata. |
| M14 | [status-M14.md](status-M14.md) | Scorecard and variant harness primitives. |
| M15 | [status-M15.md](status-M15.md) | Prompt-safe context interfaces. |
| M16 | [status-M16.md](status-M16.md) | OpenUdon adoption phase 1. |
| M17 | [status-M17.md](status-M17.md) | OpenUdon adoption phase 2. |
| M18 | [status-M18.md](status-M18.md) | Ramen adapter spike. |
| M19 | [status-M19.md](status-M19.md) | Ramen desired-state authoring. |
| M20 | [status-M20.md](status-M20.md) | Cross-repo compatibility gate. |
| M21 | [status-M21.md](status-M21.md) | Public API stabilization. |
| M22 | [status-M22.md](status-M22.md) | Documentation and examples. |
| M23 | [status-M23.md](status-M23.md) | Internal normalization and maintenance hardening. |
| M24 | [status-M24.md](status-M24.md) | Shared interactive iCoT extraction. |
| M25 | [status-M25.md](status-M25.md) | Adaptive evidence-grounded frontier authoring. |
| M26 | [status-M26.md](status-M26.md) | Parallel-lane harness migration. |
| M27 | [status-M27.md](status-M27.md) | Security-alternative prompt-context v2. |
| I01 | [status-I01.md](status-I01.md) | Interview, transcript, prompt, and iCoT integrity. |
| D01 | [status-D01.md](status-D01.md) | Persistence and report hardening. |
| M28 | [status-M28.md](status-M28.md) | API lifecycle-ranking boundary relocation. |
| M29 | [status-M29.md](status-M29.md) | Complete: neutral engine, retained iCoT facade and verified publication. |
| M30 | [status-M30.md](status-M30.md) | Pending: retire the iCoT facade; Ramen frozen at its pinned Authoring. |

## Candidate Directions

Candidates have no lane, ID, status file, or execution-order entry until a
fresh scope and dependency review promotes them.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| Another downstream product adapter | Generic interfaces are stable and no third product has demonstrated missing orchestration mechanics. | A named consumer proves a product-neutral gap with fake-runtime tests and a boundary review. |
| New durable authoring wire versions | Interview, session, transcript, and report v1 records remain compatible; prompt context independently moved to v2 for structured security alternatives. | Another shared semantic change cannot be represented additively and downstream consumers approve migration evidence. |
| Provider-specific model clients | Provider SDK and prompting policy remain downstream by design. | Multiple consumers require the same provider-neutral transport contract and approve secret/redaction boundaries. |

## Milestones

### M01 Harness And Boundary Setup

**Goal.** Establish Authoring as a public Go module with Ramen-style memory
bank harness docs, explicit ownership boundaries, the bound-runtime execution
model, and baseline commands.

Acceptance:

- `AGENTS.md`, `tabilet/memory-bank/product.md`, `tabilet/memory-bank/architecture.md`,
  `tabilet/memory-bank/tech-stack.md`, `tabilet/memory-bank/milestone.md`, and
  `tabilet/memory-bank/status-M01.md` exist through the tracked snapshot.
- The root checkout has the same symlink-facing harness pattern used by Ramen.
- `go.mod` exists with module path `github.com/OpenUdon/authoring`.
- Documentation states that Authoring owns shared orchestration while OpenUdon
  and Ramen own runtime adapters.

### M02 Evidence-Backed Foundation

**Goal.** Add `github.com/OpenUdon/evidence` as the first dependency and use
Evidence primitives for neutral durable records.

Acceptance:

- Authoring imports `evidence/diagnostic`, `evidence/redact`,
  `evidence/artifact`, and `evidence/digest` where shared records need those
  concepts.
- Durable trust evidence is not redefined in Authoring.
- Public docs state which Evidence packages are used and which product-specific
  approval or governance records remain downstream.

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

### M05 Draft Lifecycle And Artifact Safety

**Goal.** Add generic draft lifecycle persistence and safe generated-artifact
records.

Acceptance:

- Authoring supports draft load/save/delete, autosave lifecycle, transcript
  persistence, atomic writes, safe artifact path records, and digest summaries.
- Artifact and digest summaries use `evidence/artifact` and `evidence/digest`.
- Tests cover autosave, atomic write failure handling, path validation, and
  deterministic transcript generation.

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

### M07 Generic Progressive iCoT Loop

**Goal.** Move OpenUdon's generic progressive loop shape into Authoring.

Acceptance:

- The loop supports runtime hooks, draft attempts, readiness decisions,
  question planning, defaulted answers, transcript events, final confirmation
  callbacks, and bounded attempts.
- Tests cover success, `needs_input`, cancellation, draft errors, forced
  questions, defaulted answers, and deterministic event order.
- Product-specific draft schemas and issue codes stay downstream.

### M08 Runtime Adapter Interface

**Goal.** Stabilize the bound-runtime API used by downstream OpenUdon and Ramen
adapters.

Acceptance:

- Downstream runtimes can implement draft, review, repair, readiness, question
  planning, artifact writing, and optional document refresh hooks.
- Authoring drives sequencing only and does not embed OpenUdon, Ramen, UWS, API
  source, credential, or execution semantics.
- Fake runtime tests cover missing hooks, hook errors, cancellation, and
  attempt limits.

### M09 Decision Evidence And Confidence Policy

**Goal.** Add generic decision evidence and confidence handling.

Acceptance:

- Decision evidence records stage, slot, value, source, confidence, rationale,
  alternatives, and confirmation requirement.
- Confidence behavior distinguishes auto-accept, review, low-confidence, and
  conflict states.
- Tests cover sorting, conflict handling, redaction posture, and transcript
  event binding.

### M10 Readiness And Question Planning

**Goal.** Add generic readiness/result primitives and question planning
contracts.

Acceptance:

- Readiness and question planning support blocking and warning issues,
  suggested answers, forced-question mechanics, and stable sorting.
- Product-specific issue codes and remediation language remain downstream.
- Tests cover top issue selection, deterministic ordering, forced questions,
  and safe defaults.

Delivered:

- `readiness` defines deterministic issue normalization, readiness summaries,
  top issue selection, decision-confirmation issue projection, and question
  plan/default helpers.
- `icot.Question` aliases `readiness.Question`, and the default iCoT readiness
  policy treats blocking findings as blockers while allowing warning-only
  readiness to proceed unless a downstream `Ready` hook overrides it.

### M11 Draft Review And Repair Orchestration

**Goal.** Add product-neutral review issue/remediation containers and bounded
repair orchestration.

Acceptance:

- Authoring owns attempt limits, event recording, repair summaries, and final
  failure classification.
- Downstream runtimes own repair rules and draft mutation.
- Tests cover repair success, repair exhaustion, review errors, no-op repairs,
  and transcript events.

Delivered:

- `icot` defines review/remediation containers, repair statuses, bounded
  `RunRepair`, and `RunRuntimeRepair` over the M08 review/repair runtime hook
  shapes.
- Repair orchestration records M03 transcript events for review, repair
  attempt/success/error, no-op, and exhaustion while using M10 readiness
  summaries for issue counts and top issue fields.

### M12 Agent Mode Result Contracts

**Goal.** Add generic noninteractive result shapes.

Acceptance:

- Result contracts cover `complete`, `needs_input`, `failed`, and `canceled`.
- Results include top readiness issue, diagnostics, generated artifact
  descriptors, transcript metadata, and digest metadata.
- Tests cover deterministic JSON, missing input summaries, cancellation, and
  failure diagnostics.

Delivered:

- `report` defines `authoring.agent-result.v1` with generic statuses,
  readiness/top issue metadata, M11 repair status fields, M09 decision behavior
  summaries, diagnostics, artifact descriptors, digest records, and
  session/transcript metadata.
- Canonical JSON and normalization keep result output deterministic for
  downstream agent-mode tests.

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

### M14 Scorecard And Variant Harness Primitives

**Goal.** Add generic fixture, variant, and scorecard result primitives.

Acceptance:

- Public structs cover fixture/variant results, expected outcome comparison,
  grouped summary counters, failure-family fields, and report validation
  helpers.
- OpenUdon keeps its corpus and `openudon.icot-*` wire versions.
- Tests cover grouping, comparison, failure-family classification, and stable
  summaries.

Delivered:

- `report` now includes `authoring.scorecard.v1` scorecards, fixture/variant
  result records, outcome normalization, expected/observed comparison, grouped
  counters, failure-family summaries, and canonical scorecard JSON.
- Generic validation helpers return product-neutral diagnostics for malformed
  scorecard shapes without importing OpenUdon corpora or wire schemas.

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

### M16 OpenUdon Adoption Phase 1

**Goal.** Replace OpenUdon's generic prompt, transcript, structured-output,
lifecycle, and progressive-loop helpers with Authoring packages.

Acceptance:

- OpenUdon CLI output, JSON versions, fixtures, and package behavior are
  preserved.
- OpenUdon prompts, workflow intent, catalog planning, repair rules, and
  review packages remain in OpenUdon.
- `go test ./...`, focused iCoT tests, and `make icot-authoring-scorecard`
  pass in `../openudon`.

### M17 OpenUdon Adoption Phase 2

**Goal.** Move OpenUdon iCoT agent/report/scorecard shared mechanics onto
Authoring.

Acceptance:

- OpenUdon uses Authoring for shared agent result, report, retention,
  scorecard, and variant mechanics.
- OpenUdon-specific session schema, prompts, workflow intent, catalog planning,
  repair rules, and package artifacts remain in OpenUdon.
- OpenUdon fixture outcomes and wire versions remain unchanged unless a
  deliberate migration note records the change.

### M18 Ramen Adapter Spike

**Goal.** Add a Ramen downstream adapter spike that uses Authoring to draft
native Ramen/UWS project skeletons from an operator goal.

Acceptance:

- The spike can produce a draft `project.uws.yaml` skeleton.
- The generated project validates when enough metadata is present.
- The adapter stops before graph/plan if required mapping metadata is missing.
- The adapter is downstream in Ramen; Authoring still imports no Ramen, UWS, or
  OpenUdon packages.

### M19 Ramen Desired-State Authoring

**Goal.** Expand the Ramen adapter for real desired-state authoring metadata.

Acceptance:

- The adapter covers variables, resources, operation roles, identity fields,
  dependencies, redaction hints, and `x-ramen-desired-state`.
- Validation, graph, and plan checks are acceptance gates.
- The implementation remains downstream in Ramen; Authoring still imports no
  Ramen, UWS, OpenUdon, or API-source packages.
- Ramen-specific desired-state and reconciliation semantics remain in Ramen.

### M20 Cross-Repo Compatibility Gate

**Goal.** Prove Authoring's import boundaries and downstream compatibility.

Acceptance:

- Import-boundary tests prove Authoring imports neither OpenUdon nor Ramen.
- OpenUdon and Ramen can both import Authoring.
- Default tests are provider-free, model-free, executor-free, and
  credential-free.
- The reusable compatibility gate is documented for future exported API
  changes.

### M21 Public API Stabilization

**Goal.** Review exported package names, version constants, JSON tags, and
compatibility policy before broader downstream reliance.

Acceptance:

- Public APIs are marked as pre-1.0 unstable where appropriate.
- Version constants and JSON tags are reviewed for durable contracts.
- Migration expectations for OpenUdon and Ramen adopters are documented.
- Machine checks cover durable `authoring.*.v1` constants and exported JSON
  tags on durable records.

### M22 Documentation And Examples

**Goal.** Add Authoring README/package docs and minimal fake-runtime examples.

Acceptance:

- Docs cover manual prompting, agent mode, structured JSON fallback,
  progressive loop, and report/digest generation.
- Examples use fake runtimes and fake clients only.
- Documentation restates boundaries: no OpenUdon/Ramen semantics, no live model
  or executor requirement, and no credential storage.
- Examples compile and run in default `go test ./...`.

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

### M25 Adaptive Evidence-Grounded Frontier Authoring

**Goal.** Replace fixed, one-question-at-a-time iCoT loops with one
dependency-aware frontier-round engine and a generic interview state contract.

Acceptance:

- `interview` exposes versioned state, node, evidence, answer, and deferral
  records with open, settled, deferred, and inapplicable node statuses.
- Validation diagnoses duplicate IDs, missing dependencies and evidence,
  cycles, invalid statuses/transitions, and incomplete or invalid deferrals.
- The deterministic frontier includes every open node whose prerequisites are
  settled, ordered by priority and ID; a complete answer set can be applied
  atomically.
- `readiness.Question` has one durable forced flag and recommendation plus
  priority, concise public rationale, and evidence references. Compatibility
  aliases are Go-only and absent from the JSON wire.
- Generic, bound-runtime, and interactive entry points use one active engine
  that displays the entire numbered round before collection and performs one
  normalization/autosave after `ApplyRound`.
- The loop has no frontier breadth ceiling and stops on completion,
  cancellation, downstream-approved finalization, three consecutive
  no-progress rounds, or I01's configurable emergency round fuse.
- Tests cover more than 20 decisions, independent and dependent frontiers,
  deterministic ordering, atomic autosave, full/normal/fast visibility,
  graph failures, invalid deferrals/transitions, cancellation, and repeated
  no-progress.
- Authoring full/standalone/vet/diff/import-boundary gates and unchanged Ramen
  compatibility pass. OpenUdon consumes the frontier API in its M70 v2 wire
  migration.

### M26 Parallel-Lane Harness Migration

**Goal.** Adopt permanent domain lanes and runner-compatible task ledgers
without changing Authoring APIs or orchestration behavior.

Acceptance:

- M01-M25 remain permanent under their historical IDs.
- Future interview and durable-lifecycle/report work has explicit ownership.
- Candidate work stays unnumbered until promotion.
- Structural, module, compatibility, and unattended no-action runner checks
  pass.

### M27 Security-Alternative Prompt-Context V2

**Goal.** Preserve API authentication alternatives without flattening them in
the shared prompt-safe operation contract.

Acceptance:

- `authoring.prompt-context.v2` represents outer OR alternatives, inner AND
  symbolic bindings, and explicit anonymous alternatives.
- v1 prompt-context inputs are rejected rather than silently losing the old
  flattened credential list.
- OpenUdon and Ramen select or explicitly defer one alternative before
  emitting runnable artifacts, and cross-repository compatibility passes.

### I01 Interview, Transcript, Prompt, And iCoT Integrity

**Goal.** Make every loop entry point share one complete-round interview
transaction, semantic progress policy, chronological event stream, and
persistence-validation contract.

Ownership and dependency:

- Owns `interview`, `icot`, `prompt`, `session`, `transcript`, and readiness
  compatibility behavior for this slice.
- Exposes `session.ValidateForPersistence`; D01 consumes that contract and
  does not redefine sensitive-value policy.
- OpenUdon and Ramen own their binding callbacks, prompt wording, product
  mutations, readiness rules, and artifact formats.

Acceptance:

- `interview.ApplyRound` resolves every current frontier node exactly once by
  answer or complete deferral, appends evidence atomically, advances one round,
  resets the no-progress counter, and preserves the input on failure.
- `icot.InterviewBinding` is usable through generic, runtime, and interactive
  entry points; OpenUdon and Ramen remove duplicate settlement plumbing.
- Semantic fingerprints return explicit errors, event IDs are monotonic,
  transcript event order survives normalize/save/load, and the emergency fuse
  defaults to 1,000 rounds when `MaxRounds` is zero.
- Interactive planning reuses an unchanged provisional frontier, replans once
  after declared draft mutation, attributes legacy answers to the displayed
  frontier, and emits no phantom decision events.
- Prompt messages/frontier text/opening labels are configurable, only
  user-sourced `cancel` answers cancel, nil contexts/sessions fail
  descriptively, and readiness aliases normalize write-through.

### D01 Persistence And Report Hardening

**Goal.** Prevent unsafe durable state and close local-file, extraction, and
cancellation-classification correctness gaps.

Ownership and dependency:

- Depends on I01's completed `session.ValidateForPersistence` contract.
- Owns `lifecycle`, `structured`, `report`, and the internal cancellation
  sentinel for this slice; it does not change product artifact schemas.
- May proceed beside M28 because package ownership does not overlap.

Acceptance:

- Canonical session JSON, prompt transcript saves, and capable draft saves
  reject explicitly sensitive unredacted turns or answers without echoing
  values in errors; redacted records remain persistable.
- `AtomicWrite` writes, chmods, syncs, closes, renames, and syncs the parent
  directory, documenting the post-rename durability caveat.
- Bounded fenced-JSON extraction handles trailing prose, and cancellation
  lookalike text is not classified as cancellation.
- Existing v1 durable versions and diagnostic ordering remain unchanged.

### M28 API Lifecycle-Ranking Boundary Relocation

**Goal.** Move API lifecycle sibling ranking to its owning apitools module and
remove the source-breaking pre-1.0 package from Authoring.

This milestone uses M28 because M27 is permanently occupied by the completed
security-alternative prompt-context migration; milestone IDs are never reused.

Acceptance:

- `github.com/OpenUdon/apitools/operationlifecycle` consumes
  `apitools.OperationSummary`, names its scoring constants, and only strips a
  Google Discovery `/upload` prefix when explicit source provenance is present.
- OpenUdon and Ramen translate downstream product context to operation
  summaries and pass their full/focused suites.
- Authoring removes its old package, README/compatibility/API-surface entries,
  and retains no import of apitools.

### M29 Neutral engine with compatible iCoT facade

**Goal.** Expose Authoring's existing generic progressive, interactive, runtime,
repair and atomic interview mechanics through `github.com/OpenUdon/authoring/engine`.
Keep one implementation and preserve the existing `authoring/icot` public API
using type aliases and forwarding functions. This is additive extraction;
Authoring iCoT/icotcli retirement remains deferred.

**Dependencies.** Completed I01, D01 and M28; published baseline
`b417eb681746476cca80bdb14cde3e3c96de980c`. OpenUdon M91.1 inventory is
approved and M91.2 is complete at `1a2570232cdbe8cafecda7e580f91b8a2af12746`.

**Tasks.** M29.1 relocate the existing generic implementation and regression
suite into neutral engine ownership; M29.2 preserve all old exported names,
types, constants, sentinel identity and functions through the compatibility
facade; M29.3 prove atomic-interview, prompt/transcript, no-progress/cancellation,
repair and downstream parity; M29.4 bounded deep review, exact source/module
publication and OpenUdon reconciliation.

**Acceptance.** No duplicated controller or interview transaction; neutral
engine has no direct or transitive `authoring/icot` import. Old APIs compile
with identical signatures and remain usable by unchanged OpenUdon and Ramen.
Prompt/error wording, JSON tags and durable versions, safety and state transitions
remain unchanged. Workspace and standalone full tests/vet, race tests, import
boundary, source/fixture equivalence, compatibility script, OpenUdon scorecard
and downstream tests pass. Use fake providers and disposable local fixtures.
Persist a maximum-ten-iteration closing review and publish the exact verified
revision/module before consumer adoption. Keep completed records in the local
ledger under its existing convention; no new retirement convention is introduced.

**Downstream.** Published M29 precedes resumed OpenUdon M91.3; OpenUdon adopts
its exact module revision and qualifies M91 before M92. Kinet W08 consumes that
engine transitively through qualified OpenUdon. Ramen retains its old import and
receives compatibility verification only; no Ramen source/ledger action is included.

**Authority.** User approved the complete prerequisite proposal, extended the
current Kinet GOAL run and authorized Authoring publication on 2026-09-30.
`COMMIT_POLICY: task`; normal scoped fast-forward pushes to Authoring
`origin/main` (`git@github.com-tabilet:OpenUdon/authoring.git`) follow the
execution owner's recorded exact-diff/remote/checks policy. No live model,
account operation, deployment, UI retirement or new wire semantics.

### M30 Retire the iCoT compatibility facade

**Goal.** Remove `authoring/icot` and `authoring/icotcli`, the compatible facade
over `authoring/engine` that M29 introduced, so Authoring keeps one neutral
implementation and no iCoT-named public surface.

**Provenance.** Kinet Stage 6 (STG-06), approved for planning by the user on
2026-10-02 (Kinet decision R58; coordination in Kinet
`docs/kinet-order.md`, "Stage 6 approved plan"). Planning baseline: Authoring
`dc8f3d61970ae628fc0399b0ef42187aa62a3e5b`, clean worktree. This milestone
promotes the deferred "Authoring iCoT/icotcli retirement" decision recorded by
M29 and `COMPATIBILITY.md`; it reopens no completed milestone.

**Dependencies.** Complete M29. Ramen is frozen at its pinned Authoring
`v0.0.0-20260820042256-2f73e3526583` and is neither updated nor coded; the
approved planning exception intentionally ends compatibility with current
workspace Authoring for those imports. Ramen has not dropped the dependency;
it retains the old pin and is checked separately. This exception supersedes
the earlier migration/drop trigger for M30 only. OpenUdon has no
`authoring/icot` or `authoring/icotcli` import in its dependency graph, Kinet
imports only `session` and `lifecycle`, and W8M imports neither (verified by
search at planning time; re-verify in M30.1).

**Tasks.** M30.1 record the consumer inventory and the Ramen-frozen decision,
re-verifying imports across OpenUdon, Kinet, W8M and every workspace module and
checking Ramen at its pin read-only; M30.2 remove the two packages and their
tests, API-surface, README, COMPATIBILITY, architecture and compatibility-script
references, and update AGENTS/product/tech-stack current descriptions when the
removal is implemented; M30.3 qualify Authoring and its unchanged consumers; M30.4 bounded
deep review and downstream handoff.

**Acceptance.** `authoring/icot` and `authoring/icotcli` no longer exist. The
neutral `authoring/engine` and every other public package keep their names,
signatures, JSON tags, durable versions, prompt and transcript bytes, safety and
state behavior; no new import or dependency is added and `engine` still has no
import of a downstream product. The retirement is recorded as an intentional
pre-1.0 source break for consumers of the removed packages, and
`scripts/check-compat.sh` and the import-boundary tests no longer require Ramen
to build against workspace Authoring. Standalone and workspace full tests and
vet, and race tests, pass. Use a temporary workspace or temporary modfile
override binding the exact new Authoring source into unchanged OpenUdon and
Kinet; include Kinet explicitly because its normal Makefile uses GOWORK=off.
Record the effective Authoring directory/revision for each consumer check;
normal checks against older pins do not establish this gate. No consumer
manifest or operator go.work edit is permitted. OpenUdon and Kinet build and
pass their affected checks against the new Authoring without edits. Ramen at its frozen pin is
verified read-only (for example a build with `GOWORK=off`); if that module
cannot be resolved with the available cache, record the evidence gap rather than
altering Ramen. `~/Workspace/go.work` still lists `./ramen`, so Ramen will stop
building in that workspace after M30; updating or excluding it is a user action
this milestone does not perform.

**Verification.** Fake providers and disposable local fixtures; no live model,
account or network service. Record exact source revisions and commands.

**Downstream.** None required: Kinet and OpenUdon do not import the removed
packages. OpenUdon's remaining legacy ICOT-named identifiers stay an
OpenUdon-owned handoff, and Kinet moves its Authoring pin only if it chooses to.

**Reconciliation provenance.** Stage 6 planning handoff reconciliation F04/F06,
source priority and external review baseline not supplied, local P2. Revalidated
at Authoring `dc8f3d61970ae628fc0399b0ef42187aa62a3e5b`, including uncommitted
M30 planning drafts. Evidence: `COMPATIBILITY.md`, `scripts/check-compat.sh`,
Ramen's `go.mod`/iCoT imports and Kinet's standalone Makefile. The unchanged
Ramen standalone build passed with cached dependencies and Go 1.26.7,
GOWORK=off and readonly module resolution. Original clean baseline wording
refers to the pre-draft snapshot. The user approved these amendments on
2026-10-02. The review counter stays 0/10; no implementation is accepted.
See [prompt v8](../evolution/prompt-v8.md) and [current state v8](../evolution/result-v8.md).
The intentional source-break exception replaces M29's earlier migration
condition without reopening M29 or changing its historical evidence.

**Authority.** Planning approval only. Execution needs a separate request that
names its commit policy; planning grants no commit, push, tag or module
publication. No Ramen, udon-ui, W8M or other sibling edit; Ramen is read-only
evidence.
