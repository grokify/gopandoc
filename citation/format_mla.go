package citation

import (
	"strings"
)

// MLAFormatter formats references in MLA (Modern Language Association) style.
// MLA 9th edition format for web sources:
// Author. "Title." *Container*, Publisher, Day Month Year, URL.
type MLAFormatter struct{}

// Style returns the citation style.
func (f *MLAFormatter) Style() Style {
	return StyleMLA
}

// FormatReference formats a single reference in MLA style.
func (f *MLAFormatter) FormatReference(ref *Reference) string {
	var parts []string

	// Author (Last, First format for first author)
	if ref.HasAuthors() {
		author := FormatAuthors(ref.Authors, AuthorStyleLastFirst, 3)
		parts = append(parts, author+".")
	}

	// Title (in quotes for shorter works)
	if ref.Title != "" {
		switch ref.Type {
		case TypeBook:
			// Books are italicized, not quoted
			parts = append(parts, "*"+ref.Title+"*.")
		default:
			parts = append(parts, `"`+ref.Title+`."`)
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

	// Publisher (if different from container)
	if ref.Publisher != "" && ref.Publisher != ref.Container {
		parts = append(parts, ref.Publisher+",")
	}

	// Volume and issue for journals
	if ref.Volume != "" {
		volIssue := "vol. " + ref.Volume
		if ref.Issue != "" {
			volIssue += ", no. " + ref.Issue
		}
		parts = append(parts, volIssue+",")
	}

	// Date (Day Month Year format)
	if date := f.formatDate(ref); date != "" {
		parts = append(parts, date+",")
	}

	// Pages
	if ref.Pages != "" {
		parts = append(parts, "pp. "+ref.Pages+",")
	}

	// URL (shortened, no protocol)
	if ref.URL != "" {
		parts = append(parts, formatURLShort(ref.URL)+".")
	}

	// Accessed date
	if accessed := f.formatAccessed(ref); accessed != "" {
		// Remove the trailing period from URL and add accessed
		result := strings.Join(parts, " ")
		if strings.HasSuffix(result, ".") {
			result = result[:len(result)-1]
		}
		return result + ". Accessed " + accessed + "."
	}

	result := strings.Join(parts, " ")
	// Clean up any double periods or trailing commas
	result = strings.ReplaceAll(result, ",.", ".")
	result = strings.ReplaceAll(result, ",,", ",")
	if strings.HasSuffix(result, ",") {
		result = result[:len(result)-1] + "."
	}
	return result
}

// FormatCollection formats all references in a collection.
func (f *MLAFormatter) FormatCollection(c *Collection) string {
	// MLA works cited should be sorted alphabetically by author
	sorted := &Collection{References: make([]Reference, len(c.References))}
	copy(sorted.References, c.References)
	sorted.SortByAuthor()

	var lines []string
	for _, ref := range sorted.References {
		lines = append(lines, f.FormatReference(&ref))
	}
	return strings.Join(lines, "\n")
}

// formatDate formats the date in MLA style (Day Month Year).
func (f *MLAFormatter) formatDate(ref *Reference) string {
	t := ref.ParsedDate()
	if t.IsZero() {
		return ""
	}
	// MLA uses abbreviated months except May, June, July
	month := t.Format("Jan.")
	switch t.Month() {
	case 5:
		month = "May"
	case 6:
		month = "June"
	case 7:
		month = "July"
	case 9:
		month = "Sept."
	}
	return t.Format("2") + " " + month + " " + t.Format("2006")
}

// formatAccessed formats the accessed date in MLA style.
func (f *MLAFormatter) formatAccessed(ref *Reference) string {
	t := ref.ParsedAccessed()
	if t.IsZero() {
		return ""
	}
	month := t.Format("Jan.")
	switch t.Month() {
	case 5:
		month = "May"
	case 6:
		month = "June"
	case 7:
		month = "July"
	case 9:
		month = "Sept."
	}
	return t.Format("2") + " " + month + " " + t.Format("2006")
}
