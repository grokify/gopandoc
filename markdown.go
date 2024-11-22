package gopandoc

import (
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

func MarkdownLines(marginUnit string, marginScalar int, lines []string) string {
	lines = stringsutil.SliceCondenseSpace(lines, false, false)
	str := strings.Join(lines, "\n\n")

	marginUnit = strings.TrimSpace(marginUnit)
	if marginUnit != "" {
		header := MarginHeaderLines(NewGeometry(marginUnit, marginScalar))
		hstr := strings.Join(header, "\n")
		str = hstr + "\n\n" + str
	}
	return str
}
