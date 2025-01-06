package termx

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func NewScreen() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("error creating screen: %w", err)
	}
	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("error initializing screen: %w", err)
	}

	return screen, nil
}

func DrawLines(screen tcell.Screen, lines []string, startX, startY int, style tcell.Style) {
	for y, l := range lines {
		DrawText(screen, l, startX, startY+y, style)
	}
}

func DrawText(screen tcell.Screen, s string, startX, startY int, style tcell.Style) {
	cx := startX

	var (
		base      rune
		combining []rune
	)

	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			combining = append(combining, r)
			continue
		}

		if base != 0 {
			w := runewidth.RuneWidth(base)
			screen.SetContent(cx, startY, base, combining, style)
			cx += w
		}

		base = r
		combining = nil
	}

	if base != 0 {
		screen.SetContent(cx, startY, base, combining, style)
	}
}
