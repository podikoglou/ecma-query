package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/podikoglou/ecma-query/internal/spec"
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

var (
	getDepth     int
	getBrief     bool
	getMaxTokens int
	getChunk     int
)

func init() {
	getCmd.Flags().IntVar(&getDepth, "depth", 1, "recursion depth for inlining nested operations")
	getCmd.Flags().BoolVar(&getBrief, "brief", false, "return signature and summary only (--depth 0)")
	getCmd.Flags().IntVar(&getMaxTokens, "max-tokens", 0, "soft truncate output at approximately N tokens")
	getCmd.Flags().IntVar(&getChunk, "chunk", 1, "when previous output was truncated, request chunk N")
}

func runGet(_ *cobra.Command, args []string) error {
	query := args[0]
	s := getSpec()

	if getBrief {
		getDepth = 0
	}

	cn, multi, matches := resolve(s, query)

	if cn == nil && len(multi) > 1 {
		ExitAmbiguous(query, matches)
		return nil
	}

	if cn == nil {
		succ := suggestions(query, s)
		resp := ErrorResponse{
			Error:  "not_found",
			Query:  query,
			Wanted: query,
		}
		if len(succ) > 0 {
			resp.Suggestion = succ[0]
		}
		writeError(CodeNotFound, resp)
		return nil
	}

	if format == FormatJSON {
		resp := buildGetJSON(cn, s)
		printJSON(resp)
	} else {
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
				if cut > 0 && getChunk > 1 {
					start := cut * (getChunk - 1)
					if start >= len(lines) {
						start = 0
					}
					end := start + cut
					if end > len(lines) {
						end = len(lines)
					}
					md = strings.Join(lines[start:end], "\n")
				} else if cut > 0 {
					md = strings.Join(lines[:cut], "\n")
				}
			}
		}
		fmt.Print(md)
	}
	return nil
}

func buildGetJSON(cn *spec.ClauseNode, s *spec.Spec) GetResponse {
	resp := GetResponse{
		Kind:       kindLabel(cn),
		Name:       cn.Name,
		Section:    cn.Section,
		Title:      cn.H1Text,
		Breadcrumb: buildBreadcrumb(cn),
		URL:        clauseURL(cn.ID),
		Signature:  cn.H1Text,
	}

	resp.Summary = extractSummary(cn)

	steps := extractSteps(cn)
	if getDepth > 0 && len(steps) > 0 {
		resp.Steps = steps
	}

	if s.Outgoing[cn.ID] != nil {
		for targetID := range s.Outgoing[cn.ID] {
			if target, ok := s.ByID[targetID]; ok {
				name := target.Name
				if name == "" {
					name = target.H1Text
				}
				if name != "" {
					resp.SeeAlso = append(resp.SeeAlso, name)
				}
			}
		}
	}

	if getMaxTokens > 0 {
		text := resp.Summary
		for _, s := range resp.Steps {
			text += s
		}
		tokens := estimateTokens(text)
		if tokens > getMaxTokens {
			resp.Truncated = true
			cut := getMaxTokens
			if len(resp.Steps) > cut {
				resp.Steps = resp.Steps[:cut]
			}
		}
	}

	return resp
}
