# Status M31 — Ramen neutral-engine consumer adoption

State: Pending; all four rows pending. Whole-milestone review 0/10, not started.

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
| M31.1 Import inventory and neutral-engine mapping | `[ ]` | Inventory Ramen's `github.com/OpenUdon/authoring` imports at its current pin and map each retired `icot`/`icotcli` use onto its neutral `engine`/`interview`/`session`/`transcript`/`report` counterpart, recording an explicit disposition (re-implement, replace, or drop) for any behavior without an equivalent. |
| M31.2 Ramen migration to the adopted revision | `[ ]` | Adopt the published revision, update `go.mod`/`go.sum`, and remove every retired import while preserving the operator-visible `ramen author` and `ramen icot` contracts and their gates. |
| M31.3 Cross-repo qualification and adoption record | `[ ]` | Qualify Authoring and the migrated Ramen against the exact adopted revision and record the adoption in both ledgers. |
| M31.4 Bounded deep review and downstream handoff | `[ ]` | Review the whole milestone and hand off the satisfied `Authoring:M31 -> Ramen:M85` edge. |

## Review

Persisted iteration count: 0/10. Task completion alone does not establish
acceptance.
