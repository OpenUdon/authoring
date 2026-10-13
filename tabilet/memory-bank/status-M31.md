# Status M31 — Ramen neutral-engine consumer adoption

State: In progress; M31.1-M31.3 complete, M31.4 pending. Whole-milestone review 0/10, not started.

## Goal and dependencies

Move Ramen off the frozen legacy Authoring pin and onto the published neutral
engine, so Ramen imports no retired `authoring/icot` or `authoring/icotcli`
surface and both repositories build and pass their gates against one shared
revision. Baseline: Authoring
`cb21c3e36d19290dc46a6d2bcbde71b45ea8ad30` (M29 and M30 complete), clean
worktree, which is also `origin/main`, so the adopted revision is already
published. Ramen owns its own `go.mod`, imports, and CLI; Authoring owns this
ledger. If the migration needs a neutral-engine API that does not exist at the
adopted revision, stop and report: adding it is a new Authoring change plus
publication, outside this milestone.

## Provenance

Ramen's approved independent track (2026-10-12) carries `Authoring:M31?` as a
conditional prerequisite of `Ramen:M85`. The user resolved that edge on
2026-10-13 by keeping the unit Authoring-owned and granting sibling write scope
for `../authoring`. M30's planning exception froze Ramen at
`v0.0.0-20260820042256-2f73e3526583` and recorded that Ramen, not Authoring,
owns the migration; `COMPATIBILITY.md` records the same position. This
milestone promotes that handoff. The full specification is in
[milestone.md](milestone.md) (M31).

## Launch policies in force

Serial single coordinator; `COMMIT_POLICY: task`; `SUGGESTED_UPDATE: confirm`;
`EXTERNAL_MUTATIONS: none`; `AUTHORIZATION_GRANTS: {}`; local commits only, no
push and no publication.

## Tasks

Markers: `` `[ ]` `` pending, `` `[~]` `` in progress, `` `[+]` `` complete,
`` `[!]` `` blocked, `` `[X]` `` cancelled, and `` `[-]` `` closed historical
(must name its accepted successor). At most one general row may be in progress.

| Item | State | Notes |
|---|---|---|
| M31.1 Import inventory and neutral-engine mapping | `[+]` | R1. |
| M31.2 Ramen migration to the adopted revision | `[+]` | Ramen adopted `v0.0.0-20261002192730-cb21c3e36d19`: `go.mod` changed only that require line and `go.sum` gained that revision's two hashes (the revision declares no requirements). Seven sites moved from the retired `authoring/icot` facade to `authoring/engine`, which exports every symbol they used; the retired `icotcli` helpers were re-implemented in Ramen at `internal/icotcli` with the same `Flags`/`AddFlags`/`PromptDefaultMode` contract, following OpenUdon's `internal/authoringcli` precedent. No retired import remains and no operator-visible behavior changes (Ramen commit `ca9e86a`). |
| M31.3 Cross-repo qualification and adoption record | `[+]` | Authoring's own `go test ./...` and import-boundary check pass at the adopted revision; `scripts/check-compat.sh` now binds Ramen as a migrated consumer and keeps a standalone build at its own pin; `COMPATIBILITY.md` records the end of the frozen-consumer exception; Ramen's `tech-stack.md` records the adopted pin. |
| M31.4 Bounded deep review and downstream handoff | `[ ]` | Review the whole milestone and hand off the satisfied `Authoring:M31 -> Ramen:M85` edge. |

## M31.1 record

**Inventory.** Nine retired-import sites in Ramen, all at the frozen pin
`v0.0.0-20260820042256-2f73e3526583`:
`authoring/authoring.go` (`icot`), `internal/icot/planner.go` and
`internal/icot/run.go` (`icot`), `cmd/ramen/command_author.go` (`icot` and
`icotcli`), plus four test files (`internal/icot/planner_test.go`,
`internal/icot/run_test.go`, `cmd/ramen/command_author_test.go`,
`tests/icot_replay_test.go`). Symbols used: `ErrCanceled`, `ErrNeedsInput`,
`InterviewBinding`, `Options`, `Question`, `Result`, `RoundAnswer`, `Run`,
`RunRuntime`, `RuntimeConfig` from `icot`; `AddFlags`, `Flags`,
`PromptDefaultMode` from `icotcli`.

**Mapping.** `authoring/icot` maps one-to-one onto `authoring/engine`, the
neutral package M29 introduced and M30 kept; OpenUdon already migrated exactly
this way (`../openudon/internal/authoring/*` imports
`sharedengine "github.com/OpenUdon/authoring/engine"`).

**Material finding.** `authoring/icotcli` has no neutral replacement: M30
removed the shared CLI flag helpers outright, and no `AddFlags`, `Flags`, or
`PromptDefaultMode` exists anywhere in Authoring today. OpenUdon resolved this
by owning its own CLI layer (`openudon/internal/authoringcli`, with its own
`promptDefaultMode` and `PromptDefaultMode`). M31.2 therefore re-implements
those three helpers in a Ramen-owned package instead of importing a shared one.
No operator behavior is dropped by this disposition.

**Adoption probe.** `go get github.com/OpenUdon/authoring@cb21c3e` resolves to
`v0.0.0-20261002192730-cb21c3e36d19` over the public proxy, so the adopted
revision needs no publication. That revision requires Go ≥ 1.26.3 and the local
toolchain is go1.26.6, so no toolchain download is required.

**M31.2 target files.** `go.mod`/`go.sum`; the nine sites above; and one new
Ramen-owned CLI helper package replacing `icotcli`. The `--validate`, `--graph`,
`--plan`, `--network`, `--agent`, prompt-context v2, session, transcript, and
report contracts stay as they are.

| M31.2 Ramen migration to the adopted revision | `[ ]` | Adopt the published revision, update `go.mod`/`go.sum`, and remove every retired import while preserving the operator-visible `ramen author` and `ramen icot` contracts and their gates. |
| M31.3 Cross-repo qualification and adoption record | `[ ]` | Qualify Authoring and the migrated Ramen against the exact adopted revision and record the adoption in both ledgers. |
| M31.4 Bounded deep review and downstream handoff | `[ ]` | Review the whole milestone and hand off the satisfied `Authoring:M31 -> Ramen:M85` edge. |

## Review

Persisted iteration count: 0/10. Task completion alone does not establish
acceptance.
