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
- `operationlifecycle`: prompt-safe lifecycle operation sibling expansion.
- `structured`: provider-neutral structured JSON completion helpers.
- `icot`: progressive loop and bound-runtime interfaces.
- `readiness`: readiness summaries and question planning.
- `decision`: decision evidence and confidence policy.
- `report`: agent result, retention metadata, and scorecard records.
- `promptcontext`: prompt-safe source, operation, schema, and credential
  binding summaries.

## Versioned Records

Durable JSON records use `authoring.*.v1` version constants. The current
versioned records are sessions, transcripts, prompt transcripts, draft
envelopes, prompt contexts, agent results, report metadata, and scorecards.

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
3. Run OpenUdon and Ramen dependent tests from the parent workspace.
4. Update product adapters without changing product JSON versions unless a
   separate product migration explicitly requires it.

The reusable gate is:

```bash
./scripts/check-compat.sh
```
