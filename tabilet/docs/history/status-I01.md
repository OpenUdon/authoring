# Retired milestone I01 - Interview, Transcript, Prompt, And iCoT Integrity

**Milestone.** I01
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-I01.md
**Source specification.** tabilet/memory-bank/milestone.md#i01-interview-transcript-prompt-and-icot-integrity
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
## I01 Interview, Transcript, Prompt, And iCoT Integrity

### I01 Interview, Transcript, Prompt, And iCoT Integrity

**Goal.** Make every loop entry point share one complete-round interview
transaction, semantic progress policy, chronological event stream, and
persistence-validation contract.

Ownership and dependency:

- Owns `interview`, `icot`, `prompt`, `session`, `transcript`, and readiness
  compatibility behavior for this slice.
- Exposes `session.ValidateForPersistence`; D01 consumes that contract and
  does not redefine sensitive-value policy.
- OpenUdon and Ramen own their binding callbacks, prompt wording, product
  mutations, readiness rules, and artifact formats.

Acceptance:

- `interview.ApplyRound` resolves every current frontier node exactly once by
  answer or complete deferral, appends evidence atomically, advances one round,
  resets the no-progress counter, and preserves the input on failure.
- `icot.InterviewBinding` is usable through generic, runtime, and interactive
  entry points; OpenUdon and Ramen remove duplicate settlement plumbing.
- Semantic fingerprints return explicit errors, event IDs are monotonic,
  transcript event order survives normalize/save/load, and the emergency fuse
  defaults to 1,000 rounds when `MaxRounds` is zero.
- Interactive planning reuses an unchanged provisional frontier, replans once
  after declared draft mutation, attributes legacy answers to the displayed
  frontier, and emits no phantom decision events.
- Prompt messages/frontier text/opening labels are configurable, only
  user-sourced `cancel` answers cancel, nil contexts/sessions fail
  descriptively, and readiness aliases normalize write-through.

````

## Status record

````markdown
# Status I01 - Interview, Transcript, Prompt, And iCoT Integrity

State: Complete

## Dependencies And Ownership

I01 owns `interview`, `icot`, `prompt`, `session`, `transcript`, and readiness
compatibility changes in this remediation. D01 depends on the public
`session.ValidateForPersistence` contract. OpenUdon and Ramen own non-overlapping
product bindings, prompts, readiness policy, and artifact behavior.

## Tasks

| Item | State | Notes |
|---|---|---|
| I01.1 atomic interview settlement | `[+]` | Added complete answer/deferral `Resolution` transactions, evidence validation, iterative cycle detection, normalized internal helpers, indexed node settlement, and failure atomicity/deep-graph/20,000-node frontier tests. |
| I01.2 shared interview binding | `[+]` | Added clone-based generic `InterviewBinding`, exposed it through options and the alternative runtime interface, and migrated OpenUdon and Ramen settlement bridges. |
| I01.3 loop integrity | `[+]` | Added semantic fingerprints, explicit JSON errors, `MaxRounds` with a 1,000-round default fuse, exact answer IDs, user-only cancellation, last-round result semantics, and synchronized no-progress counters. |
| I01.4 chronology and events | `[+]` | Exported event callbacks, assigned monotonic IDs to engine/repair events, and made transcript chronology authoritative through normalization and persistence. |
| I01.5 interactive planning | `[+]` | Cached provisional plans, invalidated them when a recoverable draft hook discards its draft, replanned once after declared mutation, bound legacy answers to the displayed frontier, removed phantom final decisions, and added configurable prompt/frontier/opening text. |
| I01.6 validation and compatibility | `[+]` | Added descriptive nil handling, delegated replay-label checks, documented write-through readiness aliases and answer semantics, and added provider-free regression coverage. |

## Verification

- Authoring package and full suites pass.
- OpenUdon full suite passes.
- Ramen full suite passes.
- Authoring standalone, race, vet, diff, compatibility, OpenUdon scorecard, and
  workspace downstream gates pass.
- OpenUdon and Ramen pin the published Apitools and Authoring revisions and
  pass `GOWORK=off go test ./...` plus `GOWORK=off go vet ./...` without local
  replacements.

Commit tracking: Authoring interview/loop implementation `243121d`; coordinated
D01 `eee0c0d`, Authoring boundary removal `2f73e35`, Apitools `d51b61e`,
OpenUdon `43294bf`/`065bb84`, and Ramen `330fd89`/`721b540`.
````
