# Status M30 — Retire the iCoT compatibility facade

State: M30.1–M30.3 complete; closing review and handoff pending.

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
| M30.4 Review and downstream handoff | `[ ]` | Persist a bounded deep review (maximum ten iterations), record exact source and checks, and hand Kinet and OpenUdon the exact revision with no required consumer change. Publication needs a separate request. |

## Verification and acceptance

See the full specification in milestone.md. Default checks use fake providers
and disposable fixtures; no live model, account or network service. `engine`
keeps no downstream import. Prompt and transcript bytes, JSON tags, durable
versions, signatures of non-removed packages and state behavior are unchanged.

## Closing review

Persisted iteration count: 0/10. Task completion alone does not establish
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
