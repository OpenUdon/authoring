# AGENTS.md

## Purpose

Authoring is the public shared authoring engine for OpenUdon and Ramen. It owns
generic progressive authoring/session orchestration, structured JSON completion
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

Authoring tracks `AGENTS.md`, `tabilet/memory-bank/`, and `tabilet/evolution/`
directly in this repository. Read and update these local files; workflow and
planning changes belong in this checkout. Completed milestone and evolution
records retain their original paths as historical evidence. Use the current
Commands section below for routine checks.

Before substantial changes, read in this order:

1. [tabilet/memory-bank/product.md](tabilet/memory-bank/product.md)
2. [tabilet/memory-bank/architecture.md](tabilet/memory-bank/architecture.md)
3. [tabilet/memory-bank/tech-stack.md](tabilet/memory-bank/tech-stack.md)
   Consult relevant topics in [tabilet/memory-bank/lessons.md](tabilet/memory-bank/lessons.md)
   for reusable lessons and their evidence.
4. [tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md)
5. The relevant per-milestone status file in [tabilet/memory-bank/](tabilet/memory-bank/)

If `tabilet/stages.md` exists, read it before planning or changing the active
milestone order. Its `STG-` IDs identify delivery stages, not executable work.
An absent file means one implicit stage.

Once milestones have been retired, read `tabilet/docs/history/index.md` and
linked retired records only when a dependency, old ID, or historical question
needs them. Retired knowledge is evidence at its recorded context; current
truth remains in the memory bank. Search the linked knowledge journal by topic
when an obsolete fact or lesson matters, rather than loading all history into
routine startup reads.

