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
- `icot`: generic progressive and interactive iCoT loops, lifecycle hooks, and
  bound-runtime interfaces.
- `icotcli`: shared iCoT CLI flag plumbing and prompt/model label helpers.
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

## Downstream Use Cases

Authoring is shared loop infrastructure, not a product CLI. OpenUdon and Ramen
use it for different authoring jobs:

- OpenUdon iCoT starts from workflow automation intent. Its downstream adapter
  may guide catalog/API artifact selection or retrieval, then writes OpenUdon
  workflow authoring artifacts such as `project.md`, `workflows/intent.hcl`,
  review/eval reports, and package-oriented files. This path is oriented
  toward general workflow authors who may need more guidance from a goal to API
  artifact selection.
- Ramen iCoT starts from desired-state infrastructure intent plus explicit
  local API source metadata, such as OpenAPI, Smithy, or Discovery documents.
  Its downstream adapter writes native Ramen/UWS desired-state projects such as
  `project.uws.yaml`, then can run Ramen validation, graph, and static plan
  gates before any trusted executor boundary. This path is oriented toward
  enterprise IaC users who are expected to know the target API source category
  and want controlled validation, planning, state, approval, and apply
  semantics.
- Ramen conversion is separate from iCoT. It converts existing
  Terraform/OpenTofu configuration plus explicit local API source metadata into
  native Ramen/UWS review artifacts.

The practical distinction is that workflow authoring can be exploratory, while
IaC authoring should make the target API source category explicit. Authoring
supports the shared prompting, transcript, readiness, decision, and report
mechanics for both, but source discovery, provider mappings, generated artifact
formats, state, approval, and execution remain downstream responsibilities.
Ramen should not import OpenUdon workflow packages or expose a generic
`from-openudon`/`import openudon` path; the use cases and product semantics are
intentionally different.

### Shared Context

The two iCoT adapters share generic Authoring context and loop mechanics:

- prompt sessions, default modes, answer replay, and transcript records
- progressive loop lifecycle, autosave hooks, and final confirmation flow
- readiness issue and interactive question shapes
- decision evidence normalization and confirmation policy
- prompt-safe source, operation, schema, and credential-binding context records
- agent result, diagnostic, artifact, metadata, and scorecard report shapes
- common CLI flags for prompt mode, no-LLM, model labels, answers, reports, and
  transcripts

The downstream adapters provide the product hooks that make those records mean
something. OpenUdon binds workflow intent, catalog retrieval, request mapping,
flow review, and package artifacts. Ramen binds API-source translation,
desired-state resources, validation, graph, static plan, state, approval, and
trusted execution boundaries.

### OpenUdon iCoT

OpenUdon exposes its workflow-authoring loop from the downstream checkout:

```bash
(cd ../openudon && go run ./cmd/icot --example examples/<name>)
(cd ../openudon && go run ./cmd/openudon build --example examples/<name>)
```

The iCoT command writes the human/project brief and structured intent. The
downstream `openudon build` command then deterministically regenerates the
public workflow artifacts, review evidence, and quality reports from
`workflows/intent.hcl`.

Common downstream options:

| Option | Purpose |
|---|---|
| `--example DIR` / `--dir DIR` | Example directory to create or update. |
| `--from-example DIR` | Seed answers from an existing example. |
| `--answers PATH` | Replay a YAML or JSON session/answers file. |
| `--force` / `--yes` | Overwrite existing files, optionally without prompts. |
| `--print` | Render without writing files. |
| `--agent` | Return `needs_input` instead of prompting when incomplete. |
| `--no-llm` | Disable optional model extraction assistance. |
| `--provider`, `--model`, `--temperature` | Downstream model configuration. |
| `--prompt-mode full\|normal\|fast` | Control defaulted question behavior. |
| `--json`, `--report PATH` | Emit or save a structured report. |

Typical downstream output files and directories:

```text
project.md
workflows/intent.hcl
workflows/workflow.hcl
workflows/workflow.uws.yaml
openapi/
workflows/
expected/
expected/plan.json
expected/review.md
expected/review-handoff.json
expected/quality.json
.icot/
```

OpenUdon also owns related downstream subcommands such as `build`, `reconcile`,
`lint`, `repair`, `scorecard`, `variants`, `replay-eval`, `authoring-eval`, and
`report verify`.

### Ramen iCoT

Ramen exposes a desired-state authoring path from the downstream checkout:

```bash
(cd ../ramen && go run ./cmd/ramen icot \
  --goal "Create an Azure Cosmos DB account" \
  --api-source openapi:azure-cosmos=/abs/path/azure-cosmos.json \
  --out .ramen/icot/azure-cosmos \
  --validate \
  --graph \
  --plan)
```

Common downstream options:

| Option | Purpose |
|---|---|
| `--goal TEXT` | Desired-state project goal. |
| `--api-source KIND:ID=PATH` | Repeatable local API source input. |
| `--project-name NAME` | Optional generated project name. |
| `--out DIR` | Output directory for `project.uws.yaml`. |
| `--validate`, `--graph`, `--plan` | Run non-executing Ramen gates after drafting. |
| `--state PATH` | SQLite state path for the static plan gate. |
| `--answers PATH` | Replay answers/session JSON. |
| `--agent` | Run noninteractively and return `needs_input` when incomplete. |
| `--no-llm` | Disable optional model assistance. |
| `--provider`, `--model`, `--temperature` | Downstream model configuration. |
| `--prompt-mode full\|normal\|fast` | Control defaulted question behavior. |
| `--json`, `--report PATH` | Emit or save a structured report. |
| `--no-transcript` | Disable transcript persistence. |

Typical downstream output files:

```text
.ramen/icot/project.uws.yaml
.ramen/icot/<optional transcript/report files>
```

When requested, validation, graph, and plan results are returned in the command
result. They do not execute API operations or trusted executor actions.

### Ramen Convert

Ramen conversion is not an iCoT loop. It converts existing Terraform/OpenTofu
configuration plus explicit local API source metadata:

```bash
(cd ../ramen && go run ./cmd/ramen convert \
  --config-dir ./tf \
  --api-source openapi:azure-cosmos=/abs/path/azure-cosmos.json \
  --action create \
  --out .ramen/convert/azure-cosmos)
```

Common downstream options:

| Option | Purpose |
|---|---|
| `--config-dir DIR` | Terraform/OpenTofu configuration directory. |
| `--api-source KIND:ID=PATH` | Repeatable local API source input. |
| `--openapi ID=PATH` | OpenAPI shorthand for `--api-source openapi:ID=PATH`. |
| `--action create\|update\|delete\|replace` | Required for managed resources. |
| `--target ADDRESS` | Repeatable Terraform address filter. |
| `--out DIR` | Output directory for review artifacts. |
| `--strict` | Fail when strict-failure diagnostics remain. |

Typical downstream output files:

```text
.ramen/convert/project.md
.ramen/convert/project.uws.yaml
.ramen/convert/project.uws.hcl
.ramen/convert/workflows/workflow.uws.yaml
.ramen/convert/workflows/workflow.uws.hcl
.ramen/convert/expected/conversion.json
.ramen/convert/expected/mappings.json
.ramen/convert/expected/diagnostics.json
.ramen/convert/expected/diagnostics.md
.ramen/convert/expected/plan.json
.ramen/convert/expected/plan.md
.ramen/convert/expected/review.md
```

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
