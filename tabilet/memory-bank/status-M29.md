# Status M29 — Neutral engine and compatible iCoT facade

State: M29.1/M29.2 complete; verification pending.

## Goal and dependencies

Expose the existing generic mechanics through `authoring/engine`, preserving
old `authoring/icot` consumers through aliases/forwarders and one implementation.
Baseline: `b417eb681746476cca80bdb14cde3e3c96de980c`; I01, D01 and M28 complete.
This package owns M29; Kinet coordinates the goal without merging ledgers.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| M29.1 Neutral implementation extraction | `[+]` | Relocate generic loop, interactive/runtime/repair and atomic-binding implementation and regression tests; preserve bodies and wire/safety behavior. |
| M29.2 Compatible iCoT facade | `[+]` | Retain every exported old name through type aliases/forwarders and shared sentinel identity; no second controller. |
| M29.3 Compatibility and contract verification | `[ ]` | Workspace/standalone full test/vet, race, dependency boundary, parity, unchanged OpenUdon/Ramen consumers and OpenUdon scorecard. |
| M29.4 Review, publication and downstream handoff | `[ ]` | Persist bounded deep review, qualified source/module and exact-diff publication; reconcile OpenUdon M91.3 before resumption. |

## Verification and acceptance

See the full specification in milestone.md. Default checks use fake providers
and disposable fixtures. Engine cannot depend on icot or downstream products.
All old signatures, constants, errors, JSON versions, prompt/transcript bytes
and atomic state/approval behavior remain compatible. Retain current local
completed-ledger convention; no historical milestone is reopened or retired.

## Closing review

Persisted iteration count: 0/10. Not started. Task completion alone does not
establish acceptance. Record exact source/module, observed publication, checks,
findings and downstream reconciliation before completion.

## Approval and execution authority

The user approved the complete proposal and extended the existing Stage 5 goal
on 2026-09-30, including Authoring publication. Task commits and scoped normal
fast-forward origin/main publication use the coordinator's automatic exact-diff
policy. Authoring iCoT/icotcli retirement remains deferred. Ramen is read-only.

## M29.1 selected

Relocate the existing implementation and regression tests as one package.
A minimal alias/forwarder facade is required in the same compiling extraction
change to keep old consumers usable; M29.2 owns exhaustive compatibility checks
and facade documentation. Preserve bodies after package-name normalization.

## M29.1 completed

All 12 relocated implementation/regression files match published baseline
`b417eb681746476cca80bdb14cde3e3c96de980c` after package declaration and
package-doc name normalization. The compatibility facade preserves all old
exported names through aliases/forwarders and original sentinel values.

`go test ./engine ./icot`, standalone full tests/vet, and unchanged OpenUdon
internal authoring/elicitor/iCoT and Ramen authoring suites pass. Logs:
`/tmp/authoring-m29-1-focused.log`, `/tmp/authoring-m29-1-full.log`,
`/tmp/authoring-m29-1-vet.log`, `/tmp/openudon-m29-1-alias.log`,
`/tmp/ramen-m29-1-alias.log`. M29.2 still owns exhaustive facade compatibility
and documentation; M29.3/full downstream and M29.4 acceptance remain pending.

## M29.2 completed

The facade retains all exported pre-extraction declarations and forwards calls
without a second loop. Compile-time checks cross generic hook/runtime/repair
signatures; facade tests prove concrete type identity, shared error sentinels,
loop success/needs-input/cancellation parity, and exact prompt output/turn parity.
The dependency test rejects any neutral-engine dependency on the facade.
README, package docs and compatibility policy describe the additive API and
reflection ownership caveat. Durable persistence continues to use explicit
wire versions rather than Go reflection package identity.

Full workspace tests and diff checks pass; log
`/tmp/authoring-m29-2-compat-fixed.log`. Initial test-only mistakes (a typed
runtime wrongly asserted as `any`, and an extra blank input before a forced
prompt) were corrected against the unchanged actual API; implementation bodies
remain identical. M29.3 still owns full cross-package/standalone/race gates.
