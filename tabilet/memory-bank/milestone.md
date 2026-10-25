# Milestones

This file owns Authoring milestone sequencing, acceptance criteria, current
state, and the status-file index.

Each milestone section is a review unit: after no `[ ]`, `[~]`, or `[!]` rows
remain in its matching status file, run a deep code review and a milestone
review against the milestone acceptance criteria before moving to the next
milestone. Completed `[+]`, cancelled `[X]`, and closed-historical `[-]` rows are
non-actionable; every `[-]` row must name its accepted successor. Review-driven
fixes should be verified and committed before work starts on the next milestone.
Do not create an extra milestone commit when the review produces no changes.

The review procedure below covers one milestone. To run several in order, with
dependency and downstream reconciliation between them, the memory-bank goal
protocol bundled with the skills is optional; any equivalent works as well.
An optional `tabilet/GOAL.md` addendum narrows it without restating it.

When `tabilet/stages.md` exists, it records stable `STG-` IDs, the current
stage, and provisional later ideas. Give each active milestone specification
its stage ID. An absent stages file means one implicit stage. Stage labels do
not determine task readiness or prove milestone acceptance.

Each active milestone specification defines scope, acceptance criteria, and
optional boundaries and parallel safety: `Depends on`, `Downstream impacts`,
`Write set`, `Contracts read`, and `Parallel-safe` (defaulting to `no`). A
milestone is parallel-safe only when unordered against other active work, its
write set is disjoint, and its inputs are frozen. Parallel safety never prunes
downstream impact reconciliation.

## Milestone families

Keep the original delivery outcome stable. Small directly related fixes stay
in the current task. Substantial required repairs use amendment records;
optional features, broader coverage, and environment conveniences use option
records. An amendment must name the original unmet requirement, evidence that
it is unmet, and why a simpler implementation or verification route cannot
satisfy it. External testing-tool repair remains optional unless the owner
explicitly promotes it, even under `CORRECTIVE_PLANNING: auto`. A failed test
method alone does not establish a product defect. Mocks cannot establish
unexecuted real acceptance; an unmet original requirement needs repair or an
explicit owner disposition.

Use flat `amend-<parent>-<NN>.md` and `option-<parent>-<NN>.md` records under
`tabilet/memory-bank/`, following the child-ID rules below. Each child owns its
specification, tasks, and review in one file. Its first level-two heading begins
with its child ID, and exactly one unfenced `Parent:` declaration names the
root parent ID. Ordinary bold and trailing whitespace are accepted. Keep only
concise links in the parent and roadmap rather than duplicate specifications.
An initial option may be a short finding without rows; required amendments and
selected options need full scope, acceptance, tasks, and ordinary review.
Reclassification requires an explicit owner decision and a fresh identity
with preserved lineage, never a silent rename.

Finish independent parent work, then required amendments, then resume the
parent's deferred acceptance checks and final review. Required amendments
block parent acceptance; options do not. Preserve completed work and review
counters. Suspend affected parent rows as `[!]` with the prerequisite IDs and
resumption conditions. `Parent:` is a relationship, not a dependency on the
completed parent; do not create that circular edge. Each child has its own
review, which never replaces parent acceptance or resets its review counter.

For multi-milestone sequencing, follow the selected goal protocol's family
ordering rules. A strict order places required active amendments immediately
after their parent and before downstream milestones; family execution still
means parent work, amendments, then parent closure. Priority selection preserves
ordinary dependencies, ownership, and parallel-safety rules.

`OPTION_POLICY: confirm` is the default; `ignore` is the other value. During
parent execution, record optional discoveries briefly and defer detailed
planning, repeated inspection, and implementation. Under `confirm`, consolidate
and present one proposed implementation batch after parent closure, then wait
for an actual owner answer. Under `ignore`, defer refinement, prompts, and
execution throughout the goal. Preserve the backlog for later consideration.
Explicit option selection can override an earlier ignore setting for that
work. Approved options enter a newly approved execution horizon without
reopening the parent. Neither option policy nor launch-update policy supplies
new authority, budget, or scope.

The standalone API loop cannot select options automatically. An in-progress
option still counts toward sole execution ownership. Controller execution
requires an explicitly approved horizon binding each child's actual filename
in task and closure paths; the controller never inherits Goal grants or
expands its receipt dynamically. A suspended parent needs explicit
`resume_after` conditions naming its sole amendment blockers. Those conditions
gate pending and in-progress rows too; restore the parent only after amendment
acceptance and all other dependencies pass. Preserve grants, denials,
consumed attempts, and review counters.

