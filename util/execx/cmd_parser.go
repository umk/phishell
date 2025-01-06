package execx

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var envVarRegex = regexp.MustCompile(`^[\w\d_]+$`)

func Parse(input string) ([]Arguments, error) {
	builder := newArgumentsBuilder()

	var current strings.Builder

	var inSingleQuote, inDoubleQuote bool

	var escapeNext bool

	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		inQuote := inSingleQuote || inDoubleQuote

		if escapeNext {
			current.WriteRune(c)
			escapeNext = false
			continue
		}

		if c == '\\' && !inQuote {
			escapeNext = true
			continue
		}

		if c == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			if !inSingleQuote {
				builder.addQuoted(current.String(), false)
				current.Reset()
			}
			continue
		}

		if c == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			if !inDoubleQuote {
				builder.addQuoted(current.String(), true)
				current.Reset()
			}
			continue
		}

		if !inQuote && (c == ' ' || c == '\t') {
			if current.Len() > 0 {
				if err := builder.addPlain(current.String()); err != nil {
					return nil, err
				}
				current.Reset()
			}
			continue
		}

		current.WriteRune(c)
	}

	if escapeNext {
		return nil, fmt.Errorf("incomplete escape sequence")
	}

	if inSingleQuote || inDoubleQuote {
		return nil, fmt.Errorf("unterminated quote")
	}

	if current.Len() > 0 {
		if err := builder.addPlain(current.String()); err != nil {
			return nil, err
		}
	}

	return builder.Arguments, nil
}

func AllocArgs(args Arguments) (*Cmd, error) {
	if len(args) == 0 {
		return nil, errors.New("command line is empty")
	}

	i := 0

	for ; i < len(args)-1; i++ {
		part := args[i]
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 || !envVarRegex.MatchString(kv[0]) {
			break
		}
	}

	return &Cmd{
		Env:  args[:i],
		Cmd:  args[i],
		Args: args[i+1:],
	}, nil
}
