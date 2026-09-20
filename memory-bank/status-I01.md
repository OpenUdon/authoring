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