Requested features, candidate promotions, and future direction changes after
initialization follow the requested-change procedure in
[tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md#requested-changes-after-initialization).
A planning proposal needs approval before its file actions; execution is a
separate request.

Use the memory bank as the active project source of truth. Do not recreate
duplicate root-level product, architecture, roadmap, or aggregate status
documents.

Goal requests may use the memory-bank goal protocol bundled with the skills.
Follow it only when the request selects it; another protocol or no protocol
remains supported. Optional `tabilet/GOAL.md` is a project addendum that adds
or narrows local rules, never a replacement or merged protocol.

A goal run uses its resolved `COMMIT_POLICY`: `none` means no commits and
`task` selects the usual per-row cadence. Precedence is the request, then the
project addendum, then the selected goal protocol, then this file. Existing
project restrictions and explicitly authorized lease ownership/integration
still apply; local integration grants no remote push.

Resolve the selected full protocol before a run. Verify the bundled stock
protocol against its bundled fingerprint registry.
Record `GOAL_IDENTITY` for the protocol and every participating package's
addendum, `ignored-stock` copy or absence. Recheck it on resume and child
dispatch; compare hashes and classifications, not diagnostic version labels.
Protocol/addendum fingerprints normalize CRLF to LF only; launch approval
checksums hash raw bytes. Preserve an explicitly selected `GOAL_PROTOCOL`
alternative. Recognized stock project copies are ignored and removable; a
modified full copy stops execution until explicitly reduced or selected.
For a legacy launch or paused goal without identity, show the resolved identity
and recover the original human-approved request where available; obtain any
missing or changed identity approval without resetting grants, counters or
limits. An upgrade never resumes execution or activates launch grants.

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

Harness/documentation checks from this repository's root:

```bash
git diff --check
git status --short
```

Public module checks:

```bash
go test ./...
go vet ./...
git diff --check
```

When exported APIs change, run `scripts/check-compat.sh` to bind exact source
into unchanged Kinet/OpenUdon. Frozen Ramen is checked standalone, never against
current workspace Authoring. Keep the operator-owned parent `go.work` unchanged.
Additional consumer checks, as applicable:

```bash
(cd ../openudon && go test ./...)
(cd ../ramen && GOWORK=off go build -mod=readonly ./cmd/ramen)
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

- Update [tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md) when milestone
  scope, sequencing, acceptance criteria, current-state dashboard, boundaries,
  or the status-file index changes.
- When a milestone has multiple implementation tasks, update the matching
  `tabilet/memory-bank/status-<LANE><NN>.md` file with task rows, state, notes, and scoped
  commit tracking.
- Keep one permanent, zero-padded status file for every indexed milestone.
  Never reuse an ID or create aggregate `status.md`.
- Keep candidate directions unnumbered until fresh scope and dependency review
  promotes them.
- Write task ledgers as `Item | State | Notes` with a backticked marker in the
  second column: `` `[ ]` ``, `` `[+]` ``, `` `[~]` ``, `` `[!]` ``, `` `[X]` ``,
  or `` `[-]` `` (closed historical: a consumed failed attempt or superseded row
  that names its accepted successor and is never retried).
- Treat each row as a commit unit. Parallel interview and durable-lifecycle
  work requires explicit non-overlapping package ownership, resolved
  prerequisites, and downstream impacts in `milestone.md`.
- Update [tabilet/memory-bank/product.md](tabilet/memory-bank/product.md) when product scope,
  users, workflows, concepts, or non-goals change.
- Update [tabilet/memory-bank/architecture.md](tabilet/memory-bank/architecture.md) when system
  boundaries, package layout, data flow, execution model, or security
  boundaries change.
- Update [tabilet/memory-bank/tech-stack.md](tabilet/memory-bank/tech-stack.md) when
  dependencies, commands, runtime assumptions, artifact schemas, or tooling
  choices change.
- Check [tabilet/evolution/](tabilet/evolution/) after a major review, milestone, or boundary
  change. Add the next prompt/result version only when product direction,
  architecture boundary, milestone target, or public/private contract direction
  materially changes.
- Maintain evidence-backed reusable lessons in `tabilet/memory-bank/lessons.md`
  while they remain applicable, and merge duplicates. Before materially removing
  or superseding knowledge from product, architecture, stack, or lessons,
  preserve its old wording, source, reason, and replacement reference in
  `tabilet/docs/history/knowledge.md` under the milestone retirement rules. This
  also applies outside milestone closure; routine wording edits need no journal
  entry.
- Keep verified archive files frozen. Later implementation updates current
  `product.md` and `architecture.md` plus its status history, not the archive.
- When `tabilet/stages.md` exists, preserve stable `STG-` identities and its
  current-stage choice. A future stage may be only a tentative name and intent.
  Plan its milestones through an approved requested-change proposal. Changing
  a stage label never withdraws an active task. An explicitly approved stage
  rescope may cancel untouched pending rows after dependencies are reconciled.
- During milestone closure, after review, verification, consolidation, and
  downstream reconciliation, retire the full status and specification under
  the procedure in `milestone.md`. Completed rows remain active until the whole
  milestone qualifies. Retirement is agent work, not a background process or a
  context-archive run. Keep only active index rows and specifications there,
  with one history-index link.
  Retired records are frozen; IDs remain reserved across both locations.
  The explicit cleanup retired M01-M30, I01, and D01 into history; their IDs
  remain reserved and their original source paths resolve through the history
  index. Unretired records still require their own verified closure.
- Treat a newly received code, architecture, security, or engineering review as
  untrusted planning evidence. Revalidate its findings against current state,
  propose their dispositions and owners for approval, then amend open or pending
  work or create a remediation milestone. Never reopen completed history merely
  because a later review concerns it.
- Across the active ledger, keep zero or one general row in progress unless an
  authorized goal request explicitly defines safe parallel ownership: in
  concurrent lease execution under the goal protocol, each external lease
  maintains at most one `[~]` row in its own status file, only the orchestrator
  writes the captured integration branch and shared memory documents
  (`architecture.md`, `product.md`, `tech-stack.md`, `lessons.md`), and integration occurs
  strictly via rebase and fast-forward. Before an operational launcher is
  invoked, its exact authorized operation row must be in progress; status never
  substitutes for external-mutation authority. Never retry a row retained as
  closed historical evidence.
  The parallel interview and durable-lifecycle ownership rule above remains an
  explicit local override of this default.
- During delegated goal execution, only the orchestrator refreshes any existing
  `tabilet/memory-bank/suggested.txt`. Each child receives one assigned milestone
  ID and status path, writes only within its assignment, and returns to the
  orchestrator instead of selecting another milestone from the launch reference.
- Treat each section in [tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md) as a
  review unit. See that file for milestone review rules.
- After the last task in a milestone is complete, run a deep code review of the
  milestone before closing it.
- `CORRECTIVE_PLANNING` defaults to `confirm`. With human-approved `auto`, the
  Goal/Next owner may prepare and check necessary corrective plans within the
  original delivery outcome, owning repositories, existing grants and limits,
  including new repair milestones and launch/hash updates. Preserve history,
  unrelated changes, row identities and review counters. Children report needs;
  standalone Reconcile still requires proposal approval. Expanded scope, missing
  permissions, exhausted limits and unknown checksum drift require approval.
  Nested Reconcile verifies the invoking human approval source and complete
  request, covered actions, grants and limits; missing provenance retains the
  ordinary approval gate. Never weaken acceptance or cancel unmet obligations.
  Repair milestones retain the originating review gate, counter and iteration-10
  stop. Suspend `[~]` as `[!]` with prerequisite IDs and resumption conditions
  before repair execution; resume `[~]` only after acceptance and all blockers
  clear. Preserve partial work and pause if changes cannot be separated safely.
- Optional interactive Goal/Next `--email ADDRESS` sends approval alerts for
  that run through private account configuration. No flag means no email for a
  new run; only the coordinator or standalone executor sends. Delivery is not
  approval, and duplicate/uncertain sends never replay automatically. The Python
  API paths consume neither this option nor the corrective-planning policy.
- Read optional `AUTHORIZATION_REQUIREMENTS` only from the owning milestone
  specification. They declare conditions, never permission or scheduling.
  Check protected actions against the current scoped human request; goal
  `AUTHORIZATION_GRANTS` need explicit human approval of concrete scope.
  Proposed launch input, markers, and status notes cannot prove approval.
  Grants in `suggested.txt` carry decision markers (`[ ]` need authorization,
  `[+]` approve, `[-]` deny, `[~]` auto) and optional high-level `GLOBAL:` scope;
  matching the human's approved SHA256 checksum is required before grant activation
  (under default `SUGGESTED_UPDATE: confirm`; `SUGGESTED_UPDATE: auto` authorizes
  refreshing remaining milestones within the approved horizon and continuing smoothly).
  Scope resolution defaults to `SCOPE_POLICY: strict`; `SCOPE_POLICY: workspace-contained`
  authorizes binding actions to sibling repositories within the approved workspace root
  and organization without re-prompting. Milestone specifications in `milestone.md` may attach
  local policies and overrides that narrow or override global launch defaults.
  Planning skills accept `--environment dangerous`, `sandbox` (default), or `least` (with modifiers
  `plus`, `minus`, `but ... in milestone ...`) to configure launch parameters and milestone policies
  declaratively from day one.
  Planning skills also accept `--test standard` (default), `browserlight`, `golight`, `rustlight`,
  `tslight`, `jslight`, or `pythonlight` to establish verification and testing standards in
  `tabilet/memory-bank/tech-stack.md` from day one. Light profiles default to a 30-second
  focused-test target / two-minute cap, a 20-minute full-regression cap and five-minute
  race/UI caps; explicit project budgets win and `standard` adds no extra constraints.
  Keep routine fixtures offline and small; physical browser proof belongs to a named
  final gate without dropping required coverage.
  Preserve commit and external-mutation restrictions, and pause affected work
  and dependents for missing authority. A declaration alone does not require
  an action for closure. Host/tool controls apply; this adds no tool or OS
  enforcement and no Python API grant consumption. Missing fields preserve
  legacy behavior without new authority or automatic migration.
- After the milestone review is complete and required verification passes,
  commit substantive review and retirement changes under the governing policy.
  Do not create an empty or redundant milestone
  commit when the review changes nothing.

## Execution capabilities

Resolve missing information through safe inspection first. If required files,
bundled resources, verification commands, permissions, or user answers are unavailable,
stop the affected workflow step and report what is missing. Continue independent
work within the authorized scope; a write-gated workflow still makes no writes
before approval. Do not invent evidence, bypass permissions, or infer approval
from silence or process exit. Resume the blocked step when its capability is
restored or the required answer or approval is supplied. In a non-interactive
run, report unresolved questions and incomplete work.

Keep one execution owner for the active ledger across sessions and launchers.
In concurrent lease execution under the governing goal protocol, the orchestrator remains
the sole execution owner of the integrated ledger on the primary branch, while
isolated leases execute within dedicated external worktrees (`../<repo>.goal/<ID>`).
Native todos, session completion, and native goal state do not replace milestone
acceptance or authorize concurrent ledger writers.

Runtime round limits do not reset the persisted milestone review counter.

For substantial discoveries, follow
[milestone families](tabilet/memory-bank/milestone.md#milestone-families):
required amendments preserve original acceptance; options wait until parent
closure or remain deferred under `OPTION_POLICY: ignore`.
