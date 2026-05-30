// Package authoring provides shared authoring orchestration primitives.
//
// The root package is intentionally small. Public behavior lives in
// subpackages grouped by concern: trust, session, transcript, prompt,
// lifecycle, structured, icot, readiness, decision, report, and promptcontext.
// Downstream products bind those generic contracts to product-specific
// prompts, schemas, validation, artifacts, credentials, model clients, and
// execution boundaries.
package authoring
