package citation

import (
	"io"
	"strings"
)

// Style represents a citation/bibliography style.
type Style string

const (
	StyleSimple  Style = "simple"
	StyleMLA     Style = "mla"
	StyleAPA     Style = "apa"
	StyleChicago Style = "chicago"
)

// Formatter formats references in a specific citation style.
type Formatter interface {
	// FormatReference formats a single reference as a string.
	FormatReference(ref *Reference) string

	// FormatCollection formats all references in a collection.
	FormatCollection(c *Collection) string

	// Style returns the citation style this formatter uses.
	Style() Style
}

// NewFormatter creates a formatter for the specified style.
func NewFormatter(style Style) Formatter {
	switch style {
	case StyleMLA:
		return &MLAFormatter{}
	case StyleAPA:
		return &APAFormatter{}
	case StyleChicago:
		return &ChicagoFormatter{}
	case StyleSimple:
		fallthrough
	default:
		return &SimpleFormatter{}
	}
}

// FormatMarkdown formats a collection as a Markdown references section.
func FormatMarkdown(c *Collection, f Formatter, heading string) string {
	var sb strings.Builder
	if heading != "" {
		sb.WriteString("## ")
		sb.WriteString(heading)
		sb.WriteString("\n\n")
	}
	for i, ref := range c.References {
		sb.WriteString("- ")
		sb.WriteString(f.FormatReference(&ref))
		if i < len(c.References)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// FormatMarkdownGrouped formats a collection as a Markdown references section
// grouped by the specified grouping function.
func FormatMarkdownGrouped(c *Collection, f Formatter, heading string, groups map[string][]Reference, groupOrder []string) string {
	var sb strings.Builder
	if heading != "" {
		sb.WriteString("## ")
		sb.WriteString(heading)
		sb.WriteString("\n\n")
	}

	for i, groupName := range groupOrder {
		refs, ok := groups[groupName]
		if !ok || len(refs) == 0 {
			continue
		}

		if groupName != "" {
			sb.WriteString("**")
			sb.WriteString(groupName)
			sb.WriteString("**\n\n")
		}

		for _, ref := range refs {
			sb.WriteString("- ")
			sb.WriteString(f.FormatReference(&ref))
			sb.WriteString("\n")
		}

		if i < len(groupOrder)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// WriteMarkdown writes a formatted collection to a writer.
func WriteMarkdown(w io.Writer, c *Collection, f Formatter, heading string) error {
	_, err := io.WriteString(w, FormatMarkdown(c, f, heading))
	return err
}

// formatURLShort returns a shortened display version of a URL (no protocol).
func formatURLShort(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, "/")
	return url
}
