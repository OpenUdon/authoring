# M30 consumer inventory — 2026-10-02

Source inspection and test-inclusive import graphs use installed Go 1.26.7,
`GOTOOLCHAIN=local`, `GOPROXY=off`, and no fetch or consumer edit. The parent
workspace was used for its ten registered modules; Kinet and W8M’s active
`browser-workflows/check` module were inspected standalone. Authoring’s own
compatibility packages and examples are expected at this pre-removal baseline.

| Repository/module | Observed revision | Removed-package imports |
| --- | --- | --- |
| apitools | `8a52c3f602988b945a4b5c1960bce8c04170c63d` | None |
| authoring | `dc8f3d61970ae628fc0399b0ef42187aa62a3e5b` | Own compatibility facade/tests (to remove) |
| browsertools | `2cdd788e2f9ed38536fd48d743200e89eff7cf32` | None |
| evidence | `d324099cb15e4a53289e2179fff45b9b88d9274b` | None |
| openudon | `fbda7e9231b8b306fd1ae3ac623e9d70331b3e08` | None |
| ramen | `279a099a418bf324c693c5390ca2a64b967e5669` | `icot`, `icotcli` (frozen consumer) |
| simclaw | `e0cf8b9c68d9901f84e86cc32d3bc5bc4d84b221` | None |
| tfconfig | `61087ab455870b749c1a2c6f4dc90ef3773cfa39` | None |
| udon | `6c4fb8c80a06179f8c3e33ebe685b2d44f66d273` | None |
| uws | `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` | None |
| kinet | `2970731766af13657d464983791c6c17e6312c0d` | None |
| w8m-check | `b23fd63772e0d68f4e5ba1acb59d31b3c757d391` | None |

OpenUdon’s remaining literal strings name forbidden imports in qualification
tests; they are not imported packages. W8M has no root Go module; its active
check module supplies the relevant graph. Archived candidate modules and
historical qualification evidence are neither edited nor adopted here.

## Frozen Ramen and workspace consequence

Ramen remains at Authoring `v0.0.0-20260820042256-2f73e3526583` with no
replacement under `GOWORK=off`. `go list -mod=readonly -m -json` resolved it
to the installed module cache and retained its original sums. A fresh
`go build -mod=readonly -o <private temporary path> ./cmd/ramen` passed
standalone with Go 1.26.7 and disk-backed temporary output.

Removing the facade intentionally ends Ramen’s compatibility with current
workspace Authoring. Freezing Ramen does not drop its dependency. The operator
must separately update or exclude Ramen when using the parent workspace for
Ramen builds; M30 makes no such edit. Other consumers need no migration.

The observed parent `go.work` SHA-256 is `f49366357c1ebb2c9835451b7f88a3c1e5920644cd597c79fe406391fe0ec23c`.
Initial automatic-toolchain selection failed because offline checksum
resolution could not verify the workspace-selected toolchain. Re-running with
the already installed explicit Go 1.26.7 binary passed every graph; no download
or security setting change was used.

## Qualification handoff

M30.3 must bind the exact post-removal Authoring source into unchanged Kinet
and OpenUdon with temporary overrides. Their ordinary old-pin checks alone do
not prove that gate. No publication or consumer re-pin is authorized.
