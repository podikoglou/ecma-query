package cli

import (
	"fmt"
	"strings"

	"github.com/podikoglou/ecma-query/internal/spec"
	"github.com/spf13/cobra"
)

var tocCmd = &cobra.Command{
	Use:   "toc [SECTION]",
	Short: "Table of contents / structural navigation",
	Long:  "Without arguments: top-level outline. With a section number: children of that section.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runToc,
}

var tocDepth int

func init() {
	tocCmd.Flags().IntVar(&tocDepth, "depth", 2, "number of levels to show")
}

func runToc(_ *cobra.Command, args []string) error {
	s := getSpec()

	if len(args) > 0 && args[0] != "" {
		cn, ok := s.SectionNums[args[0]]
		if !ok {
			exitNotFound(args[0])
			return nil
		}
		entries := buildTocFromNode(cn, tocDepth)
		renderTocMD(entries, 0)
		return nil
	}

	entries := buildTocTopLevel(s, tocDepth)
	renderTocMD(entries, 0)
	return nil
}

type tocEntry struct {
	Section  string
	Title    string
	Children []tocEntry
}

func buildTocTopLevel(s *spec.Spec, maxDepth int) []tocEntry {
	var entries []tocEntry
	for _, cn := range s.Nodes {
		if cn.Parent != nil {
			continue
		}
		entry := tocEntry{
			Section: cn.Section,
			Title:   cn.H1Text,
		}
		if maxDepth > 1 && len(cn.Children) > 0 {
			entry.Children = buildTocFromNode(cn, maxDepth-1)
		}
		entries = append(entries, entry)
	}
	return entries
}

func buildTocFromNode(node *spec.ClauseNode, depth int) []tocEntry {
	if node == nil || depth <= 0 {
		return nil
	}
	var entries []tocEntry
	for _, child := range node.Children {
		entry := tocEntry{
			Section: child.Section,
			Title:   child.H1Text,
		}
		if depth > 1 && len(child.Children) > 0 {
			entry.Children = buildTocFromNode(child, depth-1)
		}
		entries = append(entries, entry)
	}
	return entries
}

func renderTocMD(entries []tocEntry, indent int) {
	prefix := strings.Repeat("  ", indent)
	for _, entry := range entries {
		if entry.Section != "" {
			fmt.Printf("%s- %s %s\n", prefix, entry.Section, entry.Title)
		} else {
			fmt.Printf("%s- %s\n", prefix, entry.Title)
		}
		if len(entry.Children) > 0 {
			renderTocMD(entry.Children, indent+1)
		}
	}
}
