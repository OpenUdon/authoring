# M30 exact-source qualification — 2026-10-02

## Qualified sources

Authoring removal source: `db4f5193bc53819be2b1bf1d8728735941e28d9e`.
Only status/technical documentation was uncommitted during the checks. All 56
surviving neutral-package source/test/fixture and manifest files are byte-identical
to pre-removal `dc8f3d61970ae628fc0399b0ef42187aa62a3e5b`. Root examples change
only the removed import/type prefixes; expected output is unchanged. No neutral
signature, JSON tag/version, prompt/transcript, dependency or implementation change.

| Consumer | Unchanged revision | Effective Authoring resolution |
| --- | --- | --- |
| OpenUdon | `fbda7e9231b8b306fd1ae3ac623e9d70331b3e08` | Temporary modfile replacement to `/home/peter/Workspace/authoring`, removal source above; original required version remains `v0.0.0-20260930234600-18056cb6b0c1`. |
| Kinet | `2970731766af13657d464983791c6c17e6312c0d` | Temporary modfile replacement to the same exact source; original required version remains `v0.0.0-20260920024745-b417eb681746`. Accepted U09 code is unchanged from `45217237f62993c278ccdc89448cbf14b690d275`. |
| Frozen Ramen | `279a099a418bf324c693c5390ca2a64b967e5669` | No replacement; cached `v0.0.0-20260820042256-2f73e3526583`, sum `h1:NjQ6S2XEk9iiK72rqWxuayd7zBBkIsVed3bVnMkCFkI=`. |

All consumer worktrees remain clean. Their manifests and operator `go.work`
are unchanged. Workspace SHA-256:
`f49366357c1ebb2c9835451b7f88a3c1e5920644cd597c79fe406391fe0ec23c`.
Kinet's OpenUdon runtime pin stays
`c2f161d762bc9f2217bbf0c34b00cdef64b0f7d0`; the OpenUdon source qualification
above is a separate consumer check, not runtime adoption.

## Commands and composed passing gate

Installed Go 1.26.7 was selected explicitly on PATH, with GOTOOLCHAIN=local,
GOPROXY=off and GOSUMDB=off. No install, fetch, live provider or account service.
Tests use fake providers and disposable fixtures.

- Authoring: workspace and GOWORK=off `go test ./...`, `go vet ./...`, and
  standalone `go test -race ./...` passed, including root import/API/absence
  checks and unchanged engine atomicity, cancellation, repair, persistence,
  redaction and no-progress suites.
- OpenUdon: GOWORK=off with
  `GOFLAGS="-modfile=<temporary-openudon.mod> -mod=readonly"`; full
  `go test ./...` and `go vet ./...` passed. The copied modfile replaces only
  Authoring with the exact checkout. GOFLAGS applies to subprocess Go builds.
- Kinet: same temporary-modfile contract. Every package outside `internal/server`
  passed the full `go test ./...` invocation. The unchanged complete server suite
  was compiled with `go test -c -o <private-disk-path>/server.test ./internal/server`
  under the same exact-source modfile, then run from `internal/server` with
  TMPDIR/GOTMPDIR=/tmp and `-test.v -test.timeout=10m`; all server tests passed.
  `go vet ./...` passed with the same source override. This composed gate covers
  all default Kinet test packages at the unchanged source.
- Ramen: GOWORK=off, GOFLAGS=-mod=readonly, cached module resolution and
  `go build -o <private-disk-path>/ramen ./cmd/ramen` passed at the frozen pin.
- Shell syntax, formatting and `git diff --check` passed. Neutral byte comparison
  and clean consumer/workspace checks passed.

## Retained attempts and evidence

Private evidence root: `/home/peter/.cache/kinet-stage6-qualification`.
Failed attempts remain recorded in [M30 status](../tabilet/memory-bank/status-M30.md):
redirected test paths failed the existing /tmp-only card normalizer; changing
TMPDIR alone was insufficient because Go 1.26 `testing.TempDir` reads GOTMPDIR;
the all-/tmp run exceeded per-user link quota despite filesystem free space;
the first direct test-binary run used the wrong package working directory.
No consumer test, golden fixture, shared wire byte or host quota was changed.
Valid same-source checks were reused; the corrected compiled server run and
explicit vet/frozen build close the remaining checks.

| Evidence file | SHA-256 |
| --- | --- |
| `authoring-m30-compat-db4f519.log` (passing Authoring/OpenUdon, failed Kinet card path) | `5fe8a659e9fa1d84e84fcffe11cb45aa494ef5bdb160c2e56395b855060881a8` |
| `authoring-m30-compat-db4f519-standard-temp.log` (insufficient environment correction) | `ff1adbf0c1d918d264badcde22b14f6e3b1de1f6b977d59af149f816b61e1572` |
| `authoring-m30-compat-db4f519-default-temp.log` (link quota failure) | `1da3d5434fbf98be9b949cd2e539e5be04993026bb14447fc887445fe9851d8f` |
| `authoring-m30-final-checks/qualification.log` (exact resolution/compile; wrong runtime working directory) | `42fa815c00d7f0e26f617c235ce30504374072273794790382fadf8979b0b657` |
| `authoring-m30-final-checks/qualification-correct-cwd.log` (complete server suite, Kinet vet, frozen Ramen resolution/build passed) | `8b29b2cce7a65cd54f03b21b16ffe08654eb02298399281397c683b10a0608f9` |

The detailed inventory is [M30 consumers](m30-consumer-inventory.md). Qualification
is local; publication, consumer re-pinning and operator workspace changes remain
separate decisions.
