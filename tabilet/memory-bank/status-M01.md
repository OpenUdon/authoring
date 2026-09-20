# Status M01

## Goal

Establish the Authoring public module harness and boundary record.

## Tasks

| Item | State | Notes |
|---|---|---|
| Create tracked harness snapshot | `[+]` | Add AGENTS, memory-bank docs, milestone index, and first status file under `../tofu/authoring`. |
| Wire checkout harness paths | `[+]` | Use the same symlink-facing `AGENTS.md`, `memory-bank/`, and `evolution/` pattern as Ramen. |
| Record execution model | `[+]` | AGENTS and architecture docs describe shared upstream orchestration with downstream bound runtimes. |
| Add module scaffold | `[+]` | Add `go.mod` for `github.com/OpenUdon/authoring`. |

## Notes

- Initial design input is preserved in `../authoring/x.md`.
- Authoring starts as a documentation/module harness. Implementation begins in
  M02 with session and transcript contracts.

## Verification

Planned checks:

```bash
go test ./...
go vet ./...
git diff --check
git -C ../tofu diff --check -- authoring
```
