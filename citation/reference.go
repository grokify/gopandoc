// Package citation provides types and formatters for bibliographic references.
// It supports multiple citation styles (MLA, APA, Chicago, Simple) and can
// load/save reference collections from JSON.
package citation

import (
	"strings"
	"time"
)

// ReferenceType represents the type of source being cited.
type ReferenceType string

const (
	TypeArticle    ReferenceType = "article"
	TypeBlog       ReferenceType = "blog"
	TypeBook       ReferenceType = "book"
	TypePodcast    ReferenceType = "podcast"
	TypeVideo      ReferenceType = "video"
	TypeWebsite    ReferenceType = "website"
	TypeWhitepaper ReferenceType = "whitepaper"
)

// Reference represents a bibliographic reference that can be formatted
// in various citation styles.
type Reference struct {
	// ID is a unique identifier used for citation keys (e.g., "willison-2026")
	ID string `json:"id"`

	// Type indicates the kind of source (blog, podcast, website, etc.)
	Type ReferenceType `json:"type"`

	// Title is the title of the work
	Title string `json:"title"`

	// Authors contains the list of authors or organizations
	Authors []Author `json:"authors,omitempty"`

	// Date is the publication date (YYYY-MM-DD format)
	Date string `json:"date,omitempty"`

	// URL is the web address of the source
	URL string `json:"url,omitempty"`

	// Publisher is the publishing organization
	Publisher string `json:"publisher,omitempty"`

	// Container is the parent work (journal, podcast name, blog name, etc.)
	Container string `json:"container,omitempty"`

	// Episode is the episode number for podcasts
	Episode string `json:"episode,omitempty"`

	// Volume is the volume number for journals
	Volume string `json:"volume,omitempty"`

	// Issue is the issue number for journals
	Issue string `json:"issue,omitempty"`

	// Pages is the page range for articles
	Pages string `json:"pages,omitempty"`

	// Accessed is the date the URL was last accessed (YYYY-MM-DD format)
	Accessed string `json:"accessed,omitempty"`

	// Tags are optional labels for grouping references
	Tags []string `json:"tags,omitempty"`

	// Description is an optional annotation or note
	Description string `json:"description,omitempty"`
}

// ParsedDate returns the Date field parsed as a time.Time.
// Returns zero time if parsing fails.
func (r *Reference) ParsedDate() time.Time {
	formats := []string{
		"2006-01-02",
		"2006-01",
		"2006",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, r.Date); err == nil {
			return t
		}
	}
	return time.Time{}
}

// ParsedAccessed returns the Accessed field parsed as a time.Time.
// Returns zero time if parsing fails.
func (r *Reference) ParsedAccessed() time.Time {
	if r.Accessed == "" {
		return time.Time{}
	}
	formats := []string{
		"2006-01-02",
		"2006-01",
		"2006",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, r.Accessed); err == nil {
			return t
		}
	}
	return time.Time{}
}

// Year returns the publication year as a string.
func (r *Reference) Year() string {
	t := r.ParsedDate()
	if t.IsZero() {
		return ""
	}
	return t.Format("2006")
}

// HasAuthors returns true if the reference has at least one author.
func (r *Reference) HasAuthors() bool {
	return len(r.Authors) > 0
}

// HasPersonAuthors returns true if the reference has at least one person author
// (as opposed to only organizational authors).
func (r *Reference) HasPersonAuthors() bool {
	for _, a := range r.Authors {
		if a.IsPerson() {
			return true
		}
	}
	return false
}

// PrimaryAuthor returns the first author, or an empty Author if none exist.
func (r *Reference) PrimaryAuthor() Author {
	if len(r.Authors) == 0 {
		return Author{}
	}
	return r.Authors[0]
}

// URLHost returns the hostname portion of the URL.
func (r *Reference) URLHost() string {
	url := r.URL
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "www.")
	if idx := strings.Index(url, "/"); idx > 0 {
		url = url[:idx]
	}
	return url
}

// HasTag returns true if the reference has the specified tag.
func (r *Reference) HasTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, t := range r.Tags {
		if strings.ToLower(strings.TrimSpace(t)) == tag {
			return true
		}
	}
	return false
}

// Validate checks that required fields are present and returns an error if not.
func (r *Reference) Validate() error {
	var missing []string
	if r.ID == "" {
		missing = append(missing, "id")
	}
	if r.Type == "" {
		missing = append(missing, "type")
	}
	if r.Title == "" {
		missing = append(missing, "title")
	}
	if len(missing) > 0 {
		return &ValidationError{
			ReferenceID: r.ID,
			Missing:     missing,
		}
	}
	return nil
}

// ValidationError is returned when a reference is missing required fields.
type ValidationError struct {
	ReferenceID string
	Missing     []string
}

func (e *ValidationError) Error() string {
	return "reference " + e.ReferenceID + " missing required fields: " + strings.Join(e.Missing, ", ")
}
