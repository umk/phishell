package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func readScript(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error reading script: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines []string
	prevEmpty := true

	for scanner.Scan() {
		line := strings.TrimRightFunc(scanner.Text(), unicode.IsSpace)

		if strings.HasPrefix(line, "#") {
			continue
		}

		if line == "" {
			if prevEmpty {
				continue
			}
			prevEmpty = true
		} else {
			prevEmpty = false
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}
