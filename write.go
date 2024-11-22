package gopandoc

import "os"

func WriteFilesLines(basename string, data []string, marginUnit string, marginScalar int) error {
	return WriteFiles(basename, []byte(MarkdownLines(marginUnit, marginScalar, data)))
}

func WriteFiles(basename string, data []byte) error {
	mdfile := basename + ".md"
	pdffile := basename + ".pdf"
	docxfile := basename + ".docx"

	err := os.WriteFile(mdfile, data, 0600)
	if err != nil {
		return err
	}

	popts := &PandocOpts{
		InputFiles: []string{mdfile},
		OutputFile: pdffile,
		FromFormat: FormatMarkdown,
		ToFormat:   FormatPDF,
		Margin:     "",
	}
	err = Exec(popts)
	if err != nil {
		return err
	}

	popts2 := &PandocOpts{
		InputFiles: []string{mdfile},
		OutputFile: docxfile,
		FromFormat: FormatMarkdown,
		ToFormat:   FormatDOCX,
		Margin:     "",
	}
	err = Exec(popts2)

	return err
}
