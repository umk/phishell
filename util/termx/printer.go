package termx

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

type Printer struct {
	base *glamour.TermRenderer
}

func NewPrinter() *Printer {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 80
	}

	base, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4),
	)

	return &Printer{base: base}
}

func (p *Printer) Printf(format string, a ...any) {
	p.Print(fmt.Sprintf(format, a...))
}

func (p *Printer) Print(s string) {
	if p.base != nil {
		f, err := p.base.Render(s)
		if err == nil {
			fmt.Print(f)
			return
		}
	}

	fmt.Println(s)
}
