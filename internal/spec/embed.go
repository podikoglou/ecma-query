package spec

// This file embeds the spec HTML file and exposes it for parsing.

import _ "embed"

// specHTML holds the embedded ECMAScript specification HTML bytes.
// It is loaded at compile time via //go:embed and serves as the sole input
// for all parsing, indexing, and rendering operations.
//
//go:embed spec.html
var specHTML []byte
