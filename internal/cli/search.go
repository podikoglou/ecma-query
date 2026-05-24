// Package cli search command: full-text search across the spec with ranked results.
package cli

import (
	"github.com/spf13/cobra"
)

// searchCmd is the "search" subcommand for full-text exploration across the spec.
var searchCmd = &cobra.Command{
	Use:   "search <QUERY>",
	Short: "Full-text search across the spec",
	Long:  "Returns ranked results with snippets for exploratory queries.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}

var (
	searchLimit int    // searchLimit is the maximum number of results (--limit).
	searchKind  string // searchKind filters results by entity kind (--kind).
)

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "maximum number of results")
	searchCmd.Flags().StringVar(&searchKind, "kind", "", "filter by entity kind: operation, method, section, type, grammar, slot")
}
