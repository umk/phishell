package execx

import (
	"errors"
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

func (c Arguments) Cmd() (*Cmd, error) {
	if len(c) == 0 {
		return nil, errors.New("command line is empty")
	}

	i := 0

	for ; i < len(c)-1; i++ {
		part := c[i]
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 || !envVarRegex.MatchString(kv[0]) {
			break
		}
	}

	return &Cmd{
		Env:  c[:i],
		Cmd:  c[i],
		Args: c[i+1:],
	}, nil
}
