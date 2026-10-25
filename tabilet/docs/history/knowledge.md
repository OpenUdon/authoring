# Superseded knowledge

## 2026-10-15 — Authoring workflow files are owned by this checkout

**Sources.** `AGENTS.md` / Memory Bank First and Commands;
`tabilet/memory-bank/milestone.md` / Current State;
`tabilet/memory-bank/tech-stack.md` / Commands.

**Reason.** The previous guidance describes an earlier sibling-owned snapshot.
Current read-only inspection finds the workflow files tracked directly by this
Authoring repository, with ordinary local files and directories. The referenced
`../tofu/authoring` directory is absent. This correction changes current workflow
navigation and commands; it does not revise completed task evidence.

**Evidence.** Observed repository HEAD
`acacd51ad3ea8141b48399e9976f6ef95ae2c148`, including the uncommitted v2.13.0
workflow upgrade. `git rev-parse --show-toplevel`, `git ls-files -- AGENTS.md
tabilet/memory-bank tabilet/evolution`, and filesystem inspection establish the
current ownership. This observation does not claim the uncommitted correction
is present in that commit.

**Replacements.** [Local workflow ownership](../../../AGENTS.md#memory-bank-first),
[local checks](../../../AGENTS.md#commands),
[current milestone state](../../memory-bank/milestone.md#current-state), and
[current stack commands](../../memory-bank/tech-stack.md#commands).

### Previous AGENTS.md / Memory Bank First

```markdown
The tracked canonical Authoring harness snapshot lives in `../tofu/authoring`.
In a normal `../authoring` checkout, `AGENTS.md`, `tabilet/memory-bank/`, and
`tabilet/evolution/` may be symlinks to this tracked snapshot so agents can use local
paths while planning history is committed in the `../tofu` repository.
```

### Previous AGENTS.md / Commands

```bash
git -C ../tofu diff --check -- authoring
git -C ../authoring status --short
```

### Previous milestone.md / Current State

```markdown
The repository exists
as a public Go module scaffold, and
`AGENTS.md`, `tabilet/memory-bank/`, and `tabilet/evolution/` are symlink-facing paths to the
tracked snapshot under `../tofu/authoring`.
```

### Previous tech-stack.md / Commands

```bash
git -C ../tofu diff --check -- authoring
```

## 2026-10-25 — Completed milestone retirement

**Sources.** `AGENTS.md` / Documentation Rules;
`tabilet/memory-bank/milestone.md` / Current State, Delivery Strategy, and
Long-term memory and retirement; `docs/m30-handoff.md` / source-record location.

**Reason.** The user explicitly requested retirement of completed milestones.
M01-M30, I01, and D01 have closed task tables; M29 and M30 have recorded passed
review 1/10 and must not be relabelled legacy. M31 remains unretired because its
accepted summary conflicts with duplicate pending rows, an older review count,
and qualification limitations. No completed task or review evidence is changed.

**Evidence.** Observed full HEAD `be023f06d0837b58a651449b6d89a00c12eb04cd`, including the existing
uncommitted approved workflow upgrade and this retirement. No claim is made
that these working-tree edits exist in that commit.

**Replacements.** [Current state](../../memory-bank/milestone.md#current-state),
[delivery strategy](../../memory-bank/milestone.md#delivery-strategy),
[retirement rules](../../memory-bank/milestone.md#long-term-memory-and-retirement),
and [retired milestone index](index.md). Current product, architecture, stack,
lessons, and all evolution snapshots are preserved. The historical source
specifications and statuses are retained in full in their individual records.

### Previous AGENTS.md / completed-ledger exception

````markdown
  Milestones M01-M30, I01, and D01 closed before this procedure. They keep their
  ledgers and specifications in place under the existing convention until an
  explicit cleanup request retires them.
````

### Previous milestone.md / Current State

````markdown

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
with the compatible iCoT facade. M30 is complete and accepted after review 1,
retiring that facade with Ramen frozen at its pinned Authoring.
The exact accepted source and qualification are recorded in
[status-M30.md](status-M30.md).
M31 is complete and accepted: Ramen adopted the
published neutral engine, imports `authoring/engine` instead of the retired
facade, and owns its interactive iCoT CLI flags locally; see
[status-M31.md](status-M31.md).
Authoring is a public Go module. Its `AGENTS.md`, `tabilet/memory-bank/`, and
`tabilet/evolution/` are tracked directly in this repository. Workflow and
planning changes use these local files and this checkout's Git commands.

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
````

### Previous milestone.md / Delivery Strategy

````markdown

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
````

### Previous milestone.md / pre-adoption retirement rule

````markdown
Milestones M01-M30, I01, and D01 closed before this procedure and keep their ledgers and
specifications in place under the existing completed-ledger convention. Retiring
them needs an explicit cleanup request and, because they carry no recorded review
count, the legacy-closure form below.
````

### Previous docs/m30-handoff.md / source-record location

````markdown
[M30 status](../tabilet/memory-bank/status-M30.md) is the package-local accepted
record; it stays in Authoring's completed ledger under the established convention.
````
