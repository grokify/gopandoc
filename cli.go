package gopandoc

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ErrOptsNil = errors.New("opts is nil")

type PandocOpts struct {
	FromFormat string
	OutputFile string
	ToFormat   string
	InputFiles []string
	Geometry   string
	Margin     string
}

func (p *PandocOpts) TrimSpace() {
	p.Geometry = strings.TrimSpace(p.Geometry)
	p.Margin = strings.TrimSpace(p.Margin)
	p.OutputFile = strings.TrimSpace(p.OutputFile)
	p.FromFormat = strings.TrimSpace(p.FromFormat)

	p.ToFormat = strings.TrimSpace(p.ToFormat)
	for i, input := range p.InputFiles {
		p.InputFiles[i] = strings.TrimSpace(input)
	}

}

func (p *PandocOpts) CLIArgs() []string {
	args := []string{}
	if len(p.Geometry) > 0 {
		args = append(args, "-V", fmt.Sprintf("'geometry:%s'", p.Geometry))
	}
	if len(p.Margin) > 0 {
		args = append(args, "-V", fmt.Sprintf("'geometry:margin=%s'", p.Margin))
	}
	if len(p.OutputFile) > 0 {
		args = append(args, "-o", p.OutputFile)
	}
	if len(p.FromFormat) > 0 {
		args = append(args, "-f", p.FromFormat)
	}
	if len(p.ToFormat) > 0 {
		args = append(args, "-t", p.ToFormat)
	}
	args = append(args, p.InputFiles...)
	return args
}

func PandocOptsExmample() PandocOpts {
	return PandocOpts{
		FromFormat: FormatMarkdown,
		ToFormat:   FormatDOCX,
		Margin:     "0.5cm",
	}
}

func Exec(opts *PandocOpts) error {
	if opts == nil {
		return ErrOptsNil
	}
	// https://www.mscharhag.com/software-development/pandoc-markdown-to-pdf
	args := opts.CLIArgs()
	// fmtutil.PrintJSON(args)
	fmt.Println(strings.Join(args, " "))
	cmd := exec.Command(CLICommand, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// pandoc -o my.docx -f markdown -t docx my.md
