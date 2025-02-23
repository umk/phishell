package persona

import (
	"fmt"
	"os"
	"path/filepath"
)

func ensurePersonaFiles(jsDir string) error {
	pkgPath := filepath.Join(jsDir, "package.json")
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

	if err := os.MkdirAll(filepath.Join(jsDir, "src"), 0644); err != nil {
		return fmt.Errorf("failed to create src: %w", err)
	}
	if err := createPackageFile(jsDir); err != nil {
		return err
	}
	if err := createIndexFile(jsDir); err != nil {
		return err
	}

	return nil
}

func createPackageFile(jsDir string) error {
	content, err := formatPackageJSON(PackageJSONParams{})
	if err != nil {
		return fmt.Errorf("failed to format package.json: %w", err)
	}
	filePath := filepath.Join(jsDir, "package.json")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}
	return nil
}

func createIndexFile(jsDir string) error {
	content, err := formatIndexTS(IndexTSParams{})
	if err != nil {
		return fmt.Errorf("failed to format index.ts: %w", err)
	}
	filePath := filepath.Join(jsDir, "src", "index.ts")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write index.ts: %w", err)
	}
	return nil
}
