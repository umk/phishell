package termx

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"unicode"

	"github.com/umk/phishell/util/errorsx"
	"golang.org/x/term"
)

// ReadLine reads a line from the terminal based on a prompt
func ReadLine(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	// Display the prompt
	fmt.Print(prompt)

	// Read the line from terminal input
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Trim the newline character and return the result
	return input[:len(input)-1], nil
}

// ReadKeySilent writes the prompt to the terminal and waits for a single key press silently.
// It returns the key pressed as a rune, handling multi-byte characters.
func ReadKeySilent(prompt ...string) (rune, error) {
	fd := int(os.Stdin.Fd())

	// Save the original terminal state to restore later
	oldState, err := term.GetState(fd)
	if err != nil {
		return 0, err
	}
	defer term.Restore(fd, oldState)

	// Put terminal into raw mode to read single characters
	if _, err := term.MakeRaw(fd); err != nil {
		return 0, err
	}

	// Write the prompt
	for _, p := range prompt {
		fmt.Print(p)
	}

	// Create a buffered reader for stdin
	reader := bufio.NewReader(os.Stdin)

	// Read a single UTF-8 encoded rune
	r, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}

	// Handle Ctrl+C (ASCII ETX or 0x03)
	if r == 0x03 {
		return 0, errorsx.ErrInterrupted
	}

	return r, nil
}

// ReadKey writes the prompt to the terminal and waits for a single key press.
// It returns the key pressed as a rune, handling multi-byte characters, and echoes it.
func ReadKey(prompt ...string) (rune, error) {
	r, err := ReadKeySilent(prompt...)
	if err != nil {
		fmt.Println() // Add a newline
		return 0, err
	}
	fmt.Println(string(r)) // Echo the character and add a newline
	return r, nil
}

// ReadKeyOfSilent writes the prompt and waits for a key press that matches one of the provided runes silently.
// It keeps prompting until a valid key is pressed and returns that rune.
func ReadKeyOfSilent(prompt string, keys ...rune) (rune, error) {
	fmt.Print(prompt)
	for {
		key, err := ReadKeySilent()
		if err != nil {
			return 0, err
		}
		// Convert to lower case if it's a letter for case-insensitive matching
		keyLower := unicode.ToLower(key)
		if slices.Index(keys, keyLower) >= 0 {
			return key, nil
		}
	}
}

// ReadKeyOf writes the prompt and waits for a key press that matches one of the provided runes.
// It keeps prompting until a valid key is pressed, echoes it, and returns that rune.
func ReadKeyOf(prompt string, keys ...rune) (rune, error) {
	key, err := ReadKeyOfSilent(prompt)
	if err != nil {
		return 0, err
	}
	fmt.Println(string(key)) // Echo the character and add a newline
	return key, nil
}

// ReadKeyOrDefaultOfSilent writes the prompt and waits for a key press that matches one of the provided runes silently.
// If the Enter key (rune '\r') is pressed, it returns the default rune.
func ReadKeyOrDefaultOfSilent(prompt string, def rune, others ...rune) (rune, error) {
	// Include '\r' (Enter key) as a valid key
	keys := append([]rune{'\r', def}, others...)
	result, err := ReadKeyOfSilent(prompt, keys...)
	if err != nil {
		return 0, err
	}
	if result == '\r' {
		return def, nil
	}
	return result, nil
}

// ReadKeyOrDefaultOf writes the prompt and waits for a key press that matches one of the provided runes.
// If the Enter key (rune '\r') is pressed, it echoes and returns the default rune.
func ReadKeyOrDefaultOf(prompt string, def rune, others ...rune) (rune, error) {
	result, err := ReadKeyOrDefaultOfSilent(prompt, def, others...)
	if err != nil {
		fmt.Println() // Add a newline
		return 0, err
	}
	fmt.Println(string(result)) // Echo the character and add a newline
	return result, nil
}
