package fsx

import (
	"os"
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func GlobsWalk(basePath string, includes, excludes []string, fn doublestar.GlobWalkFunc) error {
	for _, pattern := range includes {
		curDir, p := doublestar.SplitPattern(pattern)
		curDir = Resolve(basePath, curDir)

		fs := os.DirFS(curDir)
		if err := doublestar.GlobWalk(fs, p, func(p string, d os.DirEntry) error {
			p = path.Clean(p)
			for _, exc := range excludes {
				seg := strings.Split(p, string(os.PathSeparator))
				for _, segment := range seg {
					if segment == exc {
						return nil
					}
				}
				ok, err := doublestar.Match(exc, p)
				if err != nil {
					return err
				}
				if ok {
					return nil
				}
			}
			return fn(p, d)
		}); err != nil {
			return err
		}
	}

	return nil
}
