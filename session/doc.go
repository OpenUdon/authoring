// Package session defines durable, product-neutral authoring session records.
//
// Downstream products map their own prompts, goals, readiness checks, and
// generated artifact summaries into these structs. The package does not import
// OpenUdon, Ramen, UWS, API-source packages, provider clients, or executor
// code.
package session
