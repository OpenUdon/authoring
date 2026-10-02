# Architecture

Authoring is the shared orchestration layer for progressive authoring flows. It
owns generic loop mechanics, durable authoring records, and runtime adapter
contracts while downstream products own domain behavior.

## System Boundary

- `github.com/OpenUdon/authoring` owns generic sessions, transcripts,
  structured-output helpers, prompt/replay primitives, draft lifecycle
  persistence, dependency-aware interview/frontier orchestration,
  readiness-loop orchestration, repair-loop orchestration,
  decision evidence, report/scorecard metadata, and runtime adapter
  interfaces.
- `github.com/OpenUdon/evidence` owns generic digest, artifact, diagnostic,
  redaction, and approval primitives. Authoring depends on Evidence for shared
  durable trust records instead of redefining them.
- `github.com/OpenUdon/uws` owns public UWS document semantics, model, schema
  lookup, validation, and UWS artifact discovery.
- `github.com/OpenUdon/apitools` owns API source metadata, prompt-safe
  operation/auth/security summaries, and operation lifecycle ranking that
  depends on API-source provenance.
- `github.com/OpenUdon/openudon` owns OpenUdon authoring adapters, prompts,
  workflow intent, synthesis, quality, review, approval, and handoff artifacts.
- `github.com/OpenUdon/ramen` owns Ramen authoring adapters, native
  Ramen/UWS desired-state project generation, validation, graphing, planning,
  state, governance, and reconciliation.

Authoring must not import OpenUdon or Ramen. OpenUdon and Ramen import
Authoring and bind product-specific runtimes.

Authoring also should not import `uws` or `apitools` unless a future boundary
decision proves a generic interface is insufficient. Prompt-safe context
interfaces can describe source documents, operation candidates, schemas, and
credential binding names without naming UWS fields, API source families,
provider resources, or product wire formats.

## Execution Model

Authoring follows a bound-runtime execution model:

```text
shared session/transcript/readiness state
  -> generic authoring engine
  -> downstream runtime Draft/Review/Repair/Write hooks
  -> generic result/transcript/artifact summary
```

The upstream engine owns deterministic orchestration and stable state
transitions. Runtime-specific behavior is supplied through interfaces:

```go
type Runtime interface {
    Draft(ctx context.Context, session Session) (Draft, error)
    Review(ctx context.Context, draft Draft) ([]Diagnostic, error)
    Repair(ctx context.Context, draft Draft, issues []Diagnostic) (Draft, error)
    WriteArtifacts(ctx context.Context, draft Draft) ([]Artifact, error)
}
```

The concrete interface names may change during M02, but the ownership rule
should remain stable: Authoring drives the loop; products implement the domain.
The stabilized runtime API covers draft, review, repair, readiness, interview
binding or question planning, artifact writing, and optional document refresh
hooks. Authoring records monotonic events, presents complete dependency-ready
frontier rounds, applies resolutions to cloned product state, stops after three
consecutive semantic no-progress rounds or the configurable emergency round
fuse, normalizes deterministic state, streams generic and interactive callback
events into one chronological transcript, and returns result contracts;
downstream runtimes own prompts, schemas, issue codes, repair rules, artifact
mutation, and provider/model clients.

## Harness Layout

Private planning uses permanent `status-<LANE><NN>.md` ledgers. `M` retains
legacy/cross-cutting public contracts, `I` owns interview/frontier/session
orchestration, and `D` owns lifecycle/review/report mechanics. The unattended
runner reads the second column of `Item | State | Notes` tables. Candidates
remain unnumbered until promotion.

## Planned Package Layout

```text
trust/       Evidence-backed diagnostic, redaction, artifact, and digest names
session/      durable session state, answers, decisions, readiness issues
transcript/   transcript turns, model/provider metadata, event records
prompt/       local prompt modes, required prompts, replay scripts
lifecycle/    draft load/save/delete, autosave, atomic writes, artifact records
structured/   JSON completion contracts, fallback parsing, schema envelopes
engine/       generic loops, interview binding, and runtime adapter interfaces
interview/    versioned graph, unified evidence, and atomic round resolutions
readiness/    shared readiness issue/result and question-planning primitives
report/       result contracts, report metadata, retention, scorecard helpers
```

Package names may be adjusted during implementation. Keep product-specific
adapters out of this module. Public packages should group behavior rather than
mirror OpenUdon's internal package layout.

