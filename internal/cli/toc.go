// Package cli toc command: table of contents and structural navigation.
package cli

import (
	"github.com/spf13/cobra"
)

// tocCmd is the "toc" subcommand that prints the spec outline or children of a section.
var tocCmd = &cobra.Command{
	Use:   "toc [SECTION]",
	Short: "Table of contents / structural navigation",
	Long:  "Without arguments: top-level outline. With a section number: children of that section.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}

var tocDepth int // tocDepth is the number of levels to show (--depth).

func init() {
	tocCmd.Flags().IntVar(&tocDepth, "depth", 2, "number of levels to show")
}
