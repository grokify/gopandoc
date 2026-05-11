package citation

import (
	"strings"
)

// ChicagoFormatter formats references in Chicago Manual of Style (17th ed.) format.
// Chicago notes-bibliography style for web sources:
// Author First Last, "Title," Site Name, Month Day, Year, URL.
type ChicagoFormatter struct {
	// UseNotesStyle uses notes style (footnote format) instead of bibliography style
	UseNotesStyle bool
}

// Style returns the citation style.
func (f *ChicagoFormatter) Style() Style {
	return StyleChicago
}

// FormatReference formats a single reference in Chicago style.
func (f *ChicagoFormatter) FormatReference(ref *Reference) string {
	if f.UseNotesStyle {
		return f.formatNotes(ref)
	}
	return f.formatBibliography(ref)
}

// formatBibliography formats a reference in Chicago bibliography style.
func (f *ChicagoFormatter) formatBibliography(ref *Reference) string {
	var parts []string

	// Author (Last, First for bibliography)
	if ref.HasAuthors() {
		author := FormatAuthors(ref.Authors, AuthorStyleLastFirst, 3)
		parts = append(parts, author+".")
	}

	// Title (in quotes for articles, italicized for books/websites)
	if ref.Title != "" {
		switch ref.Type {
		case TypeBook, TypeWebsite:
			parts = append(parts, "*"+ref.Title+"*.")
		default:
			parts = append(parts, `"`+ref.Title+`."`)
		}
	}

	// Container (italicized)
	if ref.Container != "" {
		parts = append(parts, "*"+ref.Container+"*.")
	}

	// Episode for podcasts
	if ref.Episode != "" {
		parts = append(parts, "Episode "+ref.Episode+".")
	}

	// Publisher
	if ref.Publisher != "" && ref.Publisher != ref.Container {
		parts = append(parts, ref.Publisher+",")
	}

	// Date
	if date := f.formatDateBibliography(ref); date != "" {
		parts = append(parts, date+".")
	}

	// URL
	if ref.URL != "" {
		parts = append(parts, ref.URL+".")
	}

	result := strings.Join(parts, " ")
	result = strings.ReplaceAll(result, ",.", ".")
	result = strings.ReplaceAll(result, "..", ".")
	return result
}

// formatNotes formats a reference in Chicago notes style (for footnotes).
func (f *ChicagoFormatter) formatNotes(ref *Reference) string {
	var parts []string

	// Author (First Last for notes)
	if ref.HasAuthors() {
		author := FormatAuthors(ref.Authors, AuthorStyleFull, 3)
		parts = append(parts, author+",")
	}

	// Title (in quotes for articles, italicized for books/websites)
	if ref.Title != "" {
		switch ref.Type {
		case TypeBook, TypeWebsite:
			parts = append(parts, "*"+ref.Title+"*")
		default:
			parts = append(parts, `"`+ref.Title+`,"`)
		}
	}

	// Container (italicized)
	if ref.Container != "" {
		parts = append(parts, "*"+ref.Container+"*,")
	}

	// Episode for podcasts
	if ref.Episode != "" {
		parts = append(parts, "episode "+ref.Episode+",")
	}

	// Publisher
	if ref.Publisher != "" && ref.Publisher != ref.Container {
		parts = append(parts, ref.Publisher+",")
	}

	// Date
	if date := f.formatDateNotes(ref); date != "" {
		parts = append(parts, date+",")
	}

	// URL
	if ref.URL != "" {
		parts = append(parts, ref.URL)
	}

	result := strings.Join(parts, " ")
	// Clean up
	result = strings.ReplaceAll(result, ",,", ",")
	if strings.HasSuffix(result, ",") {
		result = result[:len(result)-1]
	}
	if !strings.HasSuffix(result, ".") {
		result += "."
	}
	return result
}

// FormatCollection formats all references in a collection.
func (f *ChicagoFormatter) FormatCollection(c *Collection) string {
	// Chicago bibliography should be sorted alphabetically by author
	sorted := &Collection{References: make([]Reference, len(c.References))}
	copy(sorted.References, c.References)
	sorted.SortByAuthor()

	var lines []string
	for _, ref := range sorted.References {
		lines = append(lines, f.FormatReference(&ref))
	}
	return strings.Join(lines, "\n")
}

// formatDateBibliography formats the date for bibliography style (Month Day, Year).
func (f *ChicagoFormatter) formatDateBibliography(ref *Reference) string {
	t := ref.ParsedDate()
	if t.IsZero() {
		return ""
	}
	// Check date precision
	if ref.Date == t.Format("2006") {
		return t.Format("2006")
	}
	if ref.Date == t.Format("2006-01") {
		return t.Format("January 2006")
	}
	return t.Format("January 2, 2006")
}

// formatDateNotes formats the date for notes style (Month Day, Year).
func (f *ChicagoFormatter) formatDateNotes(ref *Reference) string {
	t := ref.ParsedDate()
	if t.IsZero() {
		return ""
	}
	// Check date precision
	if ref.Date == t.Format("2006") {
		return t.Format("2006")
	}
	if ref.Date == t.Format("2006-01") {
		return t.Format("January 2006")
	}
	return t.Format("January 2, 2006")
}
