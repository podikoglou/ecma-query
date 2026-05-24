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
	Error      string        `json:"error"`                // Error is the error type string (e.g. "not_found").
	Query      string        `json:"query,omitempty"`      // Query is the input that caused the error.
	Matches    []MatchResult `json:"matches,omitempty"`    // Matches lists possible matches for disambiguation.
	Wanted     string        `json:"wanted,omitempty"`     // Wanted is the original query string.
	Suggestion string        `json:"suggestion,omitempty"` // Suggestion is a suggested alternative.
}

// MatchResult identifies a single match in a disambiguation response.
type MatchResult struct {
	Kind    string `json:"kind"`              // Kind is the entity kind (e.g. "operation", "method").
	Name    string `json:"name,omitempty"`    // Name is the display name of the match.
	Section string `json:"section,omitempty"` // Section is the section number containing the match.
}

// GetResponse is the JSON response for the "get" command.
type GetResponse struct {
	Kind       string   `json:"kind"`
	Name       string   `json:"name"`
	Section    string   `json:"section"`
	Title      string   `json:"title"`
	Breadcrumb []string `json:"breadcrumb"`
	URL        string   `json:"url"`
	Signature  string   `json:"signature"`
	Summary    string   `json:"summary,omitempty"`
	Steps      []string `json:"steps,omitempty"`
	SeeAlso    []string `json:"see_also,omitempty"`
	Truncated  bool     `json:"truncated"`
}

// SearchResponse is the JSON response for the "search" command.
type SearchResponse struct {
	Query   string         `json:"query"`
	Total   int            `json:"total"`
	Results []SearchResult `json:"results"`
}

// SearchResult is a single search result item.
type SearchResult struct {
	Kind    string `json:"kind"`
	Section string `json:"section"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	URL     string `json:"url"`
}

// TocResponse is the JSON response for the "toc" command.
type TocResponse struct {
	Root    *string    `json:"root"`
	Entries []TocEntry `json:"entries"`
}

// TocEntry represents a single entry in the table of contents.
type TocEntry struct {
	Section  string     `json:"section"`
	Title    string     `json:"title"`
	Children []TocEntry `json:"children"`
}

// XrefResponse is the JSON response for the "xref" command.
type XrefResponse struct {
	Target        string     `json:"target"`
	TargetSection string     `json:"target_section"`
	Incoming      []XrefItem `json:"incoming,omitempty"`
	Outgoing      []XrefItem `json:"outgoing,omitempty"`
}

// XrefItem represents a single cross-reference item.
type XrefItem struct {
	Section string `json:"section"`
	Title   string `json:"title"`
	Context string `json:"context,omitempty"`
}

// GrammarResponse is the JSON response for the "grammar" command.
type GrammarResponse struct {
	Name        string        `json:"name"`
	Kind        string        `json:"kind"`
	Section     string        `json:"section"`
	Productions []GrammarItem `json:"productions"`
}

// GrammarItem represents a single grammar production alternative.
type GrammarItem struct {
	LHS          string     `json:"lhs"`
	Alternatives [][]string `json:"alternatives"`
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

// printJSON marshals v to indented JSON and writes it to stdout.
func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
