package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/util/termx"
	"golang.org/x/term"
)

func (c *Cli) processPrompt(ctx context.Context) error {
	line, err := c.readLine(ctx)
	if err != nil {
		return err
	}

	// Handle empty content
	content := strings.TrimSpace(line)
	if content == "" {
		return nil
	}

	switch c.mode {
	case PrCommand:
		return c.processCommand(ctx, content)
	case PrChat:
		client := c.getClient(ctx)
		return c.session.ProcessChat(ctx, client, content)
	}

	return nil
}

func (c *Cli) getPromptAndHint(ctx context.Context, mode PromptMode, inboxSize int) (prompt string, hint string) {
	prompt = c.getPrompt(ctx, mode, inboxSize)
	hint = c.getHint(mode)

	return
}

func (c *Cli) getPrompt(ctx context.Context, mode PromptMode, inboxSize int) string {
	var prompt string

	switch mode {
	case PrCommand:
		prompt = c.getCommandPrompt()
	default:
		prompt = c.getChatPrompt(ctx)
	}

	var inbox string
	if inboxSize > 0 {
		inbox = fmt.Sprintf("(%d) ", inboxSize)
	}

	return inbox + prompt
}

func (c *Cli) getChatPrompt(ctx context.Context) string {
	client := c.getClient(ctx)

	return fmt.Sprintf("%s >>> ", client.Config.Profile)
}

func (c *Cli) getHint(mode PromptMode) string {
	if mode == PrCommand {
		return ""
	}

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	homeDir, err := os.UserHomeDir()

	if err == nil && dir == homeDir {
		return "~"
	} else if dir == "/" {
		return "/"
	} else {
		dirName := filepath.Base(dir)
		return dirName + "/"
	}
}

func (c *Cli) getCommandPrompt() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	var base string

	homeDir, err := os.UserHomeDir()

	if err == nil && dir == homeDir {
		base = "~"
	} else if dir == "/" {
		base = "/"
	} else {
		dirName := filepath.Base(dir)
		base = dirName
	}

	return base + " $ "
}

func (c *Cli) readLine(ctx context.Context) (string, error) {
	prompt, hint := c.getPromptAndHint(ctx, c.mode, c.session.Inbox.Len())

	fmt.Print(prompt)

	// Put terminal into raw mode
	s, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("error entering raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), s)

	var line []rune
	atStart := true

	// Get the terminal size
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80 // Default width
	}

	cursorCol := utf8.RuneCountInString(prompt)
	cursorRow := 0

	reader := bufio.NewReader(os.Stdin)

	for {
		n := utf8.RuneCountInString(hint)

		if len(line) == 0 {
			termx.Muted.Print(hint)
			fmt.Print(strings.Repeat("\b", n))
		}

		// Read one byte
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}

		if len(line) == 0 && hint != "" {
			fmt.Print(strings.Repeat(" ", n) + strings.Repeat("\b", n))
		}

		// Handle control characters
		switch b {
		case '\n', '\r':
			fmt.Print("\r\n")
			return string(line), nil
		case 3: // Ctrl+C
			fmt.Print("\r\n")
			return "", nil
		case 4: // Ctrl+D (EOF)
			fmt.Print("\r\n")
			return "", io.EOF
		case 127, '\b': // Backspace
			n := len(line)
			if n > 0 {
				line = line[:n-1]
				// Move the cursor back
				if cursorCol == 0 {
					// Move up one line
					cursorRow--
					cursorCol = width - 1
					// Move cursor up one line and to the end
					fmt.Printf("\x1b[1A\x1b[%dC ", cursorCol)
				} else {
					// Move on character back
					cursorCol--
					// Move the cursor back and erase character
					fmt.Print("\b \b")
				}
			}
			atStart = len(line) == 0
		case '\t': // Tab key
			if atStart {
				// Clear the line
				var n int
				if len(line) == 0 {
					n = utf8.RuneCountInString(prompt) + utf8.RuneCountInString(hint)
				} else {
					n = utf8.RuneCountInString(prompt) + len(line)
				}
				fmt.Print("\r" + strings.Repeat(" ", n))
				// Switch modes
				c.mode = c.getNextMode(ctx, c.mode)
				prompt, hint = c.getPromptAndHint(ctx, c.mode, c.session.Inbox.Len())
				fmt.Print("\r" + prompt)
				line = []rune{}

				cursorCol = utf8.RuneCountInString(prompt)
				cursorRow = 0
			}
		case 0x1b: // Escape character
			if reader.Buffered() >= 2 {
				// Potential start of an escape sequence
				seq, err := reader.Peek(2)
				if err == nil && (seq[0] == '[' || seq[0] == 'O') {
					// It's an escape sequence, consume it
					reader.ReadByte() // Consume the next byte ('[' or 'O')
					// Read the rest of the escape sequence
					for {
						b, err := reader.ReadByte()
						if err != nil {
							break
						}
						// Escape sequences typically end with a letter
						if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '~' {
							break
						}
					}
					// Ignore the escape sequence (do nothing)
					continue
				}
			}
			// Not an escape sequence, but lone escape character
			r := ""
			if len(line) == 0 && c.mode == PrCommand {
				r = "inbox" // A default command triggered by escape key
			}

			fmt.Printf("%s\r\n", r)
			return r, nil
		default:
			// Handle multi-byte UTF-8 characters
			// Since we've already read one byte, we might need to read more bytes to complete the rune
			var buf [utf8.UTFMax]byte
			buf[0] = b
			n := 1
			var r rune
			var size int

			if b >= 0x80 {
				// Multibyte UTF-8 character
				for n < utf8.UTFMax {
					if utf8.FullRune(buf[:n]) {
						r, size = utf8.DecodeRune(buf[:n])
						break
					}
					nextByte, err := reader.ReadByte()
					if err != nil {
						return "", err
					}
					buf[n] = nextByte
					n++
				}
				if !utf8.FullRune(buf[:n]) {
					// Invalid UTF-8 sequence
					continue
				}
			} else {
				// Single-byte ASCII character
				r = rune(b)
				size = 1
			}

			if r == utf8.RuneError && size == 1 {
				// Invalid rune
				continue
			}

			line = append(line, r)
			fmt.Print(string(r))
			atStart = false

			cursorCol++
			if cursorCol == width {
				cursorCol = 0
				cursorRow++

				fmt.Print("\r\n")
			}
		}
	}
}

func (c *Cli) getNextMode(ctx context.Context, current PromptMode) PromptMode {
	app := bootstrap.GetApp(ctx)

	max := int(PrChat) + len(app.Clients)

	return PromptMode((int(current) + 1) % max)
}
