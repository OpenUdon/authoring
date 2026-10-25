# Retired milestone M30 - Retire the iCoT compatibility facade

**Milestone.** M30
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M30.md
**Source specification.** tabilet/memory-bank/milestone.md#m30-retire-the-icot-compatibility-facade
**Evidence.** be023f06d0837b58a651449b6d89a00c12eb04cd
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** Historical verification and its limitations are retained verbatim in the status record; original commands are not re-run or relabelled by this retirement. At the observed evidence revision, go test ./... and go vet ./... passed on 2026-10-25; literal source preservation and retirement-envelope validation passed before removal.
**Consolidated into.** Current product, architecture, stack, and lessons remain in tabilet/memory-bank/; no implementation or current product contract changes. Superseded roadmap and retention wording is preserved in knowledge.md.
**Specification encoding.** The added level-two identity heading is a retirement-parser wrapper; the following level-three section is the complete original source, retained verbatim.

## Milestone specification

````markdown
## M30 Retire the iCoT compatibility facade

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

**M30 completion (2026-10-02).** Accepted removal source
`db4f5193bc53819be2b1bf1d8728735941e28d9e`, qualified in
`9241c481a281ab31f3e4b6ae89f9242ae9517c6c`; review 1 passed with no open
P1/P2 findings. All four task rows and required verification are complete.
[Handoff](../../docs/m30-handoff.md) records unchanged Kinet/OpenUdon consumer
checks, frozen Ramen's original standalone pin, operator workspace consequence
and separate publication/adoption decisions. M30 retires the compatibility
surfaces; this repository retains its completed full specification/status under
its existing convention. Earlier planning-provenance counts/authority describe
the proposal baseline; the explicit Stage 6 goal supplied task-commit authority.
No push, tag, publication or sibling edit.

````

## Status record

````markdown
# Status M30 — Retire the iCoT compatibility facade

State: Complete; accepted, compatibility surfaces retired; publication separate.

## Goal and dependencies

Remove `authoring/icot` and `authoring/icotcli`, the compatible facade over
`authoring/engine`, so Authoring keeps one neutral implementation and no
iCoT-named public surface. Baseline: Authoring
`dc8f3d61970ae628fc0399b0ef42187aa62a3e5b` (M29 complete), clean worktree.
Ramen is frozen at its pinned Authoring `v0.0.0-20260820042256-2f73e3526583`
and is not updated. This package owns M30; Kinet coordinates the stage without
merging ledgers.

## Provenance

Kinet Stage 6 (STG-06), approved for planning by the user on 2026-10-02 (Kinet
decision R58). The user stated that udon-ui and Ramen are neither updated nor
coded and that Authoring may be edited. Source priority and an external review
baseline were not supplied: this is a requested stage, not review intake. The
full specification is in [milestone.md](milestone.md) (M30).

## Tasks

Markers: `` `[ ]` `` pending, `` `[~]` `` in progress, `` `[+]` `` complete,
`` `[!]` `` blocked, `` `[X]` `` cancelled, and `` `[-]` `` closed historical
(must name its accepted successor). At most one general row may be in progress.

| Item | State | Notes |
| --- | --- | --- |
| M30.1 Consumer inventory and Ramen-frozen record | `[+]` | Re-verify, by import-graph search, that OpenUdon, Kinet, W8M and every workspace module outside Ramen have no `authoring/icot` or `authoring/icotcli` import; record Ramen's pin and imports; record that Ramen retains the dependency while compatibility with current workspace Authoring intentionally ends; check frozen Ramen standalone with GOWORK=off and cached readonly resolution. No migration/drop claim; record any real evidence gap. Edit nothing in Ramen. |
| M30.2 Remove the facade | `[+]` | Remove both packages and their tests; update API-surface, boundary and examples tests, README, COMPATIBILITY, architecture and `scripts/check-compat.sh` so Ramen is no longer required to build against workspace Authoring; record the intentional pre-1.0 source break and approved current-workspace compatibility exception; update AGENTS/product/tech-stack current descriptions in this implementation row. Neutral packages unchanged. |
| M30.3 Qualify Authoring and unchanged consumers | `[+]` | Standalone and workspace full tests/vet and race tests; import-boundary check. Bind exact new Authoring through temporary workspace/modfile overrides for unchanged OpenUdon and Kinet, recording effective source resolution; Kinet normal checks alone use its old pin. Separately check frozen Ramen standalone. No consumer manifest or operator workspace edit. Prove source/fixture equivalence of the neutral engine. Note the `go.work` consequence for `./ramen` as a user action. |
| M30.4 Review and downstream handoff | `[+]` | Persist a bounded deep review (maximum ten iterations), record exact source and checks, and hand Kinet and OpenUdon the exact revision with no required consumer change. Publication needs a separate request. |

