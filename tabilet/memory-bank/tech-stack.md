# Tech Stack

Authoring is a Go module that provides shared authoring orchestration
primitives.

## Language And Runtime

- Primary language: Go.
- Module path: `github.com/OpenUdon/authoring`.
- Go directive: `go 1.26.3`.
- Default tests must run without model providers, credentials, network access,
  API execution, workflow execution, or sibling private modules.

## Planned Dependencies

- `github.com/OpenUdon/evidence/diagnostic` for neutral diagnostic records.
- `github.com/OpenUdon/evidence/redact` for shared redaction helpers before
  durable transcript, report, and artifact metadata persistence.
- `github.com/OpenUdon/evidence/artifact` for safe artifact path records and
  generated artifact summaries.
- `github.com/OpenUdon/evidence/digest` for digest records and sidecars.
- Standard library packages first for core loop mechanics.

Authoring should avoid direct dependencies on `openudon`, `ramen`, `uws`, or
`apitools` unless a generic interface proves insufficient. Downstream adapters
should normally supply UWS/API-source/product behavior.

M02 adds Authoring's `trust` package as the first shared foundation. `trust`
aliases Evidence record types and delegates helper behavior to Evidence
packages, so later Authoring packages can use stable Authoring names without
creating a second durable diagnostic, redaction, artifact, or digest contract.
M03 adds `session` and `transcript` packages with version constants
`authoring.session.v1` and `authoring.transcript.v1`, deterministic
normalization, and canonical JSON helpers.
M04 adds the `prompt` package and prompt transcript envelope version
`authoring.prompt-transcript.v1`.
M05 adds the `lifecycle` package and draft envelope version
`authoring.draft.v1`.
M06 adds the `structured` package with provider-neutral client interfaces and
no model-provider dependencies.
M07 adds the `icot` package for provider-free progressive loop orchestration.
M08 adds `icot.Runtime`, optional runtime hook interfaces, and `RunRuntime` for
bound downstream adapters.
M09 adds the `decision` package for generic decision evidence and confidence
policy.
M10 adds the `readiness` package for provider-free readiness summaries,
decision-confirmation issue projection, generic question plans, and safe
default handling. `icot` imports `readiness` only for the public question alias
and default readiness policy.
M11 adds bounded repair orchestration to `icot`, including review/remediation
containers, repair result statuses, transcript events, and runtime binding
helpers. It remains provider-free and depends only on M10 readiness records and
M03 transcript events.
M12 adds the `report` package with `authoring.agent-result.v1` result JSON,
generic agent statuses, readiness/repair/decision summaries, diagnostics,
artifact descriptors, digest records, and M03 session/transcript metadata.
M13 extends `report` with `authoring.report-metadata.v1`, generic retention
classes, provider-output/archive/redaction flags, and digest sidecars using
M02 `trust` digest and artifact-path helpers.
M14 extends `report` with `authoring.scorecard.v1`, fixture/variant result
records, grouped summary counters, failure-family summaries, validation
diagnostics, and deterministic scorecard JSON.
M15 introduced `promptcontext`; M27 publishes `authoring.prompt-context.v2`
with prompt-safe source documents, operation candidates, schema hints, ordered
OR-of-AND symbolic credential binding sets, redaction-backed normalization,
canonical JSON, and import-boundary tests.
M16 verifies OpenUdon phase-1 adoption against the parent `go.work` sibling
checkout: OpenUdon delegates generic prompt, structured JSON, atomic write, and
prompt-safe context adapter mechanics to Authoring while keeping OpenUdon wire
versions and product-specific iCoT behavior downstream. Public-module
compatibility for OpenUdon remains an M20 gate until Authoring has a consumable
version.
M17 adds OpenUdon phase-2 adoption of `report` and `readiness` as validation
contracts for agent results, report retention metadata, and scorecard variant
summaries. It does not replace OpenUdon's `openudon.icot-*` JSON schemas.
M18 adds a downstream Ramen adapter spike in `../ramen/authoring` using
Authoring runtime, prompt context, lifecycle, readiness, and report packages.
Ramen owns the native project and validation dependencies, so Authoring's
public module still imports no Ramen or UWS packages. Ramen standalone module
compatibility remains an M20 gate until Authoring has a consumable module
version.
M19 expands that downstream adapter with optional graph and plan gates using
Ramen-owned packages. These imports remain in `../ramen`; Authoring's default
tests and standalone module checks stay provider-free, model-free, and
executor-free.
M20 adds a root import-boundary test and `scripts/check-compat.sh` for
Authoring standalone checks plus parent-workspace OpenUdon/Ramen checks.
M21 adds API-surface tests for durable version constants and JSON tags, plus
`COMPATIBILITY.md` for pre-1.0 migration policy.
M22 adds `README.md`, expanded root package docs, and runnable examples backed
only by fake prompts, fake structured clients, fake runtimes, and local digest
helpers.
M23 adds internal-only helper packages for normalization and shared record
normalization. These packages should depend only on the standard library and
Authoring's `trust` aliases as needed, must not become exported API, and must
preserve default provider-free/model-free/executor-free checks. Downstream
OpenUdon/Ramen changes are limited to compatibility verification and module
version bumps after Authoring is committed.
M24 adds exported interactive iCoT APIs in `icot` and shared flag helpers in
`icotcli`. `icotcli` depends only on the standard library plus Authoring's
prompt default-mode type; provider clients and product command behavior remain
downstream.
M25 adds the standard-library-only `interview` package and
`authoring.interview.v1`. `icot` uses one frontier-round engine for generic,
bound-runtime, and interactive flows; normal termination diagnoses three
consecutive rounds with no state/readiness progress, with I01 adding only an
emergency fuse. Ramen's existing one-question runtime is adapted
source-compatibly without changing Ramen.
I01 adds `interview.Resolution`, `interview.ApplyRound`, generic
`icot.InterviewBinding`, semantic progress fingerprints, monotonic event IDs,
configurable prompt/frontier text, and a 1,000-round default emergency fuse.
It adds no dependency and retains the v1 session, transcript, and interview
wire versions.
D01 adds persistence validation over `session.State`, synced atomic file
replacement, bounded legacy JSON extraction coverage, and an internal shared
cancellation sentinel. It adds no dependency and performs no live I/O beyond
caller-requested local persistence.
M28 places lifecycle-operation ranking in the sibling apitools module over
`apitools.OperationSummary`; Authoring remains independent of apitools while
OpenUdon and Ramen consume both public modules in their adapters.

