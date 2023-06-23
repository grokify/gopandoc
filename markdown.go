package gopandoc

import (
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

func MarkdownLines(marginUnit string, marginScalar uint, lines []string) string {
	lines = stringsutil.SliceCondenseSpace(lines, false, false)
	str := strings.Join(lines, "\n\n")

	marginUnit = strings.TrimSpace(marginUnit)
	if marginUnit != "" {
		header := MarginHeaderLines(NewGeometry(marginUnit, int(marginScalar)))
		hstr := strings.Join(header, "\n")
		str = hstr + "\n\n" + str
	}
	return str
}
