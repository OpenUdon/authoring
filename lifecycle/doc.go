// Package lifecycle provides generic draft, transcript, atomic-write, and
// artifact-safety helpers for Authoring flows.
//
// Product adapters own draft schemas and generated artifact contents. This
// package owns only durable local persistence mechanics and safe metadata
// records.
package lifecycle
