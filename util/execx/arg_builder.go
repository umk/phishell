package execx

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type argumentsBuilder struct {
	arguments []Arguments
	current   *Arguments
}

func newArgumentsBuilder() *argumentsBuilder {
	builder := &argumentsBuilder{
		arguments: []Arguments{nil},
	}
	builder.current = &builder.arguments[0]

	return builder
}

func (a *argumentsBuilder) addPlain(arg string) error {
	if arg == "|" {
		a.arguments = append(a.arguments, nil)
		a.current = &a.arguments[len(a.arguments)-1]
		return nil
	}

	if strings.ContainsAny(arg, "&|><") {
		return fmt.Errorf("operator is not supported: %s", arg)
	}

	a.add(arg, true)

	return nil
}

func (a *argumentsBuilder) addQuoted(arg string, expandEnvVars bool) {
	a.add(arg, expandEnvVars)
}

func (a *argumentsBuilder) add(arg string, expandEnvVars bool) {
	if expandEnvVars {
		arg = expandEnvVariables(arg)
	}

	a.current.Add(arg)
}

func expandEnvVariables(s string) string {
	var result strings.Builder
	var varName strings.Builder
	inVar := false

	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		if inVar {
			if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
				varName.WriteRune(c)
			} else {
				expanded := expandEnvVariable(varName.String())
				result.WriteString(expanded)

				varName.Reset()

				inVar = false
				result.WriteRune(c)
			}
		} else {
			if c == '$' {
				inVar = true
			} else {
				result.WriteRune(c)
			}
		}
	}

	if inVar {
		expanded := expandEnvVariable(varName.String())
		result.WriteString(expanded)
	}

	return result.String()
}

func expandEnvVariable(varName string) string {
	if varName == "" {
		return "$"
	} else {
		return os.Getenv(varName)
	}
}
