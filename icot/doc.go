// Package icot provides a generic progressive authoring loop.
//
// The loop owns sequencing, prompt/default handling, bounded attempts, and
// transcript events. Downstream runtimes own draft schemas, readiness checks,
// question planning, answer application, final artifact generation, and product
// semantics.
package icot
