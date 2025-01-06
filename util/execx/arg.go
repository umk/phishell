package execx

import (
	"fmt"
	"strings"
	"unicode"
)

type Arguments []string

func (c Arguments) String() string {
	var parts []string
	for _, s := range c {
		if strings.ContainsFunc(s, unicode.IsSpace) {
			s = strings.ReplaceAll(s, "'", "\\'")
			s = fmt.Sprintf("'%s'", s)
		}
		parts = append(parts, s)
	}

	return strings.Join(parts, " ")
}

func (c *Arguments) Add(arg string) {
	*c = append(*c, arg)
}

func (c Arguments) Get(n int) string {
	if len(c) <= n {
		return ""
	}

	return c[n]
}
