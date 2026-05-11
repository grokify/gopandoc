// citegen generates formatted reference lists from JSON files.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/grokify/gopandoc/citation"
	"github.com/spf13/cobra"
)

var (
	inputFile string
	style     string
	heading   string
	tag       string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "citegen",
	Short: "Generate formatted reference lists from JSON files",
	Long: `citegen generates formatted reference lists from JSON files.

It supports multiple citation styles (simple, MLA, APA, Chicago) and can
filter references by tag. Output is Markdown-formatted.

Examples:
  citegen generate -i references.json -s simple
  citegen generate -i references.json -s mla --heading "Works Cited"
  citegen generate -i references.json -s apa --tag aws
  citegen validate -i references.json
  citegen tags -i references.json`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate formatted references",
	Long:  `Generate a formatted reference list in the specified citation style.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		coll, err := citation.LoadJSON(inputFile)
		if err != nil {
			return fmt.Errorf("loading %s: %w", inputFile, err)
		}

		// Filter by tag if specified
		if tag != "" {
			refs := coll.GetByTag(tag)
			coll = citation.NewCollection()
			for _, r := range refs {
				coll.Add(r)
			}
		}

		// Create formatter
		var s citation.Style
		switch strings.ToLower(style) {
		case "mla":
			s = citation.StyleMLA
		case "apa":
			s = citation.StyleAPA
		case "chicago":
			s = citation.StyleChicago
		default:
			s = citation.StyleSimple
		}
		formatter := citation.NewFormatter(s)

		// Generate output
		output := citation.FormatMarkdown(coll, formatter, heading)
		fmt.Println(output)
		return nil
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate references JSON file",
	Long:  `Validate that all references have required fields and no duplicate IDs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		coll, err := citation.LoadJSON(inputFile)
		if err != nil {
			return fmt.Errorf("loading %s: %w", inputFile, err)
		}

		errs := coll.Validate()
		if len(errs) > 0 {
			fmt.Fprintln(os.Stderr, "Validation errors:")
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  - %v\n", e)
			}
			return fmt.Errorf("validation failed with %d errors", len(errs))
		}
		fmt.Printf("Validated %d references, no errors.\n", coll.Len())
		return nil
	},
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags in the collection",
	Long:  `List all unique tags used in the reference collection with counts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		coll, err := citation.LoadJSON(inputFile)
		if err != nil {
			return fmt.Errorf("loading %s: %w", inputFile, err)
		}

		tags := coll.Tags()
		fmt.Println("Tags:")
		for _, t := range tags {
			refs := coll.GetByTag(t)
			fmt.Printf("  %s (%d)\n", t, len(refs))
		}
		return nil
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Input JSON file containing references (required)")
	_ = rootCmd.MarkPersistentFlagRequired("input")

	// Generate command flags
	generateCmd.Flags().StringVarP(&style, "style", "s", "simple", "Citation style: simple, mla, apa, chicago")
	generateCmd.Flags().StringVar(&heading, "heading", "References", "Section heading (empty string for no heading)")
	generateCmd.Flags().StringVarP(&tag, "tag", "t", "", "Filter references by tag")

	// Add subcommands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(tagsCmd)
}
