// Package promptcontext defines product-neutral, prompt-safe context records.
//
// Downstream OpenUdon and Ramen adapters translate API metadata, UWS schemas,
// and credential bindings into these narrow summaries. This package does not
// import apitools, UWS, OpenUdon, Ramen, provider clients, or executor code.
package promptcontext
