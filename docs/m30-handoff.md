# M30 accepted source handoff — 2026-10-02

Authoring's iCoT/icotcli compatibility surfaces are retired. Accepted removal
source: `db4f5193bc53819be2b1bf1d8728735941e28d9e`. Qualification record commit:
`9241c481a281ab31f3e4b6ae89f9242ae9517c6c`, with identical implementation.
[M30 status](../tabilet/memory-bank/status-M30.md) is the package-local accepted
record; it stays in Authoring's completed ledger under the established convention.
[Qualification](m30-qualification.md) records exact consumer revisions/resolution,
commands, composed passing gates and failed evidence. Deep review passed 1/10
with no open P1/P2 findings.

## Consumer dispositions

- Kinet at `2970731766af13657d464983791c6c17e6312c0d` and OpenUdon at
  `fbda7e9231b8b306fd1ae3ac623e9d70331b3e08` passed their required checks
  against this exact new Authoring source with temporary modfile overrides.
  Neither needs an implementation or manifest change. Normal dependency pins
  remain unchanged; publication and pin adoption need separate authority.
- Ramen at `279a099a418bf324c693c5390ca2a64b967e5669` remains frozen at
  `v0.0.0-20260820042256-2f73e3526583` and retains its iCoT imports. Its
  readonly standalone build passed offline without replacement. It intentionally
  cannot build against current workspace Authoring. The operator controls any
  parent `go.work` exclusion/update; this milestone changes neither file nor Ramen.
- OpenUdon's remaining legacy ICOT names are its owner's deferred naming
  handoff. No runtime/wire or authority change is implied by renaming later.
- W8M's active check module has neither removed-package import. It is untouched.
  Its owner must separately reconcile a changed Kinet binary before adoption;
  older qualification does not qualify a Stage 6 binary.

All surviving neutral APIs, implementation/fixture bytes, durable JSON versions,
prompts/transcripts, safety and state behavior are preserved. Kinet's OpenUdon
runtime pin remains `c2f161d762bc9f2217bbf0c34b00cdef64b0f7d0`. All other sibling
ledgers retain their owners. No push, publication, tag or deployment was performed.
