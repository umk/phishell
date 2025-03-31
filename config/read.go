package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"maps"

	"github.com/umk/phishell/util/marshalx"
)

type ConfigFile struct {
	Default *string `yaml:"default" validate:"omitempty,min=1"`

	Profiles map[string]*ConfigFileProfile `yaml:"profiles"`
}

type ConfigFileProfile struct {
	Preset *string `yaml:"preset" validate:"omitempty,min=1"`

	BaseURL string `yaml:"baseurl" validate:"omitempty,url"`
	Key     string `yaml:"key" validate:"omitempty"`
	Model   string `yaml:"model" validate:"omitempty"`

	Prompt string `yaml:"prompt" validate:"omitempty"`

	Concurrency int `yaml:"concurrency" validate:"omitempty,min=1"`
	ContextSize int `yaml:"context" validate:"omitempty,min=2000"`

	Indexing ConfigFileIndexing `yaml:"indexing" validate:"dive"`
}

type ConfigFileIndexing struct {
	ChunkToks   int `yaml:"chunkToks" validate:"omitempty,min=1"`
	OverlapToks int `yaml:"overlapToks" validate:"omitempty,min=1"`
}

// loadConfigFiles reads configuration files and returns combined configuration.
func loadConfigFiles(currentDir string) (*ConfigFile, error) {
	config := ConfigFile{Profiles: make(map[string]*ConfigFileProfile)}

	currentDir = filepath.Clean(currentDir)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get user home directory: %w", err)
	}

	homeConfigPath := filepath.Join(homeDir, ".phishell.yaml")
	if err := LoadConfigFile(homeConfigPath, &config); err != nil {
		return nil, err
	}

	if homeDir != currentDir {
		dirConfigPath := filepath.Join(currentDir, ".phishell.yaml")
		if err := LoadConfigFile(dirConfigPath, &config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

// LoadConfigFile reads a YAML configuration file and updates combined config.
func LoadConfigFile(path string, config *ConfigFile) error {
	data, err := os.ReadFile(path)
	if err != nil {
		// If the file does not exist, skip without error
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("cannot read %s: %w", path, err)
	}

	var current ConfigFile
	if err := marshalx.UnmarshalYAMLStruct(data, &current); err != nil {
		return fmt.Errorf("cannot parse YAML: %w", err)
	}

	// Remove default profile if references profile in a local config.
	if config.Default != nil {
		if _, ok := current.Profiles[*config.Default]; ok {
			config.Default = nil
		}
	}

	// Copy settings from current config to combined one.

	if current.Default != nil {
		config.Default = current.Default
	}

	maps.Copy(config.Profiles, current.Profiles)

	return nil
}

func setProfileFromFileProfileOrPreset(target *Profile, source *ConfigFileProfile) error {
	if source.Preset != nil {
		preset, ok := presets[*source.Preset]
		if !ok {
			return fmt.Errorf("preset not found: %s", *source.Preset)
		}

		if err := setServiceFromProfile(target, &preset); err != nil {
			return err
		}
	}

	if err := setServiceFromProfile(target, source); err != nil {
		return err
	}

	if target.BaseURL == "" {
		target.BaseURL = presetOpenAI.BaseURL
	}

	return nil
}

func setServiceFromProfile(target *Profile, source *ConfigFileProfile) error {
	if source.BaseURL != "" {
		target.BaseURL = source.BaseURL
	}
	if source.Key != "" {
		target.Key = source.Key
	}
	if source.Model != "" {
		target.Model = source.Model
	}

	if source.Prompt != "" {
		target.Prompt = source.Prompt
	}

	if source.Concurrency > 0 {
		target.Concurrency = source.Concurrency
	}
	if source.ContextSize > 0 {
		target.ContextSize = source.ContextSize
	}

	if source.Indexing.ChunkToks > 0 {
		target.Indexing.ChunkToks = source.Indexing.ChunkToks
	}
	if source.Indexing.OverlapToks > 0 {
		target.Indexing.OverlapToks = source.Indexing.OverlapToks
	}

	return nil
}
