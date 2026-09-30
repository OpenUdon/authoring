# Status M29 — Neutral engine and compatible iCoT facade

State: M29.1–M29.3 complete; M29.4 closing review/publication in progress.

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
| M29.3 Compatibility and contract verification | `[+]` | Workspace/standalone full test/vet, race, dependency boundary, parity, unchanged OpenUdon/Ramen consumers and OpenUdon scorecard. |
| M29.4 Review, publication and downstream handoff | `[~]` | Persist bounded deep review, qualified source/module and exact-diff publication; reconcile OpenUdon M91.3 before resumption. |

## Verification and acceptance

See the full specification in milestone.md. Default checks use fake providers
and disposable fixtures. Engine cannot depend on icot or downstream products.
All old signatures, constants, errors, JSON versions, prompt/transcript bytes
and atomic state/approval behavior remain compatible. Retain current local
completed-ledger convention; no historical milestone is reopened or retired.

## Closing review

Persisted iteration count: 1/10. Iteration 1 passed; no open P1/P2 findings. Task completion alone does not
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

## M29.3 completed — compatibility qualification

Qualified implementation source: `18056cb6b0c1007dd567a4a825a6b4311a357185`.
All 66 old exported declarations/function signatures match the exact published
baseline. All 12 moved implementation/regression bodies remain identical
apart from package-name/doc normalization; durable versions and tracked fixture
bytes were not modified.

`./scripts/check-compat.sh` passed Authoring workspace/standalone full tests/vet,
import boundaries, OpenUdon full tests and Ramen full tests (including its
304.817-second corpus package). Full Authoring race tests, separate standalone vet,
APItools full tests and OpenUdon's 103-pass/zero-failure scorecard and report
verification pass. Logs are `/tmp/authoring-m29-3-compat.log`,
`/tmp/authoring-m29-3-race.log`, `/tmp/authoring-m29-3-standalone-vet.log`,
`/tmp/apitools-m29-3-compat.log`, `/tmp/openudon-m29-3-scorecard.log`.

Frozen compatibility evidence:
`/var/tmp/authoring-m29-compat-0w5o3ztc/manifest.json` binds exact Authoring,
OpenUdon and Ramen source/archive hashes, consumer module metadata, toolchain
and completed log digests. Frozen Authoring full test/vet, OpenUdon full tests
and Ramen's affected authoring package pass with GOWORK off, GOPROXY off and
only Authoring replaced by its frozen source; other consumer pins are unchanged.
The installed Go 1.26.6 binary is used explicitly for consumer checks (the
initial wrapper invocation refused its checksum configuration before testing).
No mutable sibling substitution or model/runtime account operation qualified
this source. Ramen and APItools source/ledgers remain unchanged. Closing review
and exact publication remain M29.4 work.

## Closing review iteration 1 — passed

Reviewed exact source/fixture equivalence, all 66 public declarations, alias
method/type identity, generic forwarder arguments/returns, sentinels, prompt and
transcript behavior, atomic binding, cancellation/no-progress/repair safety,
neutral dependency direction and unchanged consumers. No open P1/P2 findings;
no code fix was needed. Reflection package ownership and conventional sentinel
non-reassignment are documented compatibility limits. The full mandatory
verification and frozen compatibility gates passed; no live operation occurred.

Qualified implementation remains `18056cb6b0c1007dd567a4a825a6b4311a357185`;
qualification manifest SHA-256 `e15cb0cfb36d7c8f62223c89aa3fb08b102c92672565d38de495bae3f75e73d0`. M29.4 is in progress for scoped
publication and exact downstream handoff. This passed review does not imply
that a remote or consumer was already updated.
