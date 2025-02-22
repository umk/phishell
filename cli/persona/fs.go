package persona

import (
	"fmt"
	"os"
	"path/filepath"
)

func createPersonaFiles(pdir string) error {
	pkgPath := filepath.Join(pdir, "package.json")
	stat, err := os.Stat(pkgPath)
	if err == nil {
		if stat.IsDir() {
			return fmt.Errorf("expected file but %s is a directory", pkgPath)
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check package.json: %w", err)
	}

	if err := createPackageFile(pdir); err != nil {
		return err
	}
	if err := createIndexFile(pdir); err != nil {
		return err
	}

	return nil
}

func createPackageFile(pdir string) error {
	pkg := filepath.Join(pdir, "package.json")
	cont, err := formatPackageJSON(PackageJSONParams{})
	if err != nil {
		return fmt.Errorf("failed to format package.json: %w", err)
	}
	if err := os.WriteFile(pkg, []byte(cont), 0644); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}
	return nil
}

func createIndexFile(pdir string) error {
	src := filepath.Join(pdir, "src")
	if err := os.MkdirAll(src, 0644); err != nil {
		return fmt.Errorf("failed to create src: %w", err)
	}
	cont, err := formatIndexTS(IndexTSParams{})
	if err != nil {
		return fmt.Errorf("failed to format index.ts: %w", err)
	}
	idx := filepath.Join(src, "index.ts")
	if err := os.WriteFile(idx, []byte(cont), 0644); err != nil {
		return fmt.Errorf("failed to write index.ts: %w", err)
	}
	return nil
}
