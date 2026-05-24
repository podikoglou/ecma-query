package cli

import (
	"fmt"
	"os"
)

// ErrorCode is an exit code for error conditions.
type ErrorCode int

// Exit code constants.
const (
	CodeSuccess   ErrorCode = 0 // CodeSuccess indicates success.
	CodeNotFound  ErrorCode = 1 // CodeNotFound indicates the entity was not found.
	CodeAmbiguous ErrorCode = 2 // CodeAmbiguous indicates the query matched multiple entities.
)

// MatchResult identifies a single match in a disambiguation response.
type MatchResult struct {
	Kind    string
	Name    string
	Section string
}

func exitNotFound(query string) {
	fmt.Println("not found: " + query)
	os.Exit(int(CodeNotFound))
}

func exitAmbiguous(query string, matches []MatchResult) {
	fmt.Println("ambiguous: " + query)
	for _, m := range matches {
		line := "  - " + m.Kind
		if m.Name != "" {
			line += " " + m.Name
		}
		if m.Section != "" {
			line += " (" + m.Section + ")"
		}
		fmt.Println(line)
	}
	os.Exit(int(CodeAmbiguous))
}
