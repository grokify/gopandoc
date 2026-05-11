# GoPandoc

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/gopandoc/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/gopandoc/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/gopandoc/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/gopandoc/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/gopandoc/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/gopandoc/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gopandoc
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gopandoc
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gopandoc
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gopandoc
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fgopandoc
 [loc-svg]: https://tokei.rs/b1/github/grokify/gopandoc
 [repo-url]: https://github.com/grokify/gopandoc
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gopandoc/blob/main/LICENSE

Go helper library for [Pandoc](https://pandoc.org/) document generation.

## Features

- **Document Generation**: Helper functions for generating Pandoc-compatible markdown
- **Page Layout**: Geometry and margin configuration utilities
- **Multi-format Output**: Write to markdown, PDF, and DOCX formats
- **Citation Management**: Bibliographic references with multiple citation styles

## Installation

```bash
go get github.com/grokify/gopandoc
```

## Packages

### gopandoc (root)

Core Pandoc helper functions for document generation.

```go
import "github.com/grokify/gopandoc"

// Configure page margins
geom := gopandoc.NewGeometry("in", 1)
header := gopandoc.MarginHeaderLines(geom)

// Generate markdown with margins
content := gopandoc.MarkdownLines("in", 1, []string{
    "# My Document",
    "Some content here.",
})

// Write to multiple formats
gopandoc.WriteFiles("document.md", content, true, true) // DOCX and PDF
```

### citation

Bibliographic reference management with multiple citation styles.

```go
import "github.com/grokify/gopandoc/citation"

// Load references from JSON
coll, err := citation.LoadJSON("references.json")

// Format in different styles
simple := citation.NewFormatter(citation.StyleSimple)
mla := citation.NewFormatter(citation.StyleMLA)
apa := citation.NewFormatter(citation.StyleAPA)
chicago := citation.NewFormatter(citation.StyleChicago)

// Generate markdown references section
output := citation.FormatMarkdown(coll, mla, "Works Cited")

// Filter by tag
awsRefs := coll.GetByTag("aws")

// Validate references
if errs := coll.Validate(); len(errs) > 0 {
    // Handle validation errors
}
```

#### Reference JSON Format

```json
{
  "references": [
    {
      "id": "willison-2026",
      "type": "blog",
      "title": "The Software Factory",
      "authors": [{"given": "Simon", "family": "Willison"}],
      "date": "2026-02-07",
      "url": "https://simonwillison.net/2026/Feb/7/software-factory/",
      "container": "Simon Willison's Weblog",
      "tags": ["software-factory", "ai"]
    }
  ]
}
```

#### Supported Citation Styles

| Style | Description | Example |
|-------|-------------|---------|
| Simple | Readable format | Willison, Simon. "The Software Factory." February 2026. |
| MLA | MLA 9th edition | Willison, Simon. "The Software Factory." *Simon Willison's Weblog*, 7 Feb. 2026. |
| APA | APA 7th edition | Willison, S. (2026, February 7). The software factory. |
| Chicago | Chicago Manual of Style | Willison, Simon. "The Software Factory." *Simon Willison's Weblog*, February 7, 2026. |

## CLI Tools

### citegen

Generate formatted reference lists from JSON files.

```bash
# Install
go install github.com/grokify/gopandoc/cmd/citegen@latest

# Generate references
citegen generate -i references.json -s simple
citegen generate -i references.json -s mla --heading "Works Cited"
citegen generate -i references.json -s apa --tag aws

# Validate references
citegen validate -i references.json

# List all tags
citegen tags -i references.json
```

## External References

1. [Add styling rules in pandoc tables for odt/docx output (table borders) (Stack Overflow)](https://stackoverflow.com/questions/17858598/add-styling-rules-in-pandoc-tables-for-odt-docx-output-table-borders)
2. [How do I add custom formatting to docx files generated in Pandoc? (Stack Overflow)](https://stackoverflow.com/questions/70513062/how-do-i-add-custom-formatting-to-docx-files-generated-in-pandoc)
