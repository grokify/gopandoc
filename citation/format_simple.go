package citation

import "strings"

// SimpleFormatter formats references in a simple, readable style.
// Example: Willison, Simon. "The Software Factory." February 2026. https://simonwillison.net/...
type SimpleFormatter struct{}

// Style returns the citation style.
func (f *SimpleFormatter) Style() Style {
	return StyleSimple
}

// FormatReference formats a single reference.
func (f *SimpleFormatter) FormatReference(ref *Reference) string {
	var parts []string

	// Author
	if ref.HasAuthors() {
		author := FormatAuthors(ref.Authors, AuthorStyleLastFirst, 3)
		parts = append(parts, author+".")
	}

	// Title (in quotes for articles/blogs, plain for websites)
	if ref.Title != "" {
		switch ref.Type {
		case TypeWebsite:
			parts = append(parts, ref.Title+".")
		default:
			parts = append(parts, `"`+ref.Title+`."`)
		}
	}

	// Container (italicized in markdown)
	if ref.Container != "" {
		parts = append(parts, "*"+ref.Container+"*.")
	}

	// Episode for podcasts
	if ref.Episode != "" {
		parts = append(parts, "Episode "+ref.Episode+".")
	}

	// Publisher
	if ref.Publisher != "" && ref.Publisher != ref.Container {
		parts = append(parts, ref.Publisher+".")
	}

	// Date
	if date := f.formatDate(ref); date != "" {
		parts = append(parts, date+".")
	}

	// URL
	if ref.URL != "" {
		parts = append(parts, formatURLShort(ref.URL))
	}

	return strings.Join(parts, " ")
}

// FormatCollection formats all references in a collection.
func (f *SimpleFormatter) FormatCollection(c *Collection) string {
	var lines []string
	for _, ref := range c.References {
		lines = append(lines, f.FormatReference(&ref))
	}
	return strings.Join(lines, "\n")
}

// formatDate formats the date for display.
func (f *SimpleFormatter) formatDate(ref *Reference) string {
	t := ref.ParsedDate()
	if t.IsZero() {
		return ""
	}
	return t.Format("January 2006")
}
