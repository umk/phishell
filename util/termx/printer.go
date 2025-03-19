package termx

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
	"golang.org/x/term"
)

var MD MDPrinter

func init() {
	for _, c := range []*ansi.StyleConfig{
		&styles.LightStyleConfig,
		&styles.DarkStyleConfig,
	} {
		c.Item.BlockPrefix = " - "
	}

	MD.Init()
}

type MDPrinter struct {
	renderer *glamour.TermRenderer
}

func (p *MDPrinter) Init() {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 80
	}

	MD.renderer, _ = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4),
	)
}

func (p *MDPrinter) Println(a ...any) {
	p.Print(fmt.Sprintln(a...))
}

func (p *MDPrinter) Printf(format string, a ...any) {
	p.Print(fmt.Sprintf(format, a...))
}

func (p *MDPrinter) Print(s string) {
	if p.renderer != nil {
		f, err := p.renderer.Render(s)
		if err == nil {
			fmt.Print(f)
			return
		}
	}

	fmt.Println(s)
}
