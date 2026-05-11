package citation

import (
	"strings"
)

// APAFormatter formats references in APA (American Psychological Association) style.
// APA 7th edition format for web sources:
// Author, A. A. (Year, Month Day). Title of work. Site Name. URL
type APAFormatter struct{}

// Style returns the citation style.
func (f *APAFormatter) Style() Style {
	return StyleAPA
}

// FormatReference formats a single reference in APA style.
func (f *APAFormatter) FormatReference(ref *Reference) string {
	var parts []string

	// Author (Last, F. M. format)
	if ref.HasAuthors() {
		author := FormatAuthors(ref.Authors, AuthorStyleLastInitial, 20)
		parts = append(parts, author)
	}

	// Date (Year, Month Day) or (Year) or (n.d.)
	date := f.formatDate(ref)
	if date != "" {
		parts = append(parts, "("+date+").")
	} else {
		parts = append(parts, "(n.d.).")
	}

	// Title (sentence case, italicized for standalone works)
	if ref.Title != "" {
		title := f.toSentenceCase(ref.Title)
		switch ref.Type {
		case TypeBlog, TypeArticle:
			// Article titles not italicized
			parts = append(parts, title+".")
		default:
			// Standalone works italicized
			parts = append(parts, "*"+title+"*.")
		}
	}

	// Container/Site name (italicized)
	if ref.Container != "" {
		parts = append(parts, "*"+ref.Container+"*.")
	}

	// Episode for podcasts
	if ref.Episode != "" {
		parts = append(parts, "(No. "+ref.Episode+") [Audio podcast episode].")
	}

	// Publisher (if not same as container)
	if ref.Publisher != "" && ref.Publisher != ref.Container {
		parts = append(parts, ref.Publisher+".")
	}

	// URL (full URL)
	if ref.URL != "" {
		parts = append(parts, ref.URL)
	}

	result := strings.Join(parts, " ")
	// Clean up formatting
	result = strings.ReplaceAll(result, ". .", ".")
	result = strings.ReplaceAll(result, "  ", " ")
	return result
}

// FormatCollection formats all references in a collection.
func (f *APAFormatter) FormatCollection(c *Collection) string {
	// APA references should be sorted alphabetically by author
	sorted := &Collection{References: make([]Reference, len(c.References))}
	copy(sorted.References, c.References)
	sorted.SortByAuthor()

	var lines []string
	for _, ref := range sorted.References {
		lines = append(lines, f.FormatReference(&ref))
	}
	return strings.Join(lines, "\n")
}

// formatDate formats the date in APA style (Year, Month Day).
func (f *APAFormatter) formatDate(ref *Reference) string {
	t := ref.ParsedDate()
	if t.IsZero() {
		return ""
	}

	// Check if we have full date or just year
	if ref.Date == t.Format("2006") {
		return t.Format("2006")
	}
	if ref.Date == t.Format("2006-01") {
		return t.Format("2006, January")
	}
	return t.Format("2006, January 2")
}

// toSentenceCase converts a title to sentence case (only first word and proper nouns capitalized).
// This is a simplified version that only lowercases after the first word.
func (f *APAFormatter) toSentenceCase(title string) string {
	if title == "" {
		return ""
	}

	words := strings.Fields(title)
	if len(words) == 0 {
		return ""
	}

	// Keep first word as-is
	result := []string{words[0]}

	// Lowercase subsequent words (except those that look like proper nouns/acronyms)
	for _, word := range words[1:] {
		// Keep words that are all caps (acronyms) or start with capital after punctuation
		if f.isAcronym(word) || f.isProperNoun(word) {
			result = append(result, word)
		} else {
			result = append(result, strings.ToLower(word))
		}
	}

	return strings.Join(result, " ")
}

// isAcronym checks if a word appears to be an acronym (all caps, 2+ letters).
func (f *APAFormatter) isAcronym(word string) bool {
	if len(word) < 2 {
		return false
	}
	for _, r := range word {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

// isProperNoun makes a simple check for proper nouns.
// This is imperfect but catches common cases.
func (f *APAFormatter) isProperNoun(word string) bool {
	// Words starting with capital that contain lowercase are likely proper nouns
	if len(word) < 2 {
		return false
	}
	if word[0] >= 'A' && word[0] <= 'Z' {
		// Check if it's a product name or proper noun (contains mix of cases)
		hasLower := false
		for _, r := range word[1:] {
			if r >= 'a' && r <= 'z' {
				hasLower = true
				break
			}
		}
		// Heuristic: if starts with cap and has lowercase, might be proper noun
		// but we'll be conservative and lowercase most things
		return hasLower && len(word) > 5
	}
	return false
}
