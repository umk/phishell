package stringsx

import (
	"strings"
	"unicode"
)

type Token []rune

func (t Token) String() string {
	return string(t)
}

func Tokens(id string) []Token {
	var a []Token

	for _, s := range strings.FieldsFunc(id, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if s == "" {
			continue
		}

		r := []rune(s)

		k := 0
		for i := 1; i < len(r); i++ {
			if unicode.IsDigit(r[i-1]) != unicode.IsDigit(r[i]) ||
				(unicode.IsUpper(r[i]) &&
					((i-1 >= 0 && unicode.IsLower(r[i-1])) ||
						(i+1 < len(r) && unicode.IsLower(r[i+1])))) {
				a = append(a, r[k:i])
				k = i
			}
		}
		a = append(a, r[k:])
	}

	return a
}

func DisplayName(tokens []Token) string {
	var a []rune

	for _, t := range tokens {
		if len(t) == 0 {
			continue
		}

		if len(a) > 0 {
			a = append(a, ' ')
		}

		caps := true
		for _, c := range t {
			if !unicode.IsUpper(c) {
				caps = false
				break
			}
		}

		if caps && len(t) >= 2 {
			a = append(a, t...)
		} else {
			a = append(a, unicode.ToTitle(t[0]))

			for j := 1; j < len(t); j++ {
				a = append(a, unicode.ToLower(t[j]))
			}
		}
	}

	return string(a)
}
