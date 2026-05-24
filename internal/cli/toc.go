package cli

import (
	"github.com/spf13/cobra"
)

var tocCmd = &cobra.Command{
	Use:   "toc [SECTION]",
	Short: "Table of contents / structural navigation",
	Long:  "Without arguments: top-level outline. With a section number: children of that section.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}

var tocDepth int

func init() {
	tocCmd.Flags().IntVar(&tocDepth, "depth", 2, "number of levels to show")
}
