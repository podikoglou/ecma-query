// Package cli get command: retrieves a spec entity by exact identifier.
package cli

import (
	"github.com/spf13/cobra"
)

// getCmd is the "get" subcommand — exact identifier resolution against Section number,
// anchor ID, abstract operation name, built-in method path, internal slot, or spec type.
var getCmd = &cobra.Command{
	Use:   "get <IDENTIFIER>",
	Short: "Retrieve a spec entity by exact identifier",
	Long: `Resolves an exact identifier. No fuzzy matching.

Resolution order:
  1. Section number (e.g. 7.1.4)
  2. Anchor ID (e.g. sec-toprimitive)
  3. Abstract operation name (e.g. ToNumber)
  4. Built-in method path (e.g. Array.prototype.map)
  5. Internal slot (e.g. [[Prototype]])
  6. Spec type (e.g. PropertyDescriptor)`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		ExitNotFound(args[0])
		return nil
	},
}

var (
	getDepth     int  // getDepth is the recursion depth for inlining nested operations (--depth).
	getStepsOnly bool // getStepsOnly returns only the numbered algorithm steps (--steps-only).
	getBrief     bool // getBrief returns signature and summary only, equivalent to --depth 0.
	getMaxTokens int  // getMaxTokens is the soft token truncation limit for output (--max-tokens).
	getChunk     int  // getChunk requests a specific chunk of previously truncated output (--chunk).
)

func init() {
	getCmd.Flags().IntVar(&getDepth, "depth", 1, "recursion depth for inlining nested operations")
	getCmd.Flags().BoolVar(&getStepsOnly, "steps-only", false, "return only the numbered algorithm steps")
	getCmd.Flags().BoolVar(&getBrief, "brief", false, "return signature and summary only (--depth 0)")
	getCmd.Flags().IntVar(&getMaxTokens, "max-tokens", 0, "soft truncate output at approximately N tokens")
	getCmd.Flags().IntVar(&getChunk, "chunk", 1, "when previous output was truncated, request chunk N")
}
