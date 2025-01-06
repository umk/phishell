package stringsx

import (
	"unicode"
)

// TrimEmpty trims empty or whitespace-only lines from the start and end of a slice of lines.
func TrimEmpty(lines []string) []string {
	start := 0
	for start < len(lines) && Every(lines[start], unicode.IsSpace) {
		start++
	}
	end := len(lines) - 1
	for end >= start && Every(lines[end], unicode.IsSpace) {
		end--
	}
	return lines[start : end+1]
}
