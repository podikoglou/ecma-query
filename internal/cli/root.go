// Package cli root command: shorthand routing and the Execute entry point.
package cli

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

var version = "dev"

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}
	rootCmd.Version = version

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(tocCmd)
	rootCmd.AddCommand(xrefCmd)
	rootCmd.AddCommand(grammarCmd)
}

// Execute runs the root command and returns any error.
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
	Args:          cobra.ArbitraryArgs,
	RunE:          rootRun,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func rootRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return getCmd.RunE(cmd, args)
}