## Verification and acceptance

See the full specification in milestone.md. Default checks use fake providers
and disposable fixtures; no live model, account or network service. `engine`
keeps no downstream import. Prompt and transcript bytes, JSON tags, durable
versions, signatures of non-removed packages and state behavior are unchanged.

## Closing review

Persisted iteration count: 1/10 (passed). Task completion alone does not establish
acceptance. Record exact source, observed checks, findings and downstream
handoff before completion.

## Approval and execution authority

Planning only. The user approved the Stage 6 proposal on 2026-10-02, including
this Authoring milestone. Execution, task commits and any push, tag or module
publication need a separate request; none is authorized here. No Ramen, udon-ui
or W8M edit. `~/Workspace/go.work` is user-owned and outside this repository.

## Approved planning-handoff reconciliation — 2026-10-02

**Source.** Stage 6 planning handoff reconciliation, F04/F06; source priority/external review baseline not supplied, local P2. Revalidation `dc8f3d61970ae628fc0399b0ef42187aa62a3e5b` included the uncommitted M30 draft. The user approved the compatibility disposition, amendments and planning-file actions. Original clean-baseline wording describes the pre-draft snapshot. The four rows stay pending; the persisted review count remains 0/10.

**Compatibility exception.** Ramen still imports `authoring/icot` and `authoring/icotcli` at its old pin. Freezing it does not drop the dependency. M30 intentionally ends compatibility of those imports with current workspace Authoring, while leaving Ramen and its standalone dependency unchanged. Evidence: Ramen go.mod/imports, `COMPATIBILITY.md` and `scripts/check-compat.sh`. Its standalone build passed offline with installed Go 1.26.7 and private disk-backed temporary output after /tmp quota rejected linking. Recheck at execution; never replace an unavailable pin with current source.

**Qualification contract.** Temporary workspace/modfile overrides must bind exact new Authoring into unchanged Kinet/OpenUdon, with effective module directory/revision recorded. Kinet's ordinary GOWORK=off checks cannot establish that test. Preserve neutral wire, prompt/transcript, API and state behavior. Update current AGENTS/product/architecture/tech-stack and compatibility docs together only when removal actually happens.

**Coordination.** Recommended serial single-owner order W12 → U09 → Authoring M30. M30 has no Kinet code prerequisite, and its ledger remains authoritative. Return exact accepted source/checks for owner reconciliation; no consumer re-pin, publication or sibling edit is authorized. Direction: [prompt v8](../evolution/prompt-v8.md) / [result v8](../evolution/result-v8.md).

## Kinet downstream reconciliation — 2026-10-02

The confirmed serial Stage 6 goal authorizes task commits, no publication, and
execution of M30 after Kinet acceptance. W12 is accepted at
`31288de3be4e20def6296bf66e1985169b6fe09b`; U09 is accepted at
`45217237f62993c278ccdc89448cbf14b690d275` after review iteration 1, including
sandboxed Chrome local/hosted qualification. The latter is M30's exact unchanged
Kinet source baseline; UI-only changes add no Authoring import, dependency re-pin,
wire/schema or execution-authority change. Re-read Kinet's final worktree/commit
before qualification, use a temporary exact-source override, and report observed
effective resolution. No M30 row starts through this reconciliation. Preserve
Ramen's old standalone pin and the user-owned workspace; no sibling edits or
publication are authorized.

## Execution resumption — 2026-10-02

