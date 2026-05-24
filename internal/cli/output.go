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
	FormatJSON     OutputFormat = "json"
	FormatMarkdown OutputFormat = "md"
	FormatEBNF     OutputFormat = "ebnf"
)

// ErrorCode is an exit code for structured error responses.
type ErrorCode int

// Exit code constants.
const (
	CodeSuccess   ErrorCode = 0
	CodeNotFound  ErrorCode = 1
	CodeAmbiguous ErrorCode = 2
)

var format = FormatJSON

// ErrorResponse is a structured error returned to stdout.
type ErrorResponse struct {
	Error   string        `json:"error"`
	Query   string        `json:"query,omitempty"`
	Matches []MatchResult `json:"matches,omitempty"`
}

// MatchResult identifies a single match in a disambiguation response.
type MatchResult struct {
	Kind    string `json:"kind"`
	Name    string `json:"name,omitempty"`
	Section string `json:"section,omitempty"`
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
