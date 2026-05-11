package citation

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Collection holds a set of references with methods for loading, saving,
// and querying.
type Collection struct {
	References []Reference `json:"references"`
}

// NewCollection creates an empty Collection.
func NewCollection() *Collection {
	return &Collection{
		References: []Reference{},
	}
}

// LoadJSON loads a Collection from a JSON file.
func LoadJSON(filename string) (*Collection, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return ReadJSON(f)
}

// ReadJSON reads a Collection from a JSON reader.
func ReadJSON(r io.Reader) (*Collection, error) {
	var c Collection
	dec := json.NewDecoder(r)
	if err := dec.Decode(&c); err != nil {
		return nil, fmt.Errorf("decoding JSON: %w", err)
	}
	return &c, nil
}

// SaveJSON saves the Collection to a JSON file with indentation.
func (c *Collection) SaveJSON(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	return c.WriteJSON(f)
}

// WriteJSON writes the Collection as indented JSON to a writer.
func (c *Collection) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	return nil
}

// Add adds a reference to the collection.
func (c *Collection) Add(ref Reference) {
	c.References = append(c.References, ref)
}

// Get returns the reference with the given ID, or nil if not found.
func (c *Collection) Get(id string) *Reference {
	for i := range c.References {
		if c.References[i].ID == id {
			return &c.References[i]
		}
	}
	return nil
}

// GetByTag returns all references with the specified tag.
func (c *Collection) GetByTag(tag string) []Reference {
	var result []Reference
	for _, ref := range c.References {
		if ref.HasTag(tag) {
			result = append(result, ref)
		}
	}
	return result
}

// GetByType returns all references of the specified type.
func (c *Collection) GetByType(refType ReferenceType) []Reference {
	var result []Reference
	for _, ref := range c.References {
		if ref.Type == refType {
			result = append(result, ref)
		}
	}
	return result
}

// Len returns the number of references in the collection.
func (c *Collection) Len() int {
	return len(c.References)
}

// Validate checks all references and returns any validation errors.
func (c *Collection) Validate() []error {
	var errs []error
	seen := make(map[string]bool)
	for _, ref := range c.References {
		if err := ref.Validate(); err != nil {
			errs = append(errs, err)
		}
		if seen[ref.ID] {
			errs = append(errs, fmt.Errorf("duplicate reference ID: %s", ref.ID))
		}
		seen[ref.ID] = true
	}
	return errs
}

// SortByID sorts references alphabetically by ID.
func (c *Collection) SortByID() {
	sort.Slice(c.References, func(i, j int) bool {
		return c.References[i].ID < c.References[j].ID
	})
}

// SortByDate sorts references by date (newest first).
func (c *Collection) SortByDate() {
	sort.Slice(c.References, func(i, j int) bool {
		ti := c.References[i].ParsedDate()
		tj := c.References[j].ParsedDate()
		return ti.After(tj)
	})
}

// SortByAuthor sorts references alphabetically by primary author's last name.
func (c *Collection) SortByAuthor() {
	sort.Slice(c.References, func(i, j int) bool {
		ai := c.References[i].PrimaryAuthor()
		aj := c.References[j].PrimaryAuthor()
		// Compare by family name first, then given name
		if ai.Family != aj.Family {
			return ai.Family < aj.Family
		}
		if ai.Given != aj.Given {
			return ai.Given < aj.Given
		}
		// Fall back to org name
		return ai.Org < aj.Org
	})
}

// GroupByTag groups references by their tags.
// References with multiple tags appear in multiple groups.
// References with no tags appear under the empty string key.
func (c *Collection) GroupByTag() map[string][]Reference {
	groups := make(map[string][]Reference)
	for _, ref := range c.References {
		if len(ref.Tags) == 0 {
			groups[""] = append(groups[""], ref)
		} else {
			for _, tag := range ref.Tags {
				tag = strings.ToLower(strings.TrimSpace(tag))
				groups[tag] = append(groups[tag], ref)
			}
		}
	}
	return groups
}

// GroupByType groups references by their type.
func (c *Collection) GroupByType() map[ReferenceType][]Reference {
	groups := make(map[ReferenceType][]Reference)
	for _, ref := range c.References {
		groups[ref.Type] = append(groups[ref.Type], ref)
	}
	return groups
}

// IDs returns a slice of all reference IDs in the collection.
func (c *Collection) IDs() []string {
	ids := make([]string, len(c.References))
	for i, ref := range c.References {
		ids[i] = ref.ID
	}
	return ids
}

// Tags returns a sorted slice of all unique tags in the collection.
func (c *Collection) Tags() []string {
	tagSet := make(map[string]bool)
	for _, ref := range c.References {
		for _, tag := range ref.Tags {
			tag = strings.ToLower(strings.TrimSpace(tag))
			tagSet[tag] = true
		}
	}
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