## Commands

```bash
go test ./...
go vet ./...
GOWORK=off go test ./...
GOWORK=off go vet ./...
git diff --check
git -C ../tofu diff --check -- authoring
../skills/harness/tackle-memory-bank-api-loop --model lane-audit .
```

After exported API changes:

```bash
(cd ../openudon && go test ./...)
(cd ../ramen && GOWORK=off go build -mod=readonly ./cmd/ramen)
(cd ../apitools && go test ./...)
```

OpenUdon's iCoT command/harness is retired. Use its current offline Go and
neutral authoring regression gates; no removed command is part of this check.

The reusable compatibility gate is:

```bash
./scripts/check-compat.sh
```

## Artifact And Schema Expectations

- Session, transcript, and result structures should be versioned before they
  become durable cross-repo contracts.
- JSON output should be deterministic where artifacts are committed or used in
  digest-bound evidence.
- Model-provider-specific request/response DTOs should stay downstream unless
  multiple consumers need the same public contract.
- Prompt text belongs downstream. Shared packages may define prompt-safe
  transcript structures but should not own product prompts.
- Public result contracts should cover noninteractive `complete`,
  `needs_input`, `failed`, and `canceled` states.
- Report metadata should include run ID, command, commit, generation time,
  retention class, provider-output flag, archive safety, redaction-required
  flag, generated artifact descriptors, and digest sidecars.
- Structured-output helpers should expose provider-neutral structured-client
  interfaces, schema normalization, legacy JSON fallback, JSON block
  extraction, decode helpers, and fake-client tests.
- Prompt/replay helpers should keep ask/show/silent modes, required prompts,
  yes/no prompts, transcript save/load, and prompt-label assertions
  deterministic.

## Dependency Rules

- Keep Authoring independent of OpenUdon and Ramen.
- Keep live model clients behind downstream-supplied interfaces.
- Prefer interfaces for runtime behavior and plain structs for durable state.
- Use Evidence for generic trust primitives instead of reimplementing digest,
  artifact, redaction, or diagnostic helpers.
- Do not import private executor modules or Terraform/OpenTofu internals.
- Do not import OpenUdon workflow intent packages, OpenUdon iCoT wire schemas,
  Ramen desired-state packages, UWS schema/model packages, or `apitools`
  catalog packages into the generic module unless a later roadmap change
  explicitly narrows that exception.
- Keep default tests provider-free, model-free, executor-free, and
  credential-free. Use fake runtimes and fake structured clients.

## Neutral engine after M30 facade removal

The surviving loop import is `github.com/OpenUdon/authoring/engine`.
`icot`/`icotcli` and their facade-only tests are removed. No dependency, Go
requirement, neutral wire/version or implementation body changed. The explicit
`scripts/check-compat.sh` checks Authoring workspace/standalone tests, vet, race
and boundary, then uses temporary modfile overrides for unchanged Kinet and
OpenUdon, with effective source resolution printed. `GOFLAGS` carries that
override into subprocess builds. Frozen Ramen is built standalone at its old pin
with readonly module resolution; current workspace Ramen is intentionally excluded.
Run offline with an installed compatible toolchain and disk-backed temporary
output when the host's `/tmp` quota is insufficient. No parent workspace edit or
consumer manifest change is performed. Earlier numbered notes describe the
historical package names at their recorded stage.
