package fsx

import (
	"path/filepath"
	"regexp"
)

var fsInvalidChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

// Resolve resolves the given path relative to the basePath if it's not absolute.
func Resolve(basePath, p string) string {
	// If value is an absolute path, return it as is.
	if filepath.IsAbs(p) {
		return p
	}

	// Otherwise, append the relative value to the basePath.
	return filepath.Join(basePath, p)
}

// Normalize replaces characters that are not allowed in file paths with an underscore.
func Normalize(name string) string {
	return fsInvalidChars.ReplaceAllString(name, "_")
}
