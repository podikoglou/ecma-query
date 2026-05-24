// Package cli grammar command: grammar production query, defaulting to EBNF output.
package cli

import (
	"github.com/spf13/cobra"
)

// grammarCmd is the "grammar" subcommand for retrieving grammar productions.
var grammarCmd = &cobra.Command{
	Use:   "grammar <PRODUCTION-NAME>",
	Short: "Grammar production query",
	Long:  "Retrieve grammar productions. Defaults to EBNF output (instead of the global JSON default).",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}

var grammarFormat string // grammarFormat is the output format for grammar results (--format).

func init() {
	grammarCmd.Flags().StringVar(&grammarFormat, "format", "ebnf", "json, md, or ebnf")
}
