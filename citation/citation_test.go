package citation

import (
	"strings"
	"testing"
)

func TestAuthorFullName(t *testing.T) {
	tests := []struct {
		author Author
		want   string
	}{
		{Author{Given: "Simon", Family: "Willison"}, "Simon Willison"},
		{Author{Given: "Joe", Family: "Magherimov"}, "Joe Magherimov"},
		{Author{Org: "AWS"}, "AWS"},
		{Author{Given: "Martin", Family: "King", Suffix: "Jr."}, "Martin King Jr."},
		{Author{}, ""},
	}

	for _, tt := range tests {
		got := tt.author.FullName()
		if got != tt.want {
			t.Errorf("FullName() = %q, want %q", got, tt.want)
		}
	}
}

func TestAuthorLastFirst(t *testing.T) {
	tests := []struct {
		author Author
		want   string
	}{
		{Author{Given: "Simon", Family: "Willison"}, "Willison, Simon"},
		{Author{Org: "AWS"}, "AWS"},
		{Author{Given: "Simon"}, "Simon"},
		{Author{Family: "Willison"}, "Willison"},
	}

	for _, tt := range tests {
		got := tt.author.LastFirst()
		if got != tt.want {
			t.Errorf("LastFirst() = %q, want %q", got, tt.want)
		}
	}
}

func TestAuthorLastInitial(t *testing.T) {
	tests := []struct {
		author Author
		want   string
	}{
		{Author{Given: "Simon", Family: "Willison"}, "Willison, S."},
		{Author{Given: "Joe", Family: "Magherimov"}, "Magherimov, J."},
		{Author{Org: "AWS"}, "AWS"},
	}

	for _, tt := range tests {
		got := tt.author.LastInitial()
		if got != tt.want {
			t.Errorf("LastInitial() = %q, want %q", got, tt.want)
		}
	}
}

func TestFormatAuthors(t *testing.T) {
	authors := []Author{
		{Given: "Simon", Family: "Willison"},
		{Given: "Joe", Family: "Magherimov"},
		{Given: "Dan", Family: "Shapiro"},
	}

	got := FormatAuthors(authors, AuthorStyleLastFirst, 0)
	want := "Willison, Simon, Joe Magherimov, and Dan Shapiro"
	if got != want {
		t.Errorf("FormatAuthors() = %q, want %q", got, want)
	}

	// Test with max authors
	got = FormatAuthors(authors, AuthorStyleLastFirst, 2)
	want = "Willison, Simon and Joe Magherimov, et al."
	if got != want {
		t.Errorf("FormatAuthors(max=2) = %q, want %q", got, want)
	}
}

func TestReferenceYear(t *testing.T) {
	ref := Reference{Date: "2026-02-07"}
	if got := ref.Year(); got != "2026" {
		t.Errorf("Year() = %q, want %q", got, "2026")
	}

	ref = Reference{Date: "2026"}
	if got := ref.Year(); got != "2026" {
		t.Errorf("Year() = %q, want %q", got, "2026")
	}
}

func TestReferenceURLHost(t *testing.T) {
	ref := Reference{URL: "https://simonwillison.net/2026/Feb/7/software-factory/"}
	if got := ref.URLHost(); got != "simonwillison.net" {
		t.Errorf("URLHost() = %q, want %q", got, "simonwillison.net")
	}
}

func TestSimpleFormatter(t *testing.T) {
	ref := &Reference{
		ID:      "willison-2026",
		Type:    TypeBlog,
		Title:   "The Software Factory",
		Authors: []Author{{Given: "Simon", Family: "Willison"}},
		Date:    "2026-02-07",
		URL:     "https://simonwillison.net/2026/Feb/7/software-factory/",
	}

	f := &SimpleFormatter{}
	got := f.FormatReference(ref)

	// Check key components are present
	if !strings.Contains(got, "Willison, Simon") {
		t.Errorf("SimpleFormatter missing author: %s", got)
	}
	if !strings.Contains(got, `"The Software Factory."`) {
		t.Errorf("SimpleFormatter missing title: %s", got)
	}
	if !strings.Contains(got, "February 2026") {
		t.Errorf("SimpleFormatter missing date: %s", got)
	}
}

