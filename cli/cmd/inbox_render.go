package cmd

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	"github.com/umk/phishell/cli/session"
	"github.com/umk/phishell/util/stringsx"
	"github.com/umk/phishell/util/termx"
)

type inboxRender struct {
	screen  tcell.Screen
	context *Context

	selectedIndex int
	startIndex    int

	inbox []*session.InboxMessage
	items []string
}

func newInboxRender(screen tcell.Screen, context *Context) (*inboxRender, error) {
	return &inboxRender{
		screen:  screen,
		context: context,
	}, nil
}

func (r *inboxRender) drawInbox(handler func(message *session.InboxMessage) error) error {
	r.refreshItems()

	for {
		r.drawInboxLoop(r.screen)

		switch ev := r.screen.PollEvent().(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				r.selectedIndex--
			case tcell.KeyDown:
				r.selectedIndex++
			case tcell.KeyLeft:
				r.selectedIndex = 0
			case tcell.KeyRight:
				r.selectedIndex = len(r.items) - 1
			case tcell.KeyEnter:
				if r.selectedIndex >= 0 {
					s := r.inbox[r.selectedIndex]
					r.context.session.Inbox.Delete(s.ID)

					return handler(s)
				}
			case tcell.KeyEscape:
				return nil
			case tcell.KeyBackspace2:
				if r.selectedIndex >= 0 {
					s := r.inbox[r.selectedIndex]
					r.context.session.Inbox.Delete(s.ID)

					r.refreshItems()
				}
			}
		case *tcell.EventResize:
			r.screen.Sync()
		}
	}
}

func (r *inboxRender) drawInboxLoop(screen tcell.Screen) {
	width, height := screen.Size()

	// Window dimensions (increased size)
	winWidth := width
	winHeight := height

	// Window dimensions (increased size)
	startX := (width - winWidth) / 2
	startY := (height - winHeight) / 2

	header := wordwrap.WrapString(inboxHeader, uint(winWidth-2))
	headerLines := strings.Split(header, "\n")

	visibleItems := winHeight - len(headerLines) - 5

	r.refreshSelection(visibleItems)

	screen.Clear()

	termx.DrawLines(screen, headerLines, startX+1, startY+1, tcell.StyleDefault.Foreground(tcell.ColorGray))

	r.drawInboxMenu(startX+1, startY+len(headerLines)+2, winWidth-2, visibleItems)
	r.drawInboxFooter(startX+1, startY+len(headerLines)+2, winWidth-2, visibleItems)

	screen.Show()
}

func (r *inboxRender) refreshSelection(visibleItems int) {
	if len(r.items) == 0 {
		r.selectedIndex = -1
		r.startIndex = 0

		return
	}

	if r.selectedIndex < 0 {
		r.selectedIndex = 0
	} else if r.selectedIndex >= len(r.items) {
		r.selectedIndex = len(r.items) - 1
	}

	if r.startIndex > r.selectedIndex {
		r.startIndex = r.selectedIndex
	} else if r.startIndex < r.selectedIndex-visibleItems+1 {
		r.startIndex = r.selectedIndex - visibleItems + 1
	}
}

func (r *inboxRender) drawInboxMenu(startX, startY, menuWidth, visibleItems int) {
	if len(r.items) == 0 {
		termx.DrawText(r.screen, "No messages", startX, startY, tcell.StyleDefault.Foreground(tcell.ColorGray))
		return
	}

	endIndex := r.startIndex + visibleItems
	if endIndex > len(r.items) {
		endIndex = len(r.items)
	}

	// Render visible items
	for i := r.startIndex; i < endIndex; i++ {
		y := startY + (i - r.startIndex)
		item := stringsx.Truncate(r.items[i], menuWidth)
		text := fmt.Sprintf("%-*s", menuWidth, item)
		style := tcell.StyleDefault
		if i == r.selectedIndex {
			style = style.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
		}
		termx.DrawText(r.screen, text, startX, y, style)
	}
}

func (r *inboxRender) drawInboxFooter(startX, startY, menuWidth, visibleItems int) {
	footerX := startX + (menuWidth-runewidth.StringWidth(inboxFooter))/2
	footerY := startY + visibleItems + 1

	termx.DrawText(r.screen, inboxFooter, footerX, footerY, tcell.StyleDefault.Foreground(tcell.ColorYellow))
}

func (r *inboxRender) refreshItems() {
	r.inbox = r.context.session.Inbox.Messages()

	slices.SortFunc(r.inbox, func(a, b *session.InboxMessage) int {
		return -a.Date.Compare(b.Date)
	})

	n := len(r.inbox)
	r.items = make([]string, n)

	for i, m := range r.inbox {
		t := m.Date.Format(time.DateTime)
		r.items[i] = fmt.Sprintf("%s  %s", t, m.Content)
	}
}
