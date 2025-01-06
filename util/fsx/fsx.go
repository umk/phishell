package fsx

import (
	"fmt"
	"path/filepath"
	"strings"
)

// EscapesBaseDir gets a value indicating whether provided relative path escapes
// the base directory. The function panics if provided path is not relative.
func EscapesBaseDir(p string) bool {
	if filepath.IsAbs(p) {
		panic(fmt.Sprintf("expected a relative path, but got absolute: %s", p))
	}

	// Normalize the path for the current operating system
	v := filepath.Clean(p)

	i := strings.Index(v, string(filepath.Separator))
	if i >= 0 {
		v = v[:i]
	}

	return v == ".."
}

// Resolve resolves the given path relative to the basePath if it's not absolute.
func Resolve(basePath, p string) string {
	// If value is an absolute path, return it as is.
	if filepath.IsAbs(p) {
		return p
	}

	// Otherwise, append the relative value to the basePath.
	return filepath.Join(basePath, p)
}
