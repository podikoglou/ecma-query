package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

type OutputFormat string

const (
	FormatJSON     OutputFormat = "json"
	FormatMarkdown OutputFormat = "md"
	FormatEBNF     OutputFormat = "ebnf"
)

type ErrorCode int

const (
	CodeSuccess   ErrorCode = 0
	CodeNotFound  ErrorCode = 1
	CodeAmbiguous ErrorCode = 2
)

var format OutputFormat = FormatJSON
var noColor = true

type ErrorResponse struct {
	Error   string        `json:"error"`
	Query   string        `json:"query,omitempty"`
	Matches []MatchResult `json:"matches,omitempty"`
}

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

func ExitNotFound(query string) {
	writeError(CodeNotFound, ErrorResponse{
		Error: "not_found",
		Query: query,
	})
}

func ExitAmbiguous(query string, matches []MatchResult) {
	writeError(CodeAmbiguous, ErrorResponse{
		Error:   "ambiguous",
		Query:   query,
		Matches: matches,
	})
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