M03 adds `session` and `transcript` with versioned JSON contracts,
deterministic normalization, canonical JSON helpers, model/provider provenance
records, decision/readiness/artifact/diagnostic slots, and tests that map
OpenUdon-shaped and Ramen-shaped state without importing either product.
M04 adds `prompt` for local ask/show/silent default modes, required and yes/no
prompts, replay scripts, deterministic label assertions, and prompt
transcript save/load using M03 session/transcript records.
M05 adds `lifecycle` for JSON draft envelopes, autosave helpers, transcript
persistence, atomic writes, and safe artifact record/manifest helpers. The M04
prompt transcript writer now uses the shared lifecycle atomic write helper.
M06 adds `structured` for provider-neutral structured completion interfaces,
schema normalization, structured-first completion, legacy JSON fallback,
JSON block extraction, and decode helpers.
M07 adds `icot` for the generic progressive loop: bounded attempts, draft
hooks, readiness checks, question planning/application, defaulted answers,
autosave hooks, final confirmation, and ordered transcript events.
M08 stabilizes the bound-runtime adapter API on top of `icot`: downstream
runtimes implement draft, readiness, question planning/application, artifact
writing, and optional normalization, readiness policy, document refresh, draft
policy, review, and repair hooks.
M09 adds `decision` for generic decision evidence records, confidence behavior,
conflict merging, redaction posture, and transcript decision-event helpers.
M10 adds `readiness` for deterministic readiness results, blocking/warning
summaries, M09 decision-confirmation issue projection, generic question plans,
forced-question mechanics, and safe default-answer handling. `icot.Question`
is now an alias of `readiness.Question` so existing loop/runtime adapters keep
the M07/M08 surface while the public contract lives in the readiness package.
M11 extends `icot` with bounded draft review/repair orchestration on top of
the M08 runtime hook shapes and M10 readiness results. It records generic
review, repair, exhaustion, no-op, and failure events while downstream
products keep repair rules and draft mutation logic.
M12 adds `report` for noninteractive agent-result contracts: stable
`complete`, `needs_input`, `failed`, and `canceled` statuses, M10 readiness
summaries/top issues, M11 repair status fields, M09 decision behavior
summaries, diagnostics, generated artifact descriptors, digest records, and
M03 session/transcript metadata.
M13 extends `report` with shared report metadata: run ID, command, commit,
generation time, retention class, provider-output flag, archive-safety flag,
redaction-required propagation through M02 `trust` helpers, and digest
sidecars for reports and generated artifacts.
M14 extends `report` with provider-free scorecard and variant-harness
primitives: fixture and variant outcomes, expected/observed comparison,
grouped counters, failure-family summaries, deterministic scorecard JSON, and
generic shape validation diagnostics.
M15 adds `promptcontext` for prompt-safe source document summaries, operation
candidate summaries, schema hints, and symbolic credential binding names. M27
versions that record as `authoring.prompt-context.v2`: operation security is
an ordered outer OR of inner AND binding sets, including explicit empty
anonymous alternatives; no flattened v1 decoder is retained.
Adapters translate UWS/API-source/product metadata into these shapes; Authoring
still imports no UWS, apitools, OpenUdon, or Ramen packages.
M16 adopts Authoring in OpenUdon phase 1: OpenUdon's internal prompt session,
structured JSON helpers, atomic lifecycle writes, and prompt-safe API context
translation now delegate to public Authoring packages. OpenUdon keeps its
product-specific progressive-loop sequencing, prompts, catalog planning,
workflow intent schema, repair rules, transcript/report wire versions, and
package artifacts downstream.
M17 adopts Authoring in OpenUdon phase 2 by validating OpenUdon agent,
retention, and scorecard reports through public `report` and `readiness`
contracts. OpenUdon still emits its own `openudon.icot-*` JSON versions and
keeps fixture corpus, expected top-issue policy, authoring-eval categories,
and report bodies downstream.

## Data Flow

OpenUdon flow:

```text
project brief + API metadata
  -> authoring session
  -> OpenUdon runtime draft workflow intent
  -> OpenUdon review/repair
  -> OpenUdon package artifacts
```

OpenUdon adoption keeps OpenUdon-specific session schemas, prompts, workflow
intent, catalog planning, repair rules, review packages, and `openudon.icot-*`
wire versions in OpenUdon while replacing generic helpers with Authoring
packages.

Ramen flow:

```text
operator goal + API metadata + optional Terraform-like hints
  -> authoring session
  -> Ramen runtime draft desired-state intent
  -> Ramen review/repair
  -> native project.uws.yaml with x-ramen-desired-state
```

