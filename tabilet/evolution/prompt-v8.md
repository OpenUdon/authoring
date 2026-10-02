# Prompt V8 — Intentional iCoT facade retirement

Approved Stage 6 planning-handoff amendments, 2026-10-02; successor to [v7](prompt-v7.md). Retire `authoring/icot` and `authoring/icotcli` through [M30](../memory-bank/status-M30.md) as an intentional pre-1.0 source break. Preserve the neutral engine and all non-removed package APIs, durable versions/tags, prompt/transcript bytes, safety and state behavior.

Ramen retains its old Authoring pin and both facade imports; it does not drop the dependency. The approved exception ends compatibility with current workspace Authoring for those packages. Check frozen Ramen standalone at its actual pin, and qualify unchanged Kinet/OpenUdon against exact new Authoring through temporary dependency overrides without editing consumers or operator go.work. Record effective resolution and exact source. Update current compatibility/API/architecture/product/stack instructions when removal is implemented.

Keep four M30 rows pending until separately requested execution. Kinet coordinates W12 → U09 → Authoring M30 with one execution owner and package-local ledgers; M30 is code-independent of the Kinet milestones. Task commit policy and external-mutation authority belong to the later goal request. Planning does not authorize publication or consumer pin adoption. See [result v8](result-v8.md) for current state.
