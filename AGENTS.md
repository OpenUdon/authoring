# AGENTS.md

## Purpose

Authoring is the public shared authoring engine for OpenUdon and Ramen. It owns
generic progressive iCoT/session orchestration, structured JSON completion
helpers, transcript/event contracts, readiness-loop primitives, and adapter
interfaces for product-specific artifact generation.

Module path:

```text
github.com/OpenUdon/authoring
```

Authoring must not own OpenUdon workflow-package semantics, Ramen desired-state
profile semantics, UWS workflow semantics, API source parsing, credential
resolution, or live execution.

## Memory Bank First

The tracked canonical Authoring harness snapshot lives in `../tofu/authoring`.
In a normal `../authoring` checkout, `AGENTS.md`, `memory-bank/`, and
`evolution/` may be symlinks to this tracked snapshot so agents can use local
paths while planning history is committed in the `../tofu` repository.

Before substantial changes, read in this order:

1. [memory-bank/product.md](memory-bank/product.md)
2. [memory-bank/architecture.md](memory-bank/architecture.md)
3. [memory-bank/tech-stack.md](memory-bank/tech-stack.md)
4. [memory-bank/milestone.md](memory-bank/milestone.md)
5. The relevant per-milestone status file in [memory-bank/](memory-bank/)

Use the memory bank as the active project source of truth. Do not recreate
duplicate root-level product, architecture, roadmap, or aggregate status
documents.

This project exposes [GOAL.md](GOAL.md), one optional protocol for goal requests
that span multiple status files. Follow it only when a request names it.

A `GOAL.md` run is a deliberate exception to the row-level commit rule below.
For that run, `COMMIT_POLICY: none` — the protocol default — means no commits,
while `COMMIT_POLICY: task` keeps the usual one-commit-per-row cadence.
Precedence is the request, then `GOAL.md`, then this file; only commits are
delegated, and only during the run.

## Boundaries

- `../evidence` owns generic digest, artifact safety, diagnostic, redaction,
  and approval primitives. Authoring may consume those primitives but must not
  redefine durable trust evidence.
- `../uws` owns public workflow semantics, document model, schema lookup, and
  validation. Authoring may be used by UWS-producing adapters, but generic
  Authoring packages should not add UWS semantics.
- `../apitools` owns API source discovery, metadata, operation summaries,
  auth/security summaries, and ranking. Authoring may consume prompt-safe API
  metadata through downstream adapters.
- `../openudon` owns OpenUdon prompts, workflow intent schema, synthesis,
  package quality, review-handoff artifacts, and trusted-runner package layout.
- `../ramen` owns Ramen authoring adapters, Terraform-like desired-state
  ergonomics, native Ramen/UWS project generation, validation, graphing,
  planning, state, and reconciliation.

Rule of thumb:

- If it is a generic loop, session, transcript, structured-output, or adapter
  interface useful to both OpenUdon and Ramen, it can belong in Authoring.
- If it names OpenUdon package paths, Ramen resources, UWS fields, API source
  families, provider mappings, credentials, accounts, or execution behavior, it
  belongs downstream or in the owning sibling module.

## Execution Model

Authoring follows the same bound-runtime pattern used by UWS.

The upstream package owns shared structs, orchestration, validation, stable
state transitions, and deterministic algorithms. Runtime-dependent or
product-dependent behavior is expressed through narrow interfaces that
downstream products implement and bind at execution time.

Example shape:

```go
type Runtime interface {
    Draft(ctx context.Context, session Session) (Draft, error)
    Review(ctx context.Context, draft Draft) ([]Diagnostic, error)
    Repair(ctx context.Context, draft Draft, issues []Diagnostic) (Draft, error)
    WriteArtifacts(ctx context.Context, draft Draft) ([]Artifact, error)
}
```

The Authoring engine should drive the loop; OpenUdon and Ramen adapters should
own prompts, domain readiness checks, draft schemas, repair rules, and artifact
writers. Prefer composition and explicit interface binding over inheritance-like
type embedding that hides product behavior.

## Commands

Initial harness/documentation checks:

```bash
git -C ../tofu diff --check -- authoring
git -C ../authoring status --short
```

Planned public module checks once Go code exists:

```bash
go test ./...
go vet ./...
git diff --check
```

When exported APIs change, run dependent checks in sibling consumers as
applicable:

```bash
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
```

## Safety

- Treat model prompts, model outputs, transcripts, answers, generated drafts,
  generated artifacts, and session files as untrusted until validated.
- Do not store secrets in prompts, transcripts, reports, examples, or generated
  artifacts. Use symbolic credential binding names only.
- Do not execute workflows, API operations, Terraform/OpenTofu behavior, or
  trusted-runner actions from Authoring.
- Keep model-provider calls optional and downstream-controlled. Generic tests
  must run without live model credentials.
- Redact provider output before durable storage when downstream adapters opt
  into transcript or report persistence.

## Documentation Rules

- Update [memory-bank/milestone.md](memory-bank/milestone.md) when milestone
  scope, sequencing, acceptance criteria, current-state dashboard, boundaries,
  or the status-file index changes.
- When a milestone has multiple implementation tasks, update the matching
  `memory-bank/status-<LANE><NN>.md` file with task rows, state, notes, and scoped
  commit tracking.
- Keep one permanent, zero-padded status file for every indexed milestone.
  Never reuse an ID or create aggregate `status.md`.
- Keep candidate directions unnumbered until fresh scope and dependency review
  promotes them.
- Write task ledgers as `Item | State | Notes` with a backticked marker in the
  second column: `` `[ ]` ``, `` `[+]` ``, `` `[~]` ``, `` `[!]` ``, or
  `` `[X]` ``.
- Treat each row as a commit unit. Parallel interview and durable-lifecycle
  work requires explicit non-overlapping package ownership, resolved
  prerequisites, and downstream impacts in `milestone.md`.
- Update [memory-bank/product.md](memory-bank/product.md) when product scope,
  users, workflows, concepts, or non-goals change.
- Update [memory-bank/architecture.md](memory-bank/architecture.md) when system
  boundaries, package layout, data flow, execution model, or security
  boundaries change.
- Update [memory-bank/tech-stack.md](memory-bank/tech-stack.md) when
  dependencies, commands, runtime assumptions, artifact schemas, or tooling
  choices change.
- Check [evolution/](evolution/) after a major review, milestone, or boundary
  change. Add the next prompt/result version only when product direction,
  architecture boundary, milestone target, or public/private contract direction
  materially changes.