Ramen adoption starts with a spike that can draft a native project skeleton and
stop before graph/plan if required mapping metadata is missing. The spike lives
in `../ramen/authoring` and binds Ramen state to Authoring's runtime,
prompt-context, lifecycle, readiness, and report contracts; Authoring still
imports no Ramen, UWS, or OpenUdon packages. Later Ramen adoption expands
variables, resources, operation roles, identity fields, dependencies,
redaction hints, and `x-ramen-desired-state`, with validation, graph, and plan
checks remaining in Ramen.
M19 performs that downstream expansion: the Ramen adapter accepts Ramen-owned
variables and resources, translates prompt-safe schema hints into default
schema/identity/request-binding/redaction metadata, preserves explicit
operation roles and dependencies, and can run Ramen validation, graph, and plan
checks as gates. The shared Authoring module still owns only orchestration and
neutral result/report contracts.
M20 adds import-boundary and cross-repo compatibility gates: Authoring tests
prove the generic module does not depend on OpenUdon, Ramen, UWS, or apitools,
and `scripts/check-compat.sh` runs Authoring plus sibling OpenUdon/Ramen
workspace checks.
M21 documents the pre-1.0 compatibility policy and adds API-surface tests for
durable version constants and exported JSON tags.
M22 adds README/package documentation and runnable fake-client/fake-runtime
examples for prompting, structured JSON fallback, progressive loops, and
report/digest generation.
M23 adds internal maintenance hardening after the first broad adoption pass:
bottom-layer normalization helpers, internal artifact/digest record
normalizers, unified unknown-severity ordering, bounded JSON-fence extraction,
and public comment cleanup. It must not add exported APIs or product semantics;
OpenUdon and Ramen remain downstream consumers validated through compatibility
checks and module bumps after Authoring changes land.
M24 adds the public shared interactive iCoT loop API and `icotcli` flag
plumbing. The loop owns opening prompts, extractor draft/disambiguation hooks,
readiness/question sequencing, deterministic prefill, draft-question hooks,
autosave, transcript callbacks, cancellation, and needs-input behavior while
downstream products still own prompts, API metadata translation, provider
clients, artifact formats, and gates.
M25 replaces the two active question-at-a-time implementations with one
frontier-round engine shared by generic, runtime-bound, and interactive entry
points. `interview` owns `authoring.interview.v1`, graph validation,
deterministic dependency-ready frontiers, answers, deferrals, and a unified
public evidence ledger. The engine has no breadth ceiling, shows the whole
round before collecting answers, applies all answers before one normalization
and autosave, and diagnoses three consecutive no-progress rounds. Durable
session/transcript shapes remain unchanged where their JSON did not change;
product-specific v2 state and wire migrations stay downstream. Evidence may
carry normalized string attributes for public, machine-readable qualifiers
needed to reproduce downstream readiness and safety policy after resume.

I01 hardens that frontier architecture. `interview.ApplyRound` is the shared
complete-frontier transaction, and `icot.InterviewBinding` projects nodes and
applies product mutations on clones for generic, bound-runtime, and interactive
entry points. Loop progress uses an overridable semantic fingerprint, event IDs
are monotonic, transcript event order is authoritative, and zero `MaxRounds`
selects the 1,000-round emergency fuse. Prompt wording remains configurable;
OpenUdon retains its workflow label downstream.

D01 depends on I01's `session.ValidateForPersistence` contract. Prompt
transcripts and capable draft values reject explicitly sensitive unredacted
turns or answers before serialization. Lifecycle writes sync the temporary
file and parent directory around rename, structured fallback extraction stays
bounded, and report cancellation classification uses an internal sentinel
rather than error-message substrings.

M28 moves operation lifecycle ranking to
`github.com/OpenUdon/apitools/operationlifecycle`, where it consumes
`apitools.OperationSummary` and can require explicit Google Discovery
provenance before normalizing `/upload`. OpenUdon and Ramen translate their
product context downstream; Authoring no longer owns or imports the ranking
package.

## Security Boundary

- Authoring stores no credential values.
- Transcripts and model outputs are untrusted and may require downstream
  redaction before persistence.
- Explicitly sensitive turns and answers must be redacted before Authoring
  canonical serialization, prompt transcript save, or capable draft save.
- Generic orchestration must not perform live API calls, workflow execution,
  Terraform/OpenTofu execution, or trusted-runner invocation.
- Default tests must be provider-free and model-free.
- Downstream adapters decide whether model calls are enabled and supply any
  provider clients explicitly.
- Safe artifact paths, digest records, diagnostic records, and redaction
  helpers come from Evidence where the record is shared across products.
- Noninteractive agent results must distinguish `complete`, `needs_input`,
  `failed`, and `canceled` without implying product-specific remediation.

## Neutral engine ownership (M29 extraction, M30 removal)

`engine` owns the generic progressive, interactive, bound-runtime, repair and
atomic interview implementation. M29's alias/forwarder facade is removed in
M30, including `icotcli`; no neutral implementation body, signature or durable
wire changes. Old pinned Ramen continues standalone, but current workspace
compatibility is intentionally ended under the approved exception. Kinet and
OpenUdon use only surviving neutral packages. Exact-source qualification uses
disposable overrides, not consumer manifests or the operator's `go.work`.
Earlier M01–M29 paragraphs describe historical extraction names; `engine` is
the current owner of their loop and interview APIs.
