package cli

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "ecma-query [command]",
	Short: "Query the ECMAScript specification",
	Long: `ecma-query is a CLI for querying the ECMAScript specification,
designed for LLM agent consumption.

Output is deterministic, structured, and parseable. All output,
including errors, goes to stdout. No ANSI codes, no prompts.`,
	Args:    cobra.ArbitraryArgs,
	RunE:    rootRun,
	Version: "0.1.0",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.Flags().StringVar((*string)(&format), "format", "json", "output format: json, md")
	rootCmd.Flags().String("spec", "latest", "spec edition: es2024, es2025, latest")
	rootCmd.Flags().Bool("no-color", true, "disable ANSI codes")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(tocCmd)
	rootCmd.AddCommand(xrefCmd)
	rootCmd.AddCommand(grammarCmd)
}

func rootRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return getCmd.RunE(cmd, args)
}
