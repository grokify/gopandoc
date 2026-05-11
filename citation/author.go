package citation

import "strings"

// Author represents an individual or organizational author.
type Author struct {
	// Given is the first/given name (e.g., "Simon")
	Given string `json:"given,omitempty"`

	// Family is the last/family name (e.g., "Willison")
	Family string `json:"family,omitempty"`

	// Org is the organization name for corporate authors (e.g., "AWS")
	Org string `json:"org,omitempty"`

	// Suffix is a name suffix (e.g., "Jr.", "III")
	Suffix string `json:"suffix,omitempty"`
}

// IsPerson returns true if this is a person author (has given or family name).
func (a Author) IsPerson() bool {
	return a.Given != "" || a.Family != ""
}

// IsOrg returns true if this is an organizational author.
func (a Author) IsOrg() bool {
	return a.Org != "" && !a.IsPerson()
}

// IsEmpty returns true if no author information is present.
func (a Author) IsEmpty() bool {
	return a.Given == "" && a.Family == "" && a.Org == ""
}

// FullName returns the author's full name in natural order (Given Family).
// For organizations, returns the organization name.
func (a Author) FullName() string {
	if a.IsOrg() {
		return a.Org
	}
	parts := []string{}
	if a.Given != "" {
		parts = append(parts, a.Given)
	}
	if a.Family != "" {
		parts = append(parts, a.Family)
	}
	if a.Suffix != "" {
		parts = append(parts, a.Suffix)
	}
	return strings.Join(parts, " ")
}

// LastFirst returns the author's name in "Family, Given" format.
// For organizations, returns the organization name.
func (a Author) LastFirst() string {
	if a.IsOrg() {
		return a.Org
	}
	if a.Family == "" {
		return a.Given
	}
	if a.Given == "" {
		return a.Family
	}
	result := a.Family + ", " + a.Given
	if a.Suffix != "" {
		result += ", " + a.Suffix
	}
	return result
}

// LastInitial returns the author's name in "Family, G." format.
// For organizations, returns the organization name.
func (a Author) LastInitial() string {
	if a.IsOrg() {
		return a.Org
	}
	if a.Family == "" {
		if a.Given != "" {
			return string(a.Given[0]) + "."
		}
		return ""
	}
	if a.Given == "" {
		return a.Family
	}
	result := a.Family + ", " + string(a.Given[0]) + "."
	if a.Suffix != "" {
		result += ", " + a.Suffix
	}
	return result
}

// InitialLast returns the author's name in "G. Family" format.
// For organizations, returns the organization name.
func (a Author) InitialLast() string {
	if a.IsOrg() {
		return a.Org
	}
	if a.Given == "" {
		return a.Family
	}
	if a.Family == "" {
		return string(a.Given[0]) + "."
	}
	result := string(a.Given[0]) + ". " + a.Family
	if a.Suffix != "" {
		result += ", " + a.Suffix
	}
	return result
}

// FormatAuthors formats a slice of authors according to the specified style.
// maxAuthors limits how many authors to show before "et al." (0 = no limit).
func FormatAuthors(authors []Author, style AuthorStyle, maxAuthors int) string {
	if len(authors) == 0 {
		return ""
	}

	var formatted []string
	for i, a := range authors {
		if maxAuthors > 0 && i >= maxAuthors {
			break
		}
		var name string
		switch style {
		case AuthorStyleLastFirst:
			if i == 0 {
				name = a.LastFirst()
			} else {
				name = a.FullName() // Subsequent authors in normal order
			}
		case AuthorStyleLastInitial:
			name = a.LastInitial()
		case AuthorStyleInitialLast:
			name = a.InitialLast()
		case AuthorStyleFull:
			name = a.FullName()
		default:
			name = a.FullName()
		}
		formatted = append(formatted, name)
	}

	result := joinAuthors(formatted)
	if maxAuthors > 0 && len(authors) > maxAuthors {
		result += ", et al."
	}
	return result
}

// joinAuthors joins author names with commas and "and" before the last.
func joinAuthors(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	case 2:
		return names[0] + " and " + names[1]
	default:
		return strings.Join(names[:len(names)-1], ", ") + ", and " + names[len(names)-1]
	}
}

// AuthorStyle specifies how to format author names.
type AuthorStyle int

const (
	// AuthorStyleFull formats as "Given Family"
	AuthorStyleFull AuthorStyle = iota
	// AuthorStyleLastFirst formats as "Family, Given"
	AuthorStyleLastFirst
	// AuthorStyleLastInitial formats as "Family, G."
	AuthorStyleLastInitial
	// AuthorStyleInitialLast formats as "G. Family"
	AuthorStyleInitialLast
)
