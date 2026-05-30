// Package trust anchors Authoring's neutral trust and safety records to the
// shared Evidence module.
//
// The package intentionally aliases Evidence record types and delegates helper
// behavior to Evidence packages. Authoring-specific session, transcript,
// lifecycle, and report packages can depend on these names without defining a
// second durable diagnostic, redaction, artifact, or digest contract.
package trust