W12/U09 acceptance is complete. The user confirmed and resumed the explicit
Stage 6 order with one execution owner, `COMMIT_POLICY: task`, and no Authoring
external mutations. The approved M30 planning drafts are included in this task's
first commit; unrelated work stays untouched. Kinet closure is
`2970731766af13657d464983791c6c17e6312c0d`, with unchanged accepted source
`45217237f62993c278ccdc89448cbf14b690d275`. This repository retains completed
specifications/statuses under its existing convention; the Kinet retirement
envelope is not imposed here. Only M30.1 is selected; review count stays 0/10.

### M30.1 evidence

[Consumer inventory](../../docs/m30-consumer-inventory.md) records full observed
consumer revisions and test-inclusive dependency graphs. Every workspace module
outside the Authoring facade itself and frozen Ramen has no removed-package
import, including unchanged OpenUdon/Kinet and W8M's active check module. Frozen
Ramen standalone build and effective old-pin resolution passed offline with
readonly manifests and installed Go 1.26.7. Initial offline automatic-toolchain
and W8M-root attempts are not acceptance evidence; explicit installed toolchain
and the actual W8M module corrected those checks. No consumer or operator
workspace file was changed. M30.1 is complete; review remains 0/10.

### M30.2 evidence

Both compatibility packages and their tests are removed. Examples now use
`engine`; root API/boundary tests guard neutral records and the absence of both
retired package directories. Current instructions, README, compatibility and
memory-bank descriptions record the approved intentional source break. The
compatibility script binds unchanged Kinet/OpenUdon through temporary modfiles,
including subprocess builds, and keeps Ramen standalone at its frozen pin.
Workspace and standalone `go test ./...`, workspace `go vet ./...`, shell syntax
and `git diff --check` passed offline with installed Go 1.26.7. All 56 surviving
neutral-package source/test/fixture and manifest files are byte-identical to
`dc8f3d61970ae628fc0399b0ef42187aa62a3e5b`; the intentional compatibility-script
change is outside that equivalence scope. Full race and exact-consumer gates
remain M30.3, not yet acceptance evidence. No neutral implementation, dependency,
consumer manifest or operator workspace change.

### M30.3 qualification attempt — environment correction

The first offline exact-source run completed Authoring full workspace/standalone
and race checks, and OpenUdon full test/vet at unchanged
`fbda7e9231b8b306fd1ae3ac623e9d70331b3e08`. Both consumer module resolutions
bound `/home/peter/Workspace/authoring` at source
`db4f5193bc53819be2b1bf1d8728735941e28d9e`. Kinet's full suite failed only
`TestExternalAPIKeepsOriginalRequestTransientAndRejectWritesNoLedger`: its
existing shared-card canonicalizer assumes `/tmp/Test...` and did not normalize
the alternative TMPDIR. This is failed evidence, not accepted qualification;
no fixture was regenerated and no Kinet code changed. The corrected invocation
keeps TMPDIR=/tmp and places only Go build artifacts in disk-backed GOTMPDIR.
Ramen was not reached in that failed script run. Review remains 0/10.

The second attempt with TMPDIR=/tmp still failed that same fixture because
Go 1.26's `testing.TempDir` uses GOTMPDIR directly (installed toolchain
`src/testing/testing.go`, `makeTempDir`). The first correction was insufficient
and is not accepted evidence. Both TMPDIR and GOTMPDIR now use /tmp for the
unchanged full compatibility gate; the inspected filesystem has 2.7 GiB free.
No fixtures, consumer code or manifests were changed. Large explicit native
qualification outputs remain in the private disk-backed cache.

The third script invocation hit the host's per-user /tmp disk quota while
linking OpenUdon's elicitor test binary, despite filesystem-wide free space.
It is failed evidence, not a code or compatibility failure. Existing valid
OpenUdon full test/vet evidence from the first two invocations is reused. The
remaining Kinet server suite is compiled with the exact-source temporary modfile
on disk and run separately with standard /tmp fixture paths; Kinet's other test
packages passed at the same source in both earlier invocations. This separates
large link output from runtime fixture paths without editing consumer tests or
fixtures. Kinet vet and frozen Ramen remain explicit final checks.

