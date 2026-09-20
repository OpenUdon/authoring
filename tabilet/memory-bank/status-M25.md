# Status M25 - Adaptive Evidence-Grounded Frontier Authoring

State: Complete

## Scope

Add a product-neutral interview graph and replace Authoring's fixed
question-at-a-time loops with one dependency-ready frontier-round engine.
Authoring stores only concise public rationale and evidence; product prompts,
source discovery, artifact schemas, and workflow behavior remain downstream.

## Tasks

| Item | State | Notes |
|---|---|---|
| M25 adaptive frontier authoring | `[+]` | Added `authoring.interview.v1`, graph/evidence/answer/deferral validation and transitions, deterministic frontiers and atomic answers, the cleaned durable question shape, one unlimited frontier engine across generic/runtime/interactive entry points, three-round no-progress diagnosis, more-than-20/mode/graph regression coverage, docs/evolution v3, required verification, and unchanged Ramen compatibility. |
| M25 review remediation | `[+]` | Streams generic loop events into interactive transcripts at emission time (`a3ec037`), adds normalized durable evidence attributes for downstream resume-time safety qualifiers (`18ba806`), and stops frontier answer collection immediately when an operator cancels. |

## Completion Notes

- Durable `authoring.session.v1` and `authoring.transcript.v1` records are
  unchanged because their JSON shapes did not change.
- Deprecated Go-only readiness question aliases keep the unchanged Ramen
  adapter source-compatible, but the aliases do not serialize and are not part
  of the v2 downstream wire contract.
- OpenUdon's v2 adapter and reports are migrated separately in OpenUdon M70.

## Verification

- `go test ./...`
- `go vet ./...`
- `GOWORK=off go test ./...`
- `GOWORK=off go vet ./...`
- `git diff --check`
- Authoring import-boundary test
- `(cd ../ramen && go test ./...)`
- `(cd ../ramen && go vet ./...)`
- `git -C ../tofu diff --check -- authoring`
