package gopandoc

import (
	"io"
	"os"

	"github.com/grokify/mogo/os/fileext"
)

func WriteFilesLines(basename string, data []string, marginUnit string, marginScalar int, stdout io.Writer, stderr io.Writer) error {
	return WriteFiles(basename, []byte(MarkdownLines(marginUnit, marginScalar, data)), stdout, stderr)
}

func WriteFiles(basename string, data []byte, stdout io.Writer, stderr io.Writer) error {
	fileMkdn := basename + "." + fileext.ExtMarkdown
	filePDF := basename + "." + fileext.ExtPDF
	fileDOCX := basename + "." + fileext.ExtDOCX

	if err := os.WriteFile(fileMkdn, data, 0600); err != nil {
		return err
	}

	popts := &PandocOpts{
		InputFiles: []string{fileMkdn},
		OutputFile: filePDF,
		FromFormat: FormatMarkdown,
		ToFormat:   FormatPDF,
		Margin:     "",
	}
	if err := Exec(popts, stdout, stderr); err != nil {
		return err
	}

	popts2 := &PandocOpts{
		InputFiles: []string{fileMkdn},
		OutputFile: fileDOCX,
		FromFormat: FormatMarkdown,
		ToFormat:   FormatDOCX,
		Margin:     "",
	}
	return Exec(popts2, stdout, stderr)
}
