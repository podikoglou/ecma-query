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

	var root *spec.ClauseNode
	if len(args) > 0 && args[0] != "" {
		var err error
		root, err = resolveTocRoot(s, args[0])
		if err != nil {
			ExitNotFound(args[0])
			return nil
		}
	} else {
		root = s.Root
	}

	entries := buildToc(root, 0, tocDepth)

	if format == FormatJSON {
		var rootSection *string
		if root != nil && root.Section != "" {
			rootSection = &root.Section
		}
		printJSON(TocResponse{
			Root:    rootSection,
			Entries: entries,
		})
	} else {
		renderTocMD(entries, 0)
	}
	return nil
}

func resolveTocRoot(s *spec.Spec, query string) (*spec.ClauseNode, error) {
	if cn, ok := s.SectionNums[query]; ok {
		return cn, nil
	}
	return nil, fmt.Errorf("section not found: %s", query)
}

func buildToc(node *spec.ClauseNode, depth, maxDepth int) []TocEntry {
	if node == nil || depth >= maxDepth {
		return nil
	}
	var entries []TocEntry
	for _, child := range node.Children {
		entry := TocEntry{
			Section: child.Section,
			Title:   child.H1Text,
		}
		if depth+1 < maxDepth && len(child.Children) > 0 {
			entry.Children = buildToc(child, depth+1, maxDepth)
		}
		entries = append(entries, entry)
	}
	return entries
}

func renderTocMD(entries []TocEntry, indent int) {
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
