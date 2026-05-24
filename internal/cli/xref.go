// Package cli xref command: finds incoming and outgoing references for an entity.
package cli

import (
	"github.com/spf13/cobra"
)

// xrefCmd is the "xref" subcommand for resolving cross-references.
var xrefCmd = &cobra.Command{
	Use:   "xref <IDENTIFIER>",
	Short: "Cross-reference resolution",
	Long:  "Find what references a given entity and what it references.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}

var xrefDirection string // xrefDirection filters references: incoming, outgoing, or both (--direction).

func init() {
	xrefCmd.Flags().StringVar(&xrefDirection, "direction", "both", "incoming, outgoing, or both")
}