func TestMLAFormatter(t *testing.T) {
	ref := &Reference{
		ID:        "willison-2026",
		Type:      TypeBlog,
		Title:     "The Software Factory",
		Authors:   []Author{{Given: "Simon", Family: "Willison"}},
		Date:      "2026-02-07",
		Container: "Simon Willison's Weblog",
		URL:       "https://simonwillison.net/2026/Feb/7/software-factory/",
	}

	f := &MLAFormatter{}
	got := f.FormatReference(ref)

	// Check key components are present
	if !strings.Contains(got, "Willison, Simon") {
		t.Errorf("MLAFormatter missing author: %s", got)
	}
	if !strings.Contains(got, "*Simon Willison's Weblog*") {
		t.Errorf("MLAFormatter missing container: %s", got)
	}
	if !strings.Contains(got, "7 Feb. 2026") {
		t.Errorf("MLAFormatter missing date: %s", got)
	}
}

func TestAPAFormatter(t *testing.T) {
	ref := &Reference{
		ID:        "willison-2026",
		Type:      TypeBlog,
		Title:     "The Software Factory",
		Authors:   []Author{{Given: "Simon", Family: "Willison"}},
		Date:      "2026-02-07",
		Container: "Simon Willison's Weblog",
		URL:       "https://simonwillison.net/2026/Feb/7/software-factory/",
	}

	f := &APAFormatter{}
	got := f.FormatReference(ref)

	// Check key components are present
	if !strings.Contains(got, "Willison, S.") {
		t.Errorf("APAFormatter missing author: %s", got)
	}
	if !strings.Contains(got, "(2026, February 7)") {
		t.Errorf("APAFormatter missing date: %s", got)
	}
}

func TestCollectionValidation(t *testing.T) {
	c := NewCollection()
	c.Add(Reference{ID: "test1", Type: TypeBlog, Title: "Test"})
	c.Add(Reference{ID: "test1", Type: TypeBlog, Title: "Duplicate"}) // duplicate ID
	c.Add(Reference{Type: TypeBlog, Title: "Missing ID"})             // missing ID

	errs := c.Validate()
	if len(errs) != 2 {
		t.Errorf("Validate() returned %d errors, want 2", len(errs))
	}
}

func TestCollectionGetByTag(t *testing.T) {
	c := NewCollection()
	c.Add(Reference{ID: "r1", Type: TypeBlog, Title: "T1", Tags: []string{"aws", "ai"}})
	c.Add(Reference{ID: "r2", Type: TypeBlog, Title: "T2", Tags: []string{"ai"}})
	c.Add(Reference{ID: "r3", Type: TypeBlog, Title: "T3", Tags: []string{"other"}})

	got := c.GetByTag("ai")
	if len(got) != 2 {
		t.Errorf("GetByTag(ai) returned %d refs, want 2", len(got))
	}

	got = c.GetByTag("AWS") // case insensitive
	if len(got) != 1 {
		t.Errorf("GetByTag(AWS) returned %d refs, want 1", len(got))
	}
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		style Style
		want  Style
	}{
		{StyleSimple, StyleSimple},
		{StyleMLA, StyleMLA},
		{StyleAPA, StyleAPA},
		{StyleChicago, StyleChicago},
		{"unknown", StyleSimple}, // default
	}

	for _, tt := range tests {
		f := NewFormatter(tt.style)
		if got := f.Style(); got != tt.want {
			t.Errorf("NewFormatter(%q).Style() = %q, want %q", tt.style, got, tt.want)
		}
	}
}

func TestFormatMarkdown(t *testing.T) {
	c := NewCollection()
	c.Add(Reference{
		ID:      "test1",
		Type:    TypeBlog,
		Title:   "Test Article",
		Authors: []Author{{Given: "John", Family: "Doe"}},
		Date:    "2026-01-15",
	})

	f := NewFormatter(StyleSimple)
	got := FormatMarkdown(c, f, "References")

	if !strings.HasPrefix(got, "## References") {
		t.Errorf("FormatMarkdown missing heading: %s", got)
	}
	if !strings.Contains(got, "- Doe, John") {
		t.Errorf("FormatMarkdown missing list item: %s", got)
	}
}
