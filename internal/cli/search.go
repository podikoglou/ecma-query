package cli

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <QUERY>",
	Short: "Full-text search across the spec",
	Long:  "Returns ranked results with snippets for exploratory queries.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var searchLimit int
var searchKind string

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "maximum number of results")
	searchCmd.Flags().StringVar(&searchKind, "kind", "", "filter by entity kind: operation, method, section, type, grammar, slot")
}
