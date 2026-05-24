package cli

import (
	"fmt"

	"github.com/podikoglou/ecma-query/internal/spec"
	"github.com/spf13/cobra"
)

var xrefCmd = &cobra.Command{
	Use:   "xref <IDENTIFIER>",
	Short: "Cross-reference resolution",
	Long:  "Find what references a given entity and what it references.",
	Args:  cobra.ExactArgs(1),
	RunE:  runXref,
}

var xrefDirection string

func init() {
	xrefCmd.Flags().StringVar(&xrefDirection, "direction", "both", "incoming, outgoing, or both")
}

func runXref(_ *cobra.Command, args []string) error {
	query := args[0]
	s := getSpec()

	cn, multi, matches := resolve(s, query)
	if cn == nil && len(multi) > 1 {
		exitAmbiguous(query, matches)
		return nil
	}
	if cn == nil {
		exitNotFound(query)
		return nil
	}

	xrefs := buildXrefPairs(s, cn)
	renderXrefMD(cn, xrefs)
	return nil
}

type xrefPair struct {
	item xrefItem
	isIn bool
}

type xrefItem struct {
	Section string
	Title   string
	Context string
}

func buildXrefPairs(s *spec.Spec, cn *spec.ClauseNode) []xrefPair {
	var pairs []xrefPair

	if xrefDirection == "both" || xrefDirection == "incoming" {
		for sourceID := range s.Incoming[cn.ID] {
			if source, ok := s.ByID[sourceID]; ok {
				pairs = append(pairs, xrefPair{
					item: xrefItem{
						Section: source.Section,
						Title:   source.H1Text,
					},
					isIn: true,
				})
			}
		}
	}

	if xrefDirection == "both" || xrefDirection == "outgoing" {
		for targetID := range s.Outgoing[cn.ID] {
			if target, ok := s.ByID[targetID]; ok {
				context := ""
				for _, xr := range s.XRefs {
					if xr.Source == cn.ID && xr.Target == targetID {
						context = xr.Context
						break
					}
				}
				pairs = append(pairs, xrefPair{
					item: xrefItem{
						Section: target.Section,
						Title:   target.H1Text,
						Context: context,
					},
					isIn: false,
				})
			}
		}
	}

	return pairs
}

func renderXrefMD(cn *spec.ClauseNode, pairs []xrefPair) {
	fmt.Printf("## %s\n\n", cn.H1Text)
	if cn.Section != "" {
		fmt.Printf("Section: %s\n\n", cn.Section)
	}

	var incoming, outgoing []xrefPair
	for _, p := range pairs {
		if p.isIn {
			incoming = append(incoming, p)
		} else {
			outgoing = append(outgoing, p)
		}
	}

	if len(incoming) > 0 {
		fmt.Println("### Referenced by")
		for _, p := range incoming {
			if p.item.Section != "" {
				fmt.Printf("- %s %s\n", p.item.Section, p.item.Title)
			} else {
				fmt.Printf("- %s\n", p.item.Title)
			}
		}
		fmt.Println()
	}

	if len(outgoing) > 0 {
		fmt.Println("### References")
		for _, p := range outgoing {
			if p.item.Section != "" {
				fmt.Printf("- %s %s\n", p.item.Section, p.item.Title)
			} else {
				fmt.Printf("- %s\n", p.item.Title)
			}
		}
		fmt.Println()
	}
}
