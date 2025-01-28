package execx

import (
	"errors"
	"io"
	"strings"
	"sync"
	"unicode/utf8"
)

type outputWrapper struct {
	mu sync.Mutex

	base   io.Writer // base output writer
	maxLen int       // max size of the output tracked by the wrapper in bytes

	init bool // indicates if there was any output
	tail bool // indicates if the output was truncated

	lines    []string
	linesLen int // total length of the lines in bytes

	cur struct {
		b   strings.Builder // current line until it exceeds max length
		len int             // length of the current line in bytes
	}
	incompl []byte // bytes that appear to be a head of multi-byte character

	err error // error that occurred during writing
}

func newOutputWrapper(base io.Writer, maxLen int) *outputWrapper {
	if maxLen <= 0 {
		panic("maximum log length cannot be less or equal to zero")
	}

	return &outputWrapper{
		base:   base,
		maxLen: maxLen,
	}
}

func (w *outputWrapper) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.err == nil {
		if err := w.processOutput(p); err != nil {
			w.err = err

			w.resetOnError()
		}
	}

	if w.base != nil {
		if n, err := w.base.Write(p); err != nil {
			return n, err
		}
	}

	return len(p), nil
}

func (w *outputWrapper) finalize() error {
	if !w.init {
		return nil
	}

	if len(w.incompl) > 0 {
		return errors.New("invalid UTF-8 encoding")
	}

	w.appendCurrent()

	w.init = false

	return nil
}

func (w *outputWrapper) Get() ProcessOutput {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.err == nil {
		if err := w.finalize(); err != nil {
			w.err = err

			w.resetOnError()
		}
	}

	if w.err != nil {
		return ProcessOutput{err: w.err}
	}

	return ProcessOutput{
		lines: w.lines,
		tail:  w.tail,
		err:   w.err,
	}
}

func (w *outputWrapper) processOutput(p []byte) error {
	w.init = true

	b := p
	if len(w.incompl) > 0 {
		b = append(w.incompl, p...)
		w.incompl = nil
	}

	for i := 0; i < len(b); {
		if b[i] == '\n' {
			i++
			if i < len(b)-1 && b[i+1] == '\r' {
				i++
			}

			w.appendCurrent()
		} else if b[i] < utf8.RuneSelf {
			if w.cur.len < w.maxLen {
				w.cur.b.WriteByte(b[i])
				w.cur.len++

				if w.cur.len == w.maxLen {
					w.cur.b.Reset()
				}
			}
			i++
		} else if _, width := utf8.DecodeRune(b[i:]); width > 1 {
			if w.cur.len < w.maxLen {
				w.cur.b.Write(b[i : i+width])
				w.cur.len += width

				if w.cur.len >= w.maxLen {
					w.cur.b.Reset()
				}
			}
			i += width
		} else if !utf8.FullRune(b[i:]) {
			w.incompl = b[i:]
			break
		} else {
			return errors.New("invalid UTF-8 encoding")
		}
	}

	return nil
}

func (w *outputWrapper) appendCurrent() {
	if w.cur.len > w.cur.b.Len() {
		// The line exceeds the max length, so everything is discarded.
		w.lines = nil
		w.linesLen = 0

		w.tail = true

		w.resetCur()

		return
	}

	// The line fits the max length, so it is added to the list.
	w.lines = append(w.lines, w.cur.b.String())
	w.linesLen += w.cur.b.Len() + 1

	w.resetCur()

	// The list is truncated if it exceeds the max length.
	var j int
	for w.linesLen > w.maxLen {
		w.linesLen -= len(w.lines[j]) + 1
		w.lines[j] = ""

		j++
	}

	if j > 0 {
		w.lines = w.lines[j:]
		w.tail = true
	}
}

func (w *outputWrapper) resetCur() {
	w.cur.b.Reset()
	w.cur.len = 0
}

func (w *outputWrapper) resetOnError() {
	w.tail = false

	w.resetCur()

	w.lines = nil
	w.incompl = nil
}