Retire children under their original filenames, preserving their full source,
parent declaration, and permanent identity. Source status and Source
specification both name the original child file. Required amendment reviews
and acceptance must pass before parent retirement; deferred options do not
block it. Preserve earlier rows and specification prose. A declined rowless
option may retire as cancelled or explicitly superseded with its authorized
disposition, Review `not applicable`, and Review iterations `0`, without
inventing implementation rows. Supersession must identify the accepted
successor that retains outstanding obligations. All other retirement gates
remain in force.

## Authorization requirements

Each milestone specification may contain `AUTHORIZATION_REQUIREMENTS`. It
declares conditions on actions; it never grants permission or schedules them.
Keep the declaration in that specification only. Status files reference it.
The following optional shape is a template, not authorization:

```yaml
AUTHORIZATION_REQUIREMENTS:
  git.commit:
    via: goal-policy
    executor: assigned-agent
  git.push:
    via: explicit
    executor: coordinator
    scope:
      repository: "<owning repository>"
      remote: "<exact remote URL>"
      ref: "<exact destination ref>"
      mode: fast-forward
  cli.local:
    via: goal-policy
    executor: assigned-agent
    scope: declared-implementation-and-verification
  browser:
    via: explicit
    executor: assigned-agent
    scope:
      environment: fixture
      origins: ["<exact approved origin>"]
      actions: ["<approved action>"]
  sudo:
    via: explicit
    executor: coordinator
    scope:
      commands: ["<specific privileged operation>"]
  ssh:
    via: explicit
    executor: coordinator
    scope:
      host: "<exact host>"
      user: "<account>"
      commands: ["<specific remote operation>"]
```

`via` is `goal-policy` or `explicit`. Only local commits and ordinary local
implementation/verification may use `goal-policy`; their governing task scope
and commit policy still apply. Other listed actions require explicit human
authority. Inspect and resolve exact targets before seeking approval. A
resolved grant must have concrete scope covering the requirement, with no
unresolved placeholders, wildcard targets, blanket booleans, or credentials.

Under the optional bundled goal protocol, `AUTHORIZATION_GRANTS` belongs
to the resolved human-approved request, not to this file. Requirements,
proposed grants, status notes, and `suggested.txt` are evidence, not approval.
Check authority before the action; pause affected work and dependents when it
is missing. Independent authorized work may continue under the governing
order and ownership rules. Required unperformed task or acceptance actions
prevent closure; listing an action here alone does not require its performance.
Outside a goal, use the current scoped human request and applicable project
rules. Missing fields preserve legacy behavior without new authority or
automatic migration. These are instruction-level boundaries, not tool or OS
enforcement.

A changed input to a protected action changes its scope unless the governing
goal protocol's human-approved standing revision envelope covers it. That
envelope covers only bounded, unprivileged verification within its approved
workspace and isolation boundaries; it never extends sudo, SSH, live-browser,
or push authority.

## Standing corrective planning

`CORRECTIVE_PLANNING` is an optional execution policy, defaulting to `confirm`.
Standing `auto` cannot weaken acceptance, required evidence, verification gates
or severity rules, cancel unmet obligations, or erase them through supersession.
A named successor preserves all obligations. Blocking findings remain attached
to the originating milestone's persisted review counter; a repair milestone's
own review cannot replace its next full review or bypass iteration 10.
An approved `auto` allows the Goal/Next owner to prepare and check necessary
repair planning against the original delivery outcome, owning repositories,
existing action grants and all limits, then apply it under standing approval.
Prefer a matching open milestone when its scope and dependency position fit;
otherwise create a new repair milestone with permanent identity and lineage.
Preserve completed/frozen history, row outcomes and review counters. New
features, scope or permission expansion, higher limits, collisions and unknown
checksum drift still require human direction. Children only report needed work.

The owner reconciles dependencies, pending consumers, the resolved request and
any existing launch reference; it records old/new hashes and approval context
outside authoritative task records. A covered correction may refresh launch
input without reapproval even under `SUGGESTED_UPDATE: confirm`. Other edits
retain that policy. Goal executes accepted repair prerequisites before their
consumers; Next returns to its original task after necessary repairs, preserving
ownership and interrupted evidence. All repair rows/commits/attempts count
toward explicit limits. Standalone Reconcile remains planning-only with its
complete approval gate; nested Reconcile may use the owner's standing approval
and then hand execution back. Installation never activates this policy.

