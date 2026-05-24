package cli

import (
	"github.com/spf13/cobra"
)

var grammarCmd = &cobra.Command{
	Use:   "grammar <PRODUCTION-NAME>",
	Short: "Grammar production query",
	Long:  "Retrieve grammar productions. Defaults to EBNF output (instead of the global JSON default).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var grammarFormat string

func init() {
	grammarCmd.Flags().StringVar(&grammarFormat, "format", "ebnf", "json, md, or ebnf")
}
