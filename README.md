# Authoring

Authoring is the shared Go module for product-neutral authoring orchestration
used by OpenUdon and Ramen.

It owns generic sessions, transcripts, prompt/replay helpers, draft lifecycle
persistence, structured JSON fallback, readiness/question planning, progressive
iCoT loops, decision evidence, report metadata, scorecard records, and
prompt-safe context shapes.

It does not own OpenUdon workflow package semantics, Ramen desired-state
semantics, UWS document semantics, API-source parsing, credential resolution,
model-provider clients, live execution, governance, state, or reconciliation.

## Packages

- `trust`: Evidence-backed diagnostic, artifact, digest, and redaction names.
- `session`: durable prompt session state and answers.
- `transcript`: transcript turns, events, diagnostics, artifacts, and model
  provenance.
- `prompt`: local prompting, default modes, replay scripts, and prompt
  transcripts.
- `lifecycle`: draft envelopes, atomic writes, autosave, and artifact helpers.
- `structured`: provider-neutral structured JSON completion and legacy JSON
  fallback.
- `icot`: generic progressive loops and bound-runtime interfaces.
- `readiness`: readiness summaries, blocking/warning sorting, and question
  planning.
- `decision`: decision evidence, confidence behavior, and confirmation policy.
- `report`: agent result contracts, report metadata, scorecards, and variants.
- `promptcontext`: prompt-safe source documents, operations, schema hints, and
  symbolic credential binding names.

## Boundary

Authoring is intentionally upstream of products. It imports Evidence for
neutral trust primitives and must not import OpenUdon, Ramen, UWS, or API-source
packages. Downstream adapters translate product metadata into Authoring
contracts and own product prompts, validation, artifacts, credentials, model
clients, execution, and state.

Default tests and examples use fake runtimes and fake clients only. They do not
require credentials, model providers, API calls, workflow execution,
Terraform/OpenTofu execution, or trusted-runner access.

## Compatibility

Authoring is pre-1.0. Versioned durable records use `authoring.*.v1` constants,
but exported APIs may still change while OpenUdon and Ramen adoption settles.
See [COMPATIBILITY.md](COMPATIBILITY.md) for the current policy and migration
expectations.

## Examples

Runnable examples in `examples_test.go` cover manual prompting, agent-style
result/digest records, structured JSON fallback with a fake client, and a
progressive loop with a fake runtime.

## Checks

```bash
go test ./...
go vet ./...
git diff --check
```

From the Authoring checkout inside the parent workspace, run the cross-repo
gate:

```bash
./scripts/check-compat.sh
```
