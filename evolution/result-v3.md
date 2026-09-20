# Result V3 - Adaptive Evidence-Grounded Authoring

M25 introduces `authoring.interview.v1` and one active frontier-round engine.
The interview package validates graph structure, evidence references, status
transitions, and structured deferrals, and returns all currently ready nodes in
deterministic priority/ID order. Independent answers can be applied atomically.

Generic, runtime-bound, and interactive iCoT entry points now share the same
engine. It renders the complete numbered round before input, applies the round
as one mutation followed by one normalization/autosave, supports more than 20
decisions without a breadth ceiling, and diagnoses three consecutive
no-progress rounds. Full asks every question, normal visibly accepts safe
recommendations, and fast hides safe defaults while still showing forced or
missing decisions.

`readiness.Question` now serializes one `forced` flag, one `recommendation`,
priority, concise rationale, and evidence references. Ramen remains unchanged
through source-only compatibility aliases. OpenUdon owns the separate v2
session/report/replay migration and product-specific adaptive decision tree.
