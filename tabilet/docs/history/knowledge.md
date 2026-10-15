# Superseded knowledge

## 2026-10-15 — Authoring workflow files are owned by this checkout

**Sources.** `AGENTS.md` / Memory Bank First and Commands;
`tabilet/memory-bank/milestone.md` / Current State;
`tabilet/memory-bank/tech-stack.md` / Commands.

**Reason.** The previous guidance describes an earlier sibling-owned snapshot.
Current read-only inspection finds the workflow files tracked directly by this
Authoring repository, with ordinary local files and directories. The referenced
`../tofu/authoring` directory is absent. This correction changes current workflow
navigation and commands; it does not revise completed task evidence.

**Evidence.** Observed repository HEAD
`acacd51ad3ea8141b48399e9976f6ef95ae2c148`, including the uncommitted v2.13.0
workflow upgrade. `git rev-parse --show-toplevel`, `git ls-files -- AGENTS.md
tabilet/memory-bank tabilet/evolution`, and filesystem inspection establish the
current ownership. This observation does not claim the uncommitted correction
is present in that commit.

**Replacements.** [Local workflow ownership](../../../AGENTS.md#memory-bank-first),
[local checks](../../../AGENTS.md#commands),
[current milestone state](../../memory-bank/milestone.md#current-state), and
[current stack commands](../../memory-bank/tech-stack.md#commands).

### Previous AGENTS.md / Memory Bank First

```markdown
The tracked canonical Authoring harness snapshot lives in `../tofu/authoring`.
In a normal `../authoring` checkout, `AGENTS.md`, `tabilet/memory-bank/`, and
`tabilet/evolution/` may be symlinks to this tracked snapshot so agents can use local
paths while planning history is committed in the `../tofu` repository.
```

### Previous AGENTS.md / Commands

```bash
git -C ../tofu diff --check -- authoring
git -C ../authoring status --short
```

### Previous milestone.md / Current State

```markdown
The repository exists
as a public Go module scaffold, and
`AGENTS.md`, `tabilet/memory-bank/`, and `tabilet/evolution/` are symlink-facing paths to the
tracked snapshot under `../tofu/authoring`.
```

### Previous tech-stack.md / Commands

```bash
git -C ../tofu diff --check -- authoring
```
