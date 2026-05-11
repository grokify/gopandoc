# Release Notes: gopandoc v0.4.0

**Release Date:** 2026-05-11

This release introduces a new `citation` package for bibliographic reference management with support for multiple citation styles, along with the `citegen` CLI tool for generating formatted reference lists from JSON.

## Highlights

- **New `citation` package** — Manage bibliographic references with types for authors, references, and collections
- **Four citation styles** — Simple, MLA (9th ed.), APA (7th ed.), and Chicago Manual of Style
- **`citegen` CLI tool** — Generate formatted references from JSON with filtering and validation
- **JSON-based workflow** — Define references once in JSON, output in any style

## New Features

### Citation Package

The `citation` package provides a complete solution for managing bibliographic references in Go applications.

#### Core Types

```go
import "github.com/grokify/gopandoc/citation"

// Reference represents a bibliographic entry
ref := citation.Reference{
    ID:        "willison-2026",
    Type:      citation.TypeBlog,
    Title:     "The Software Factory",
    Authors:   []citation.Author{{Given: "Simon", Family: "Willison"}},
    Date:      "2026-02-07",
    URL:       "https://simonwillison.net/2026/Feb/7/software-factory/",
    Container: "Simon Willison's Weblog",
    Tags:      []string{"ai", "software-factory"},
}

// Author with flexible name handling
author := citation.Author{
    Given:  "Simon",
    Family: "Willison",
}
author.FullName()    // "Simon Willison"
author.LastFirst()   // "Willison, Simon"
author.LastInitial() // "Willison, S."
```

#### Collection Management

```go
// Load from JSON
coll, err := citation.LoadJSON("references.json")

// Filter by tag
awsRefs := coll.GetByTag("aws")

// Filter by type
blogs := coll.GetByType(citation.TypeBlog)

// Validate all references
if errs := coll.Validate(); len(errs) > 0 {
    // Handle validation errors
}

// Sort options
coll.SortByAuthor()
coll.SortByDate()
coll.SortByID()

// Save back to JSON
coll.SaveJSON("references.json")
```

#### Citation Formatters

Four built-in formatters produce properly formatted citations:

| Style | Output Example |
|-------|----------------|
| **Simple** | Willison, Simon. "The Software Factory." *Simon Willison's Weblog*. February 2026. simonwillison.net/... |
| **MLA** | Willison, Simon. "The Software Factory." *Simon Willison's Weblog*, 7 Feb. 2026, simonwillison.net/... |
| **APA** | Willison, S. (2026, February 7). The software factory. *Simon Willison's Weblog*. https://simonwillison.net/... |
| **Chicago** | Willison, Simon. "The Software Factory." *Simon Willison's Weblog*, February 7, 2026. https://simonwillison.net/... |

```go
// Create a formatter
formatter := citation.NewFormatter(citation.StyleMLA)

// Format a single reference
formatted := formatter.FormatReference(&ref)

// Generate markdown section
markdown := citation.FormatMarkdown(coll, formatter, "Works Cited")
```

### citegen CLI Tool

The `citegen` command-line tool generates formatted reference lists from JSON files.

#### Installation

```bash
go install github.com/grokify/gopandoc/cmd/citegen@latest
```

#### Usage

```bash
# Generate references in different styles
citegen generate -i references.json -s simple
citegen generate -i references.json -s mla --heading "Works Cited"
citegen generate -i references.json -s apa --heading "References"
citegen generate -i references.json -s chicago

# Filter by tag
citegen generate -i references.json -s mla --tag aws

# Validate references
citegen validate -i references.json

# List all tags
citegen tags -i references.json
```

#### Reference JSON Format

```json
{
  "references": [
    {
      "id": "aws-mantle-podcast",
      "type": "podcast",
      "title": "Amazon Bedrock Mantle and Developing at the Speed of AI",
      "authors": [{"given": "Joe", "family": "Magherimov"}],
      "date": "2026-01-26",
      "url": "https://podcasts.apple.com/us/podcast/...",
      "publisher": "AWS",
      "container": "AWS Podcast",
      "episode": "753",
      "tags": ["aws", "ai-development"]
    },
    {
      "id": "willison-2026",
      "type": "blog",
      "title": "The Software Factory",
      "authors": [{"given": "Simon", "family": "Willison"}],
      "date": "2026-02-07",
      "url": "https://simonwillison.net/2026/Feb/7/software-factory/",
      "container": "Simon Willison's Weblog",
      "tags": ["software-factory"]
    }
  ]
}
```

#### Supported Reference Types

- `article` — Journal or news articles
- `blog` — Blog posts
- `book` — Books
- `podcast` — Podcast episodes
- `video` — Video content
- `website` — Websites or web applications
- `whitepaper` — Technical whitepapers

## Dependencies

- Added `github.com/spf13/cobra` v1.10.2 for CLI
- Updated `github.com/grokify/mogo` to v0.74.4

## Upgrading

This release is fully backward compatible. No changes required for existing code using gopandoc.

To use the new citation features:

```go
import "github.com/grokify/gopandoc/citation"
```

## Documentation

- [README](README.md) — Updated with citation package documentation
- [pkg.go.dev](https://pkg.go.dev/github.com/grokify/gopandoc/citation) — API documentation

## Contributors

- John Wang (@grokify)
- Claude Opus 4.5 (AI pair programmer)
