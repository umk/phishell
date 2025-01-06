package stringsx

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func Truncate(s string, n int) string {
	const ellipsis = "..."

	if len(s) < n {
		return s
	}

	var b strings.Builder
	b.Grow(n * 2)

	c := 0

	for _, r := range s {
		c += runewidth.RuneWidth(r)

		if c > n-len(ellipsis) {
			b.WriteString(ellipsis)

			break
		}

		b.WriteRune(r)
	}

	return b.String()
}

func Every(s string, test func(rune) bool) bool {
	for _, r := range s {
		if !test(r) {
			return false
		}
	}

	return true
}
