# Authoring Compatibility Policy

Authoring is pre-1.0. Public packages are intended for OpenUdon and Ramen
adoption, but exported APIs may still change while those adopters migrate.

## Public Surface

The public package groups are behavior-based:

- `trust`: Evidence-backed diagnostic, redaction, artifact, and digest names.
- `session`: durable prompt session state.
- `transcript`: durable transcript turns, events, and model provenance.
- `prompt`: local prompt modes, replay scripts, and prompt transcripts.
- `lifecycle`: draft persistence, atomic writes, and artifact helpers.
- `structured`: provider-neutral structured JSON completion helpers.
- `engine`: progressive loops, atomic interview binding and bound-runtime interfaces.
- `interview`: dependency-aware interview graph and unified evidence records.
- `readiness`: readiness summaries and question planning.
- `decision`: decision evidence and confidence policy.
- `report`: agent result, retention metadata, and scorecard records.
- `promptcontext`: prompt-safe source, operation, schema, and credential
  binding summaries.

## Versioned Records

Durable JSON records use `authoring.*.v1` version constants. The current
versioned records are interviews, sessions, transcripts, prompt transcripts,
draft envelopes, prompt contexts, agent results, report metadata, and
scorecards.

The frontier migration intentionally replaces the readiness question wire
aliases (`force_ask`, `allow_default`, `default_answer`, `suggested_answer`,
`default_source`, and `grouped`) with `forced`, `recommendation`, `priority`,
`rationale`, and `evidence_refs`. Deprecated Go-only fields remain temporarily
for source compatibility with the unchanged Ramen adapter, but they are not
serialized. Product v1 session/report decoders must not be reused for a v2
product migration.

Within pre-1.0, version constants and JSON tags should not change silently.
If a JSON shape changes incompatibly, update the version constant, document the
migration, and run downstream OpenUdon and Ramen checks.

## Boundary Rules

Authoring must not import OpenUdon, Ramen, UWS, or API-source packages.
Downstream adapters translate product data into Authoring contracts and own
product-specific prompts, schemas, validation, artifact formats, execution,
credentials, governance, state, and reconciliation.

Default Authoring tests must remain provider-free, model-free, executor-free,
credential-free, and live-network-free.

## Migration Expectations

OpenUdon and Ramen adopters should wrap Authoring APIs behind product-owned
adapter code where product wire compatibility matters. During pre-1.0 changes,
the expected migration path is:

1. Update Authoring.
2. Run Authoring tests, vet, and import-boundary checks.
3. Qualify unchanged OpenUdon and Kinet with temporary exact-source overrides;
   separately qualify frozen Ramen standalone at its old pin.
4. Update product adapters without changing product JSON versions unless a
   separate product migration explicitly requires it.

Readiness compatibility aliases are write-through normalized: all aliases
reflect the durable `forced` and `recommendation` values after normalization,
and `DefaultSource` is always `"recommendation"`. `engine.Result.Answers` holds
the most recently collected round: it is partial when returning `needs_input`
and otherwise contains the last successfully applied round.

The move of API lifecycle ranking from
`github.com/OpenUdon/authoring/operationlifecycle` to
`github.com/OpenUdon/apitools/operationlifecycle` is an intentional pre-1.0
source-breaking boundary correction. Consumers pass `apitools.OperationSummary`
records and retain API-source provenance; no Authoring durable JSON version
changed for this move.

The reusable gate is:

```bash
./scripts/check-compat.sh
```

## Neutral engine extraction and facade retirement

M29 introduced `github.com/OpenUdon/authoring/engine` and compatible `icot`
aliases/forwarders. M30 removes `authoring/icot` and `authoring/icotcli` as an
intentional pre-1.0 source break. This is the approved compatibility exception:
frozen Ramen retains Authoring `v0.0.0-20260820042256-2f73e3526583` and its old
imports, with no source/manifest change. It has not dropped the dependency;
compatibility with current workspace Authoring intentionally ends.

Other public APIs, generic signatures, methods, constants, sentinel identity,
JSON shapes/versions, prompt/error wording and state transitions are unchanged.
`engine` remains the one neutral implementation. Reflection reports `engine`
ownership as established in M29; durable storage uses explicit versioned JSON,
not Go package identity. Error sentinels must not be reassigned by consumers.

The explicit compatibility script checks Authoring standalone/workspace and
binds the exact current source into unchanged OpenUdon and Kinet with disposable
modfiles, including subprocess Go builds. Ramen is checked only standalone at
its original cached pin. The operator-owned parent `go.work` remains unchanged;
Ramen cannot build there after removal until the operator separately excludes or
updates it. [M30 inventory](docs/m30-consumer-inventory.md) records observed
consumers. No publication or downstream re-pin follows from local acceptance.
