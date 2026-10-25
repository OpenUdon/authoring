# Retired milestone M25 - Adaptive Evidence-Grounded Frontier Authoring

**Milestone.** M25
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M25.md
**Source specification.** tabilet/memory-bank/milestone.md#m25-adaptive-evidence-grounded-frontier-authoring
**Evidence.** be023f06d0837b58a651449b6d89a00c12eb04cd
**Worktree.** includes uncommitted changes
**Review.** legacy
**Review iterations.** not recorded
**Verification.** Historical verification and its limitations are retained verbatim in the status record; original commands are not re-run or relabelled by this retirement. At the observed evidence revision, go test ./... and go vet ./... passed on 2026-10-25; literal source preservation and retirement-envelope validation passed before removal.
**Consolidated into.** Current product, architecture, stack, and lessons remain in tabilet/memory-bank/; no implementation or current product contract changes. Superseded roadmap and retention wording is preserved in knowledge.md.
**Specification encoding.** The added level-two identity heading is a retirement-parser wrapper; the following level-three section is the complete original source, retained verbatim.
**Legacy closure.** This milestone closed before adoption of the persisted review gate. Its original closure evidence is the complete literal status record; no review iteration or retrospective implementation review is invented.

## Milestone specification

````markdown
## M25 Adaptive Evidence-Grounded Frontier Authoring

### M25 Adaptive Evidence-Grounded Frontier Authoring

**Goal.** Replace fixed, one-question-at-a-time iCoT loops with one
dependency-aware frontier-round engine and a generic interview state contract.

Acceptance:

- `interview` exposes versioned state, node, evidence, answer, and deferral
  records with open, settled, deferred, and inapplicable node statuses.
- Validation diagnoses duplicate IDs, missing dependencies and evidence,
  cycles, invalid statuses/transitions, and incomplete or invalid deferrals.
- The deterministic frontier includes every open node whose prerequisites are
  settled, ordered by priority and ID; a complete answer set can be applied
  atomically.
- `readiness.Question` has one durable forced flag and recommendation plus
  priority, concise public rationale, and evidence references. Compatibility
  aliases are Go-only and absent from the JSON wire.
- Generic, bound-runtime, and interactive entry points use one active engine
  that displays the entire numbered round before collection and performs one
  normalization/autosave after `ApplyRound`.
- The loop has no frontier breadth ceiling and stops on completion,
  cancellation, downstream-approved finalization, three consecutive
  no-progress rounds, or I01's configurable emergency round fuse.
- Tests cover more than 20 decisions, independent and dependent frontiers,
  deterministic ordering, atomic autosave, full/normal/fast visibility,
  graph failures, invalid deferrals/transitions, cancellation, and repeated
  no-progress.
- Authoring full/standalone/vet/diff/import-boundary gates and unchanged Ramen
  compatibility pass. OpenUdon consumes the frontier API in its M70 v2 wire
  migration.

````

## Status record

````markdown
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
````
