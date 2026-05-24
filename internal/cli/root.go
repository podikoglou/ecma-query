// Package cli root command: global flags, shorthand routing, and the Execute entry point.
package cli

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags (goreleaser).
// Falls back to debug.ReadBuildInfo (go install) then "dev".
var version = "dev"

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}
	rootCmd.Version = version

	rootCmd.PersistentFlags().StringVar((*string)(&format), "format", "json", "output format: json, md")

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

// rootCmd is the top-level cobra command for ecma-query.
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

// rootRun is the RunE handler for the root command. If no subcommand is given
// it shows help. Otherwise it falls through to getCmd (shorthand for "get").
func rootRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return getCmd.RunE(cmd, args)
}
