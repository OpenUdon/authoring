# Retired milestone M22 - Documentation And Examples

**Milestone.** M22
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M22.md
**Source specification.** tabilet/memory-bank/milestone.md#m22-documentation-and-examples
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
## M22 Documentation And Examples

### M22 Documentation And Examples

**Goal.** Add Authoring README/package docs and minimal fake-runtime examples.

Acceptance:

- Docs cover manual prompting, agent mode, structured JSON fallback,
  progressive loop, and report/digest generation.
- Examples use fake runtimes and fake clients only.
- Documentation restates boundaries: no OpenUdon/Ramen semantics, no live model
  or executor requirement, and no credential storage.
- Examples compile and run in default `go test ./...`.

````

## Status record

````markdown
# Status M22

## Goal

Add Authoring documentation and fake-runtime examples.

## Tasks

| Item | State | Notes |
|---|---|---|
| Add README and package docs | `[+]` | Added `README.md`, expanded root package docs, and linked compatibility policy. |
| Add manual prompting example | `[+]` | Added runnable `Example_manualPrompting` using in-memory input/output only. |
| Add agent and structured-output examples | `[+]` | Added runnable structured JSON fallback example with a fake client and no provider credentials. |
| Add progressive loop and report examples | `[+]` | Added runnable fake-runtime progressive loop and report/digest examples. |

## Notes

- Examples must not require model credentials, API execution, workflow
  execution, Terraform/OpenTofu execution, or trusted runner access.
- Examples are compiled and executed by `go test ./...`.

## Verification

Completed checks:

```bash
(cd ../authoring && go test ./... && go vet ./...)
(cd ../authoring && GOWORK=off go test ./... && GOWORK=off go vet ./...)
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
(cd ../authoring && git diff --check)
git -C ../tofu diff --check -- authoring
```
````
