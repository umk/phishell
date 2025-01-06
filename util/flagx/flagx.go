package flagx

import "strings"

// Strings is a custom type that holds multiple string values.
type Strings []string

// String returns a textual representation of the Strings.
func (s *Strings) String() string {
	return strings.Join(*s, ",")
}

// Set appends a new string to the slice.
func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}
