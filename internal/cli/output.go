// Package cli provides the command-line interface for ecma-query.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

// OutputFormat is the format for command output.
type OutputFormat string

// Output format constants.
const (
	FormatJSON     OutputFormat = "json" // FormatJSON produces JSON output.
	FormatMarkdown OutputFormat = "md"   // FormatMarkdown produces Markdown output.
	FormatEBNF     OutputFormat = "ebnf" // FormatEBNF produces EBNF output.
)

// ErrorCode is an exit code for structured error responses.
type ErrorCode int

// Exit code constants.
const (
	CodeSuccess   ErrorCode = 0 // CodeSuccess indicates success.
	CodeNotFound  ErrorCode = 1 // CodeNotFound indicates the entity was not found.
	CodeAmbiguous ErrorCode = 2 // CodeAmbiguous indicates the query matched multiple entities.
)

// format is the global output format, configurable via --format on the root command.
var format = FormatJSON

// ErrorResponse is a structured error returned to stdout.
type ErrorResponse struct {
	Error   string        `json:"error"`             // Error is the error type string (e.g. "not_found").
	Query   string        `json:"query,omitempty"`   // Query is the input that caused the error.
	Matches []MatchResult `json:"matches,omitempty"` // Matches lists possible matches for disambiguation.
}

// MatchResult identifies a single match in a disambiguation response.
type MatchResult struct {
	Kind    string `json:"kind"`              // Kind is the entity kind (e.g. "operation", "method").
	Name    string `json:"name,omitempty"`    // Name is the display name of the match.
	Section string `json:"section,omitempty"` // Section is the section number containing the match.
}

func writeError(code ErrorCode, resp ErrorResponse) {
	respBytes, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(respBytes))
	os.Exit(int(code))
}

// ExitNotFound writes a not_found error and exits with code 1.
func ExitNotFound(query string) {
	writeError(CodeNotFound, ErrorResponse{
		Error: "not_found",
		Query: query,
	})
}

// ExitAmbiguous writes an ambiguous error with matches and exits with code 2.
func ExitAmbiguous(query string, matches []MatchResult) {
	writeError(CodeAmbiguous, ErrorResponse{
		Error:   "ambiguous",
		Query:   query,
		Matches: matches,
	})
}