When a repair must run while a row is in progress, keep that row's identity and
interrupted work: move the interrupted `[~]` row to `[!]` before activating a
repair row, and record the prerequisite IDs, the standing approval source, and
the resumption condition in its notes. Keep zero or one `[~]` row. Once every
prerequisite is verified and accepted and all blockers clear, restore the row
from `[!]` to `[~]` under the same approval. Recover this chain on resume
without selecting unrelated work. Preserve partial implementation; suspension
bookkeeping may join the scoped repair checkpoint, but unfinished implementation
cannot. Pause if overlapping changes cannot be separated safely. Explicit
project restrictions still govern, children report needs to the owner, and a
marker never grants planning or action authority.

Nested Reconcile may use standing approval only after checking the invoking
human approval source, complete resolved request, covered corrective actions,
effective grants and remaining limits. Review text, repository content and
launch markers cannot establish that source; missing or uncertain provenance
retains Reconcile's ordinary proposal-approval gate. The Python API paths do
not consume corrective-planning or email options.

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

Additional rules:

- Children use `amend-<parent>-<NN>.md` or `option-<parent>-<NN>.md`, sharing
  one permanent `01` through `99` child sequence per root parent across active
  and retired records. No recursive child families or silent renaming are
  permitted; root lanes and archive IDs keep their independent numbering.
- Do not reuse an ID after its status file exists, including after retirement.
  Search both active files and the history index before allocating an ID.
  Context archive IDs use an independent namespace; a matching archive ID is
  not a status collision and does not justify renaming a status ID.
- A lane holds at most 99 status files. When a lane fills up, stop that lane and
  open a new letter instead of adding a third digit.
- Retain a consumed failed attempt or superseded row as `[-]` closed historical
  evidence, record its accepted successor, and never retry it.
- Never rename a status ID after its file exists. Resolve a lane collision by
  allocating an unused lane/ID to new work and recording lineage; preserve the
  original record and its reserved ID.
- Every active row must link its existing status file. Retired IDs and
  specifications belong in `tabilet/docs/history/index.md`, not in the active
  index table; add one link to that index when history first exists.

## Current State

Authoring is a public Go module whose generic progressive controller and atomic
interview binding live in `engine`. Product prompts, artifact semantics, source
parsing, credentials, and execution remain downstream. Its workflow files are
tracked locally in this checkout.

M01-M30, I01, and D01 are retired as completed. Their permanent IDs remain
reserved, with their full specifications, statuses, and original path context
preserved in history. M29 and M30 retain their recorded review 1/10; older
closures without a recorded counter use legacy metadata.

M31 remains active for closure reconciliation. Its status declares acceptance
and review 1/10 but also contains three duplicate pending rows, an older 0/10
review line, and a qualification caveat. This cleanup preserves that record
unchanged and makes no new acceptance claim or task selection.

## Delivery Strategy

Future implementation work follows the approved requested-change procedure and
the bound-runtime ownership rules. This retirement selects no implementation
task, promotes no candidate, and leaves M31 reconciliation unresolved.

## Status Files

Completed records and permanently reserved IDs are listed in the
[history index](../docs/history/index.md). Only unretired work is indexed here.

| Milestone | Status File | Summary |
|---|---|---|
| M31 | [status-M31.md](status-M31.md) | Retirement pending: accepted record contains conflicting pending rows, review counts, and qualification evidence. |

## Requested changes after initialization

For stage work, `memory-bank-propose --stages` requests creation or revision of
the optional `tabilet/stages.md` overview. `--stage STG-01` selects an existing
stage; the request must say when milestones should be planned. Equivalent
plain-English references are valid. Later-stage ideas may have only an ID,
tentative name, and brief intent. Additional context can be proposed later,
without assigning milestone IDs or task rows now. Stage IDs stay stable across
renames and reorderings and are never reused. Do not advance stages or schedule
later work merely because the current milestone closes.

For execution environment and authorization policy, `--environment dangerous`,
`--environment sandbox` (default when omitted), or `--environment least` sets the
hands-free stance in launch references and milestone policies, optionally modified
by expressions like `plus <action>`, `minus <action>`, or `but <action> in milestone <ID>`.

For a requested feature, candidate promotion, or change to future direction,
inspect the current plan and implementation first. Record user intent separately
from observed facts and unresolved assumptions. If a pending milestone already
owns a small change, propose the smallest row and acceptance update that fits.
Otherwise propose a dependency-closed milestone with an unused permanent ID.
A candidate's trigger prompts a new decision and approval; promotion is never
automatic. Schedule by approved priority and dependencies, not review finding
severity. Reconcile affected downstream acceptance and optional launch input.

