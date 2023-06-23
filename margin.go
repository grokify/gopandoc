package gopandoc

import "fmt"

func MarginHeaderLines(geo Geometry) []string {
	geoString := ""
	if geo != (Geometry{}) {
		geoString = geo.String()
	}
	return []string{
		"---",
		//	"title: " + title,
		//	"",
		"header-includes:",
		" \\usepackage{geometry}",
		fmt.Sprintf(" \\geometry{%s}", geoString),
		"",
		"output:",
		"  pdf_document",
		"---",
	}
}

type Geometry struct {
	Unit   string
	Margin int
	Left   int
	Right  int
	Top    int
	Bottom int
}

func (g Geometry) String() string {
	if g.Margin > 0 {
		return fmt.Sprintf("margin=%d%s", g.Margin, g.Unit)
	}
	format := `geometry: "left=%d%s,right=%d%s,top=%d%s,bottom=%d%s"`
	return fmt.Sprintf(format, g.Left, g.Unit, g.Right, g.Unit, g.Top, g.Unit, g.Bottom, g.Unit)
}

func NewGeometry(unit string, width int) Geometry {
	return Geometry{
		Unit:   unit,
		Margin: width,
		Left:   width,
		Right:  width,
		Top:    width,
		Bottom: width,
	}
}

/*

func (bw BorderWidth) PandocMarkdownMargins() []string {

	return []string{
		"---",
		fmt.Sprintf(format, bw.Left, bw.Unit, bw.Right, bw.Unit, bw.Top, bw.Unit, bw.Bottom, bw.Unit),
		"---",
	}
}
*/
