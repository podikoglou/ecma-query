package cli

import (
	"github.com/spf13/cobra"
)

var xrefCmd = &cobra.Command{
	Use:   "xref <IDENTIFIER>",
	Short: "Cross-reference resolution",
	Long:  "Find what references a given entity and what it references.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}

var xrefDirection string

func init() {
	xrefCmd.Flags().StringVar(&xrefDirection, "direction", "both", "incoming, outgoing, or both")
}
