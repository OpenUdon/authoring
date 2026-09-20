# Product

Authoring is a shared public authoring engine for OpenUdon and Ramen. It
provides reusable progressive authoring primitives without taking ownership of
either product's domain semantics, artifact schemas, or execution policies.

The module exists so OpenUdon and Ramen can use the same iCoT/session,
structured-output, transcript, readiness, artifact-safety, report, and
repair-loop foundation while binding different downstream runtimes:

- OpenUdon uses Authoring to move from project briefs and API source metadata to
  workflow intent, UWS workflow artifacts, quality reports, review evidence,
  and trusted-runner handoff packages.
- Ramen uses Authoring to move from operator goals and Terraform-like
  desired-state ergonomics to native Ramen/UWS projects with
  `x-ramen-desired-state`, validation, graphing, and planning.

Authoring does not replace OpenUdon synthesis or Ramen project/profile logic.
It gives both products a common loop and shared contracts for interactive and
agent-assisted authoring.

The first extraction source is OpenUdon's mature `internal/authoring`,
`internal/icot`, and `internal/workflowintent` implementation. Authoring should
extract only generic mechanics from that stack: prompt modes, replay scripts,
transcript/session contracts, lifecycle persistence, structured JSON helpers,
progressive loop sequencing, decision evidence, readiness/question planning,
repair orchestration, report metadata, and scorecard primitives. OpenUdon
prompts, workflow intent schemas, API-source catalog behavior, review package
layout, trusted-runner policy, and OpenUdon wire versions stay downstream.

Ramen adoption comes after the OpenUdon extraction proves the shared contracts.
Ramen uses Authoring through a downstream adapter to draft native Ramen/UWS
project skeletons and desired-state metadata, but Ramen keeps resource mapping,
validation, graphing, planning, state, governance, and reconciliation.

## Users

- OpenUdon maintainers extracting generic iCoT and structured-output behavior
  from OpenUdon internals.
- Ramen maintainers adding Terraform-like native project authoring without
  depending on OpenUdon internals.
- Tooling maintainers who need deterministic transcripts, readiness decisions,
  repair attempts, and generated-artifact summaries.
- Downstream product adapters that want to bind product-specific prompts,
  schemas, validation, and artifact writing to a shared orchestration engine.
- Test and scorecard maintainers who need provider-free fake runtimes,
  deterministic replay, variant summaries, and failure-family reporting.

## Core Workflows

1. A downstream product starts an authoring session from a user goal, existing
   files, selected API source metadata, or saved answers.
2. The shared engine maps decisions into a dependency graph and records
   session state, transcript events, readiness findings, model/provider
   provenance, a unified evidence ledger, diagnostics, and deterministic
   decisions.
3. A downstream runtime drafts a product-specific intent or artifact model.
4. The shared loop presents every dependency-ready question as one numbered
   frontier, then applies the complete answer/deferral resolution set through
   a clone-based transaction. Full mode asks all decisions, normal mode
   visibly accepts safe recommendations, and fast mode silently accepts them
   while still asking forced or missing decisions.
5. The downstream runtime reviews and optionally repairs the draft within
   bounded product-specific rules.
6. The downstream runtime writes product-owned artifacts through safe artifact
   path and digest records.
7. The shared engine returns a stable result summary, transcript metadata,
   diagnostics, generated artifact descriptors, report metadata, and digest
   sidecars for downstream review.

## Scope

- Domain-neutral session and transcript structures.
- Progressive iCoT/readiness loop orchestration.
- Versioned dependency-aware interview graphs with nodes, evidence, answers,
  deferrals, complete-round resolutions, iterative validation, deterministic
  frontiers, and status transitions.
- Structured JSON completion helpers and fallback parsing contracts.
- Generic draft/review/repair/write interfaces.
- Generic diagnostic, event, and result-summary shapes when they are not tied
  to product semantics.
- Local prompt modes, required/yes-no prompts, replay scripts, transcript
  save/load, configurable shared messages and labels, and deterministic
  prompt-label assertions.
- Generic draft lifecycle persistence, autosave, atomic writes, safe artifact
  path records, and digest summaries.
- Decision evidence, confidence policy, readiness and question planning
  containers, review/remediation containers, report metadata, and scorecard or
  variant-harness primitives.
- Narrow prompt-safe context interfaces for source documents, operation
  candidates, schemas, and ordered OR-of-AND symbolic credential binding sets,
  with downstream adapters translating product metadata into those shapes.
- Use of `github.com/OpenUdon/evidence` primitives for digests, artifact
  safety, diagnostics, and redaction wherever the durable record shape is
  shared.

## Non-Goals

- OpenUdon workflow intent schema, prompts, package layout, quality gates,
  review-handoff manifests, or trusted-runner behavior.
- Ramen desired-state project/profile schema, resource mapping, graphing,
  planning, state, governance, or reconciliation behavior.
- UWS document semantics or schema validation.
- API source parsing, cataloging, or auth/security metadata.
- API operation lifecycle ranking, which belongs in `apitools` with the source
  provenance needed to interpret provider-specific paths safely.
- OpenUdon API-source catalog planning or Ramen API-source/resource mapping.
- LLM provider ownership, credential storage, live model calls in default
  tests, or prompt content that stores secrets.
- Workflow execution, API execution, Terraform/OpenTofu execution, or trusted
  runner invocation.
- Reimplementation of durable trust evidence already owned by
  `github.com/OpenUdon/evidence`.

## Current State

Authoring has completed the M01 harness stage, M02 Evidence-backed foundation,
M03 public session/transcript contracts, M04 local prompt/replay primitives,
M05 draft lifecycle/artifact-safety helpers, M06 structured JSON output
helpers, M07 generic progressive iCoT loop, M08 bound-runtime adapter
interface, M09 decision evidence and confidence policy, M10
readiness/question-planning contracts, M11 draft review/repair orchestration,
M12 noninteractive agent-result contracts, M13 report/retention metadata
helpers, M14 scorecard and variant-harness primitives, M15 prompt-safe context
interfaces, M16-M17 OpenUdon adoption, M18-M19 Ramen adapter adoption, M20
cross-repo compatibility gates, M21 pre-1.0 public API stabilization, M22
documentation/examples, M23 internal normalization, M24 shared interactive
iCoT extraction plus `icotcli` flag plumbing, and M25 adaptive frontier-round
authoring with `authoring.interview.v1`, followed by M26 harness lanes and M27
`authoring.prompt-context.v2` security alternatives. I01 interview/session
integrity, D01 persistence/report hardening, and M28 API lifecycle-ranking
relocation are complete and published; OpenUdon and Ramen pin the coordinated
Apitools and Authoring revisions with passing standalone gates. The loop now shares one atomic interview
binding across generic, runtime, and interactive entry points, preserves event
chronology, rejects unredacted sensitive persistence, and uses a configurable
1,000-round emergency fuse. The roadmap
expands the original M02-M06 sketch while keeping OpenUdon and Ramen product
semantics downstream.