The first direct compiled-server invocation used the repository root instead
of the package working directory that `go test` supplies. Relative workflow
fixtures were therefore unavailable; that attempt failed and is not acceptance
evidence. The unchanged binary is rerun from Kinet `internal/server` with the
standard runtime fixture environment. No code, fixture or module change.

### M30.3 completed — exact-source qualification

[Qualification record](../../docs/m30-qualification.md) records the composed
passing gate, effective consumer directory/version, full source revisions,
commands and hashes of both passing and failed evidence. Authoring full
workspace/standalone tests/vet, standalone race and import/API/absence gates
passed. Unchanged OpenUdon full tests/vet passed at the exact new Authoring;
all unchanged Kinet packages passed across its full-suite and corrected complete
compiled-server invocation, with explicit Kinet vet. Frozen Ramen resolved its
original cached pin with no replacement and built standalone. All 56 surviving
neutral-package source/test/fixture and manifest files remain byte-identical to
the pre-removal baseline. Kinet/OpenUdon/Ramen worktrees and operator go.work
are unchanged. Qualified removal source:
`db4f5193bc53819be2b1bf1d8728735941e28d9e`; only status/technical docs were dirty.
M30.3 is complete; review remains 0/10 and acceptance awaits M30.4.

## Closing review iteration 1 — started

Started after reading the persisted 0/10 count and qualification record.
Review baseline is full M30 change from
`dc8f3d61970ae628fc0399b0ef42187aa62a3e5b` through
`9241c481a281ab31f3e4b6ae89f9242ae9517c6c`, plus this status selection.
Qualified implementation source is
`db4f5193bc53819be2b1bf1d8728735941e28d9e`; intervening changes are documentation
only. Review the removed surfaces, unchanged neutral engine/API/wire/safety,
consumer/source resolution, compatibility tooling, frozen Ramen exception,
documentation, permissions and evidence. No acceptance is claimed yet.

## Closing review iteration 1 — passed

Reviewed the full M30 diff and approved planning amendments, including every
removed alias/forwarder and CLI helper, surviving engine implementation and
root examples/API/boundary tests, exact-source compatibility tooling, consumer
inventory and all qualification outcomes. No open P1, P2 or higher finding.
The removed facade has no independent implementation or required consumer
outside frozen Ramen. All 56 surviving neutral source/test/fixture and manifest
files remain byte-identical to baseline; root examples preserve expected output.
Generic type identity, sentinel ownership, JSON tags/versions, prompts,
atomicity, cancellation, persistence/redaction, repair and state behavior stay
unchanged. The deliberate removed-import source break and operator-workspace
consequence are explicit. Temporary modfiles bind new Authoring into unchanged
Kinet/OpenUdon, including child builds; frozen Ramen resolves/builds independently
at its original pin. Failed environment/working-directory attempts are retained;
the composed passing gate covers every required test package and vet/race gate.

No consumer manifest, wire/schema, dependency, runtime pin, owner ledger,
operator workspace, secret, install or publication change was made by M30.
Package-local completed-record conventions are preserved. Remaining consolidation
and exact-source handoff are documentation only under M30.4. Direction remains
the approved v8; no successor evolution version is justified.

## M30.4 completed — accepted retirement and handoff

Accepted removal source: `db4f5193bc53819be2b1bf1d8728735941e28d9e`.
Qualification record: `9241c481a281ab31f3e4b6ae89f9242ae9517c6c`, identical
implementation. Review 1/10 passed; all tasks/checks complete.
[Exact-source handoff](../../docs/m30-handoff.md) returns consumer dispositions
to Kinet and OpenUdon without repinning or sibling writes. Kinet's coordinator
can close Stage 6 after reconciling this record; M17/STG-07 remain excluded.
Current truth and status index are consolidated. The v8 planning snapshots and
completed M29 evidence retain their recorded historical context. Authoring's
existing convention keeps this full accepted status/specification in its own
ledger; no Kinet retirement envelope or ID allocation is introduced. Publication
requires separate authority and was not performed. No required work remains.
````
