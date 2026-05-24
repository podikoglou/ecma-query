package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
	RunE: runGet,
}

var getMaxTokens int

func init() {
	getCmd.Flags().IntVar(&getMaxTokens, "max-tokens", 0, "soft truncate output at approximately N tokens")
}

func runGet(_ *cobra.Command, args []string) error {
	query := args[0]
	s := getSpec()

	cn, multi, matches := resolve(s, query)

	if cn == nil && len(multi) > 1 {
		exitAmbiguous(query, matches)
		return nil
	}

	if cn == nil {
		succ := suggestions(query, s)
		msg := "not found: " + query
		if len(succ) > 0 {
			msg += "\nsuggestion: " + succ[0]
		}
		fmt.Println(msg)
		os.Exit(int(CodeNotFound))
		return nil
	}

	md, err := s.RenderNode(cn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "render error:", err)
		os.Exit(1)
	}
	if getMaxTokens > 0 {
		tokens := estimateTokens(md)
		if tokens > getMaxTokens {
			lines := strings.Split(md, "\n")
			cut := 0
			cur := 0
			for i, line := range lines {
				cur += estimateTokens(line)
				if cur > getMaxTokens {
					cut = i
					break
				}
			}
			if cut > 0 {
				md = strings.Join(lines[:cut], "\n")
			}
		}
	}
	fmt.Print(md)
	return nil
}