Present one complete proposal with intended outcome, owner rows, acceptance,
verification, dependencies, downstream effects, and exact file actions. Write
planning changes only after approval and a fresh check of affected files,
worktree changes, and active and retired IDs. Preserve non-pending outcomes,
review counters, local policy, and frozen history. Record target behavior here
and in pending status rows until implemented; current architecture describes
only what the repository establishes as fact. Planning does not implement or
accept the requested behavior. Use `memory-bank-propose` where installed, or
follow this procedure directly with another agent.

**Expensive checks.** Some acceptance depends on a check that is expensive or
protected; see the [verification-attempt gate](#verification-attempt-gate).
For such a milestone:

- Limit acceptance to behavior the milestone owns, proved with the smallest
  workload. Put actual pairing with real producers or consumers in the
  milestone that adopts them.
- When the check needs host capabilities (privileges, operating-system or
  kernel features, devices, browsers, or external accounts), make a cheap
  readiness probe its first row.
- Name the uncounted local diagnostic loop, and order cases so that an
  unrelated failure cannot hide the milestone's own cases.

To defer already approved work through a stage rescope, present the exact
untouched pending rows and affected downstream work for approval. After
approval, retain their identity and original task text, mark them `[X]`, and
record the authority, reason, destination stage, and dependency disposition.
Cancellation proves no delivered acceptance. A future replacement receives a
new approved identity and names its cancelled lineage. Preserve `[~]`, `[!]`,
`[+]`, `[X]`, and `[-]` rows and all review evidence. Withdraw affected pending
dependents or keep their prerequisites in the active horizon; an unresolved
dependency cannot be deferred by changing a stage label. An entirely cancelled
milestone remains active until its normal closure and retirement disposition.

## Candidate Directions

Candidates have no lane, ID, status file, or execution-order entry until a
fresh scope and dependency review promotes them.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| Another downstream product adapter | Generic interfaces are stable and no third product has demonstrated missing orchestration mechanics. | A named consumer proves a product-neutral gap with fake-runtime tests and a boundary review. |
| New durable authoring wire versions | Interview, session, transcript, and report v1 records remain compatible; prompt context independently moved to v2 for structured security alternatives. | Another shared semantic change cannot be represented additively and downstream consumers approve migration evidence. |
| Provider-specific model clients | Provider SDK and prompting policy remain downstream by design. | Multiple consumers require the same provider-neutral transport contract and approve secret/redaction boundaries. |

## Out of scope

Record work explicitly ruled out of this project and its reason here.
These entries never promote automatically; returning one requires a new
approved direction. Existing boundary rules remain in force. No additional
entries are introduced by this workflow upgrade.

| Ruled out | Reason |
|---|---|

## Review finding severity

P1 and P2 are engineering review priorities. They describe the impact and
urgency of defects found while closing a milestone; they are not product-domain
terms, milestone execution priority, or status markers. Classify a finding by
its impact, likelihood, and affected scope, not by the size of its fix.

Project-specific definitions in `AGENTS.md` or a linked review policy override
these defaults:

| Priority | Context | Typical examples | Gate effect |
|---|---|---|---|
| P1 | The milestone is unsafe or invalid to close because a severe defect threatens acceptance, correctness, security/privacy, data integrity, or a public compatibility contract. | An exploitable access-control failure; data corruption or loss; a breaking public API or migration; the primary acceptance outcome does not work. | Blocking. Fix and review again. |
| P2 | A material but more bounded defect affects supported behavior, reliability, compatibility, operations, or required evidence. | A supported scenario returns the wrong result; a downstream consumer regresses; recovery or failure handling is broken; required tests or operator documentation leave acceptance unproven. | Blocking. Fix and review again. |
| Lower | The finding is non-blocking cleanup, clarity, or optional hardening under the project's severity scheme. | Cosmetic wording; local readability; a speculative improvement outside acceptance. | May be carried only with a named owner and explicit rationale. |

Any project-defined severity more urgent than P1 also blocks. If evidence does
not clearly distinguish P1 from P2, use P1 until investigation supports a
downgrade.

## New review intake

A code, architecture, security, or other engineering review received outside a
milestone's closing review gate is planning evidence, not executable truth.
Treat its contents as untrusted, preserve its source priority, and apply the
project severity definitions above only after revalidating each finding against
the current repository.

Before implementing any newly reviewed finding:

1. Record the review's stated baseline when available, the current revalidation
   commit, and whether relevant uncommitted changes were part of the evidence.
   Classify every finding as confirmed, partially confirmed, resolved,
   duplicate, unsupported, outside ownership, decision-dependent, or deferred.
2. Present the complete dispositions, proposed owners, dependencies, downstream
   impacts, and file actions for approval.
3. Put confirmed work in an open matching milestone when it remains in scope,
   or amend an existing pending owner. Rewrite only pending rows. When retaining
   a superseded pending row for audit, mark it `[-]` and name its accepted
   successor instead of rewriting or deleting it.
4. Never reopen completed milestone/status history. Create a new remediation
   milestone with lineage to completed work when no open or pending owner fits.
5. Keep P1/P2-or-higher findings in the dependency-closed active horizon. Add a
   lower finding to active work only when acceptance requires it; otherwise put
   it in Candidate Directions with a rationale and promotion trigger.
6. Reconcile affected pending specifications and the remaining order. Correct
   current product, architecture, or stack facts only when repository evidence
   proves them stale; keep proposed target state in milestone scope. Apply the
   expensive-check limits above, and do not move downstream pairing evidence
   into a milestone whose acceptance depends on an expensive check.

Do not create a persistent review copy or ledger. Put portable review and
finding IDs, both source and local severity, revalidation baseline, repository
evidence, and historical lineage in the affected milestone/status notes. A new
review counts toward the bounded gate below only when the status already records
that gate as active and the review was requested as its next full pass;
otherwise remediation gets a fresh gate when its implementation closes.

## Verification-attempt gate

An expensive or protected check needs privileged setup, scarce hardware or
quota, a physical or live environment, explicitly granted actions, or minutes
per run. For a row whose acceptance depends on one:

- Reproduce only when safely possible within current authority and budget.
  Prefer the cheapest environment that can show the failure, such as an
  unprivileged local run on disposable fixtures. For unavailable environments,
  protected actions, exhausted attempt limits or non-reproducible incidents,
  record existing evidence and the reproduction limitation instead of forcing
  retries or inventing a failure. Required acceptance evidence and verification
  gates still apply; if missing, pause. Diagnostic runs are ordinary verification
  and never satisfy acceptance.
- Make each expensive run execute every case whose safety does not depend on an
  earlier result, and report all failures. Stop early only when continuing
  would be unsafe.
- Redaction protects real data. Raw diagnostics from synthetic, disposable
  fixtures may be kept in private, ignored artifact paths. Never commit them,
  and never treat them as evidence of effects, non-dispatch, retry eligibility,
  or authority.
- Each expensive run is one attempt. Persist the row's attempt count in its
  status notes; a session change does not reset it. Stop after three
  consecutive failed attempts at the same stage (or the limit the milestone
  or a human-approved goal budget sets), or when the failing stage lies outside
  the behavior the milestone owns. Run the milestone's applicable cheap
  readiness/preflight checks before each counted attempt.
  Mark the row `[!]` with a short diagnosis and options, such as a correction, a
  rescope or an environment change, and ask the user. Approval to continue sets
  a new explicit limit. Never repeat an unchanged failed attempt.

## Milestone review procedure

When the last open row in a milestone's status file closes as `[+]`, `[X]`, or
`[-]` during an agent session, perform the review before ending the turn and
before moving to the next milestone. A `[-]` row counts as closed only when its
notes identify the consumed attempt or supersession and its accepted successor:

1. Re-read the milestone scope and acceptance criteria here. Confirm the code or
   docs meet the acceptance line; do not rely on the status file alone.
2. Run the required verification before review, including proportionate checks
   for affected consumers and any applicable compatibility, migration,
   rollback, security, concurrency, or failure paths.
3. Run a bounded deep-review and fix gate using the review finding severity
   context above.
   - The initial deep-review pass is iteration 1. Read the persisted count and
     findings first. Record the iteration as started before review; resume an
     interrupted pass at that same number. Runtime round limits are separate.
     Pin the milestone's base commit, inspect the full diff including staged,
     unstaged and relevant untracked work under the commit policy, read the
     `git log` range, and review the full milestone diff for correctness,
     regressions, failure semantics, boundary drift, stale docs, and missing
     tests. Review acceptance fidelity separately: requirements missing or
     partial, behavior nobody asked for, and requirements implemented wrongly,
     each quoting its acceptance line. Record findings grouped by concern.
     An empty implementation diff still needs acceptance review; do not
     manufacture changes. Lenses are coverage obligations: combine them or
     run batches within host capacity. When an interactive request selects
     `MODEL_POLICY`, follow the Goal/Next bundled model-policy reference for
     fresh read-only whole-milestone reviewers; the owner alone records results
     and fixes under `DEFAULT`. Absence preserves host model settings.
   - Record the iteration number and findings in the current status notes so a
     continuation cannot reset the counter. If the pass finds no P1, P2, or
     higher-severity issue, the gate passes. Otherwise, when the current
     iteration is below 10, fix every such finding in the current milestone,
     rerun affected verification, and review the whole milestone again,
     including the fixes.
   - Run at most 10 iterations; a session or reviewer change does not reset the
     counter. If iteration 10 still finds a blocking issue, add `[!]` review
     rows recording the remaining findings and iteration-limit blocker. Do not
     start another automatic fix-review cycle or move downstream; stop and ask
     for user direction.
   - Carry a lower-severity finding forward only with a named pending owner and
     explicit rationale.
4. Reconcile the memory bank. Update `product.md` if the milestone changed
   product scope, domain terminology, concept relationships, or business
   invariants. Update `architecture.md` or `tech-stack.md` if it changed
   boundaries, dependencies, commands, data flow, or runtime assumptions.
   Then run an environment retro on mistakes found during the milestone:
   mechanical issues become proposed deterministic checks; judgement issues
   become lessons or review rules. Include missing navigation or information and
   instructions that change nothing. Report these candidates by severity for a
   later approved proposal; create no unapproved task rows.
5. Check `tabilet/evolution/`. Add the next `prompt-vN.md` and `result-vN.md` only when
   product direction, architecture boundary, milestone target, or public/private
   contract direction materially changes.
6. Revisit candidate directions affected by the milestone. Update their reason
   or trigger; when a trigger is now true, propose a reconciled milestone and
   obtain approval before allocating its permanent ID and status file.
7. After the gate passes, reconcile affected active downstream dependencies and
   acceptance against the delivered outcome. Then consolidate and retire under
   the procedure below. During an ordered goal, the goal protocol finishes its
   downstream reconciliation before performing this retirement step.
8. Run required verification again, including the retirement links and record,
   then commit substantive closure changes under the governing commit policy.
   Do not create an extra milestone commit when the review changes nothing.
9. Report a short review summary: what was verified, the review-fix iteration
   count, what memory-bank files changed, any review commit, and whether an
   evolution bump was made. Include the retired record and durable lessons.

## Goal lease cleanup

When a goal uses external lease worktrees, follow the selected goal protocol's
cleanup gate. Remove one only after recorded successful closure, matching
ownership, no active writer, verified integration and a clean worktree. Inspect
ignored data as well as tracked changes. Before freezing closure records,
preserve required qualification/recovery artifacts in an approved durable
location with their hashes and citations. Never rewrite frozen citations or
delete their required evidence; preserve and report uncertain leftovers.
This rule does not authorize cleanup of existing leases during an upgrade.

## Long-term memory and retirement

The explicit cleanup request retired M01-M30, I01, and D01 on 2026-10-25.
Their original paths and complete source records remain available in history.
Use recorded review evidence where present; legacy closure applies only when
a milestone predates the persisted gate and no review count was recorded.

This project retires a milestone automatically after closure. “Automatically”
means the agent performs this procedure during the milestone review workflow
above, including when using `memory-bank-next` or `memory-bank-goal`. There is
no background process or separate archive invocation. Individual completed,
cancelled, or closed-historical rows remain in their active status file until
the whole milestone qualifies; unresolved work or missing closure evidence
keeps it active. Terminal rows alone do not prove acceptance.

Existing projects adopt this contract explicitly, with a compatible API runner
if used; upgrading installed skills alone does not merge project instructions
or move files. An explicit cleanup request may retire older closed milestones
only when their closure evidence is available, either a recorded review under
this contract or the legacy closure described below.

### Consolidate before retiring

The history index links retired milestones and superseded workflow knowledge.
An unretired record remains active until its own closure evidence qualifies it
for retirement.

Keep current facts in product, architecture, and stack documents. Maintain
[lessons.md](lessons.md) for applicable lessons and decision rationale with
evidence links; merge duplicates and avoid a chronological session log. Maintain
relevant lessons during ordinary work, not only at closure, and keep them active
while applicable even after their supporting milestone retires. This removes
accumulated history from active context, not a fixed number of tokens or files;
genuinely active work and useful knowledge can still grow.

Before materially superseding or removing knowledge from those documents,
append the old wording to `tabilet/docs/history/knowledge.md`. Use a unique, descriptive
dated heading, the original document and heading, retirement reason, supporting
evidence, and a link to its replacement (or an explicit reason there is none).
Preserve the old excerpt in a fenced markdown block. This applies outside
milestones too. Create this journal only when needed, link it from the history
index, and append corrections rather than rewriting old entries. Routine edits
need no journal entry; Git, when present, holds intermediate revisions.

Consolidation is ordinary project maintenance. A separate context snapshot is
optional and never a prerequisite for retirement; the archive skill's clean
baseline rule must not force commits or interrupt a no-commit goal.

### Retire a closed milestone

1. Require a passed review within the persisted 10-iteration limit, recorded
   verification, consolidated knowledge, and reconciled downstream work. No
   `[ ]`, `[~]`, or `[!]` row may remain. Every `[-]` row names its accepted
   successor. A cancelled or superseded milestone also needs an authorized
   disposition; it must not be described as delivered acceptance.
2. Create `tabilet/docs/history/status-<LANE><NN>.md` using the envelope below. Keep the
   full final specification and status document in separate literal markdown
   fences, including all task text, notes, acceptance, and review evidence.
   Choose fences longer than any fence in the source. Original relative paths
   inside these literal documents retain their source-document meaning.
3. Add one history index row using `Milestone | Outcome | Retired | Record |
   Summary`; the ID, outcome, date, and record link must match the envelope.
   Outcomes are `completed`, `cancelled`, or `superseded`. The record link is
   relative to the index, for example `[M01](status-M01.md)`.
4. Remove the active status file, specification section, and active index row.
   Before removal, validate every envelope field and compare both retained
   documents with their complete active sources. For Git evidence, use the
   literal full output of `git rev-parse --verify HEAD`, never an abbreviated
   hash from a log display. Stop with the active sources intact if validation
   fails; a malformed retired record is not completed closure.
   Keep one link from this file to `../docs/history/index.md`. Repair maintained
   incoming links and remaining dependencies. Never rewrite frozen snapshots;
   resolve their old paths through the record's original-path metadata.
5. Refresh an existing disposable goal suggestion to contain only remaining
   active work, or remove it when empty. Do not create one without a compatible
   goal protocol. The resolved goal in the conversation remains authoritative
   over its disposable launch input.
6. Verify the record and links before handoff. On interruption, reconcile the
   source and destination before continuing; duplicate IDs or a missing record
   are invalid, and an existing retired record must never be overwritten.

Use this envelope, with actual values. The `Evidence` field is the observed
full Git commit, not a claim that uncommitted work exists in that commit.
`Worktree` is `clean`, `includes uncommitted changes`, or `unversioned`; without
Git both provenance fields say `unversioned`. The date uses `YYYY-MM-DD` in UTC.
Closure metadata precedes the two sections and uses one line per field.

`````markdown
# Retired milestone M01 - <title>

**Milestone.** M01
**Outcome.** completed
**Retired.** <YYYY-MM-DD>
**Source status.** tabilet/memory-bank/status-M01.md
**Source specification.** tabilet/memory-bank/milestone.md#<original-anchor>
**Evidence.** <full commit or unversioned>
**Worktree.** <clean, includes uncommitted changes, or unversioned>
**Review.** passed
**Review iterations.** <1 through 10>
**Verification.** <commands, results, and supporting evidence>
**Consolidated into.** <current-document and lesson links, or no current-truth change>

## Milestone specification

````markdown
<Complete final milestone specification, not a summary.>
````

## Status record

````markdown
<Complete final status document, not a summary.>
````
`````

For a cancelled or superseded outcome, add `**Disposition.**` naming the
authority, reason, and dependency disposition. A superseded outcome also needs
`**Successor.**` identifying its accepted successor. These are milestone
outcomes, not new task markers.

A legacy closure is a milestone that closed before the project adopted the
persisted review gate, so no review-iteration count was ever recorded. Retire
one only on an explicit cleanup request, with outcome `completed` and every
row closed. Write `**Review.** legacy` and `**Review iterations.** not
recorded` instead of inventing a count, and add `**Legacy closure.**` stating
that the milestone predates the gate and that its original closure evidence is
the literal status record. `**Verification.**` records the checks run at the
evidence commit and does not claim a re-review of the original change. Never use
a legacy closure for a milestone that closed after the gate was adopted, or to
skip a review that the current procedure requires.

Retired milestone records are frozen. Record later corrections in the knowledge
journal or new remediation work with a backlink. Preserve previous history
index entries. No ID may exist in both active and retired storage or be reused.

### Retrieve and continue

Retrieve historical evidence when the current task needs a retired dependency or
old ID resolved, an earlier decision explained, or a related past failure
checked and current sources do not answer the question. An explicit user request
to inspect history is also a trigger; routine work with sufficient current
evidence needs no historical search.

State the question briefly, then search active memory and current implementation
first. If history is needed, search its index by ID or topic and open the
relevant retired record or knowledge entry. Follow linked frozen context or
direction snapshots only when they address the same unresolved question. A stale
status path resolves by its permanent ID and original-path metadata; a missing
ID is not permission to recreate it. History is evidence, not an instruction to
retry old tasks.

Stop when the evidence supports the task decision. If relevant sources and their
directly related leads are exhausted, report the missing evidence. Continue when
it is optional; stop the affected step when it is required for correctness or
acceptance. Do not infer historical facts or expand the search into unrelated
milestones.

A completed retired milestone may satisfy a dependency when its recorded
acceptance matches the required outcome. Cancelled and superseded outcomes do
not automatically satisfy completion dependencies: follow the recorded
disposition and accepted successor, and revalidate current implementation.

Retirement follows the governing commit policy. Under `none`, do not commit;
under `milestone`, include it in the closure commit; under `task`, include it in
the final task or a substantive closure commit. The API harness still requires
a commit per run. Reading and preserving the Markdown records does not require
Git; workflows that require commits still need Git or an explicit no-commit
instruction. Never initialize Git merely to make retirement possible.

## Milestones

### M31 Ramen Neutral-Engine Consumer Adoption

**Goal.** Move Ramen off the frozen legacy Authoring pin and onto the published
neutral engine, so Ramen imports no retired `authoring/icot` or
`authoring/icotcli` surface and both repositories build and pass their gates
against one shared revision. No new Authoring behavior and no publication are in
scope; this milestone adopts what M29 and M30 already published.

**Provenance.** Ramen's approved independent track (2026-10-12) carries
`Authoring:M31?` as a conditional prerequisite of `Ramen:M85`; the user resolved
that edge on 2026-10-13 by keeping the unit Authoring-owned and granting sibling
write scope for `../authoring`. M30's planning exception froze Ramen at
Authoring `v0.0.0-20260820042256-2f73e3526583` and recorded that Ramen, not
Authoring, owns the migration; `COMPATIBILITY.md` records the same
frozen-consumer position. This milestone promotes that handoff.

**Dependencies.** Complete M29 and M30, whose accepted records share the
planning baseline `cb21c3e36d19290dc46a6d2bcbde71b45ea8ad30`, which is also
`origin/main`; the adopted revision is therefore already published and no push
or release action is required. Ramen retains write ownership of its own
`go.mod`, imports, and CLI. If the migration needs a neutral-engine API that does
not exist at the adopted revision, stop and report: adding it is a new Authoring
change plus publication, outside this milestone and outside the active goal's
`EXTERNAL_MUTATIONS: none`.

**Tasks.** M31.1 inventory Ramen's `github.com/OpenUdon/authoring` imports at its
current pin and map each retired `icot`/`icotcli` use onto its neutral
`engine`/`interview`/`session`/`transcript`/`report` counterpart, recording an
explicit disposition (re-implement, replace, or drop) for any behavior without an
equivalent; M31.2 migrate Ramen onto the adopted revision, update
`go.mod`/`go.sum`, and remove every retired import while preserving the
operator-visible `ramen author` and `ramen icot` contracts and their gates;
M31.3 qualify Authoring and the migrated Ramen against the exact adopted revision
and record the adoption in both ledgers; M31.4 bounded deep review and downstream
handoff to Ramen M85.

**Acceptance.** Ramen's `go.mod` pins the adopted published Authoring revision
and neither `authoring/icot` nor `authoring/icotcli` appears anywhere in Ramen's
import graph. Ramen builds, vets, and passes its full standalone gate with
`GOWORK=off`; `ramen author` and `ramen icot` keep their documented behavior or
every change is recorded with its disposition. Authoring's own tests,
import-boundary checks, API-surface tests, and `scripts/check-compat.sh` pass with
the migrated consumer bound to the exact revision. No Authoring public package
name, signature, JSON tag, or durable version changes, and no publication occurs.

**Verification.** Authoring: `go test ./...`, `go vet ./...`, race tests, and
`scripts/check-compat.sh` with a temporary workspace or modfile override binding
the migrated Ramen and recording the effective revision for each consumer check.
Ramen: `GOWORK=off go build ./...`, `GOWORK=off go vet ./...`, `GOWORK=off go test
./...`, and an import grep for the retired packages. Fake providers and disposable
local fixtures only; no live model, account, or network service.

**Downstream.** Satisfies `Authoring:M31 -> Ramen:M85`; Ramen M85 and its
dependent M86 become dispatchable. Ramen's `tech-stack.md` records the adopted
pin and `product.md`/`architecture.md` note only what the migration changes. No
other consumer depends on this revision; OpenUdon and Kinet are unaffected.
