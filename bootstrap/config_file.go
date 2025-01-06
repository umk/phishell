package bootstrap

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/umk/phishell/util/fsx"
	"github.com/umk/phishell/util/marshalx"
)

const configFileName = ".phishell.yaml"

type ConfigFile struct {
	Default  *string                       `yaml:"default"`
	Profiles map[string]*ConfigFileProfile `yaml:"profiles"`
}

type ConfigFileProfile struct {
	Context ConfigFileProfileContext `yaml:"-"`

	Preset *string `yaml:"preset"`

	BaseURL string `yaml:"baseurl"`
	Key     string `yaml:"key"`
	Model   string `yaml:"model"`

	Prompt *ConfigFilePrompt `yaml:"prompt"`

	Retries        int `yaml:"retries"`
	Concurrency    int `yaml:"concurrency"`
	CompactionToks int `yaml:"compactionToks"`
}

type ConfigFileProfileContext struct {
	IsGlobal bool   // Indicates whether profile came from a global config
	Dir      string // Directory where the config is located
}

type ConfigFilePrompt struct {
	Path    *string `json:"path"`
	Content *string `json:"content"`
}

// loadConfigFiles reads configuration files and returns combined configuration.
func loadConfigFiles(currentDir string) (*ConfigFile, error) {
	config := ConfigFile{Profiles: make(map[string]*ConfigFileProfile)}

	currentDir = filepath.Clean(currentDir)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get user home directory: %w", err)
	}

	homeConfigPath := filepath.Join(homeDir, configFileName)
	if err := LoadConfigFile(homeConfigPath, &config, true); err != nil {
		return nil, err
	}

	if homeDir != currentDir {
		dirConfigPath := filepath.Join(currentDir, configFileName)
		if err := LoadConfigFile(dirConfigPath, &config, false); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

// LoadConfigFile reads a YAML configuration file and updates combined config.
func LoadConfigFile(path string, config *ConfigFile, isGlobal bool) error {
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

	// Remove default config profile if profile in current file overrides
	// the profile config in the file, where this default was defined.
	if config.Default != nil {
		if _, ok := current.Profiles[*config.Default]; ok {
			config.Default = nil
		}
	}

	// Copy settings from current config to combined one.

	if current.Default != nil {
		config.Default = current.Default
	}

	for k, v := range current.Profiles {
		v.Context = ConfigFileProfileContext{
			IsGlobal: isGlobal,
			Dir:      filepath.Dir(path),
		}

		config.Profiles[k] = v
	}

	return nil
}

func setServiceFromProfileOrPreset(target *ConfigService, source *ConfigFileProfile) error {
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

	if source.Retries > 0 {
		target.Retries = source.Retries
	}

	return nil
}

func setServiceFromProfile(target *ConfigService, source *ConfigFileProfile) error {
	if source.BaseURL != "" {
		target.BaseURL = source.BaseURL
	}
	if source.Key != "" {
		target.Key.Value = source.Key
		target.Key.Source = CfConfigFile
	}
	if source.Model != "" {
		target.Model = source.Model
	}

	prompt, err := getServicePrompt(source)
	if err != nil {
		return err
	}

	if prompt != "" {
		target.Prompt = prompt
	}

	if source.Concurrency > 0 {
		target.Concurrency = source.Concurrency
	}

	if source.CompactionToks > 0 {
		target.CompactionToks = source.CompactionToks
	}

	return nil
}

func getServicePrompt(source *ConfigFileProfile) (string, error) {
	if source.Prompt == nil {
		return "", nil
	}

	if source.Prompt.Content != nil {
		return *source.Prompt.Content, nil
	}

	if source.Prompt.Path != nil {
		p := *source.Prompt.Path

		if source.Context.IsGlobal {
			p = fsx.Resolve(source.Context.Dir, p)
		} else {
			if filepath.IsAbs(p) {
				return "", errors.New("prompt path must be relative")
			}

			if fsx.EscapesBaseDir(p) {
				return "", errors.New("prompt file path escapes the root directory")
			}

			p = filepath.Join(source.Context.Dir, p)
		}

		content, err := os.ReadFile(p)
		if err != nil {
			return "", err
		}

		return string(content), nil
	}

	return "", errors.New("neither path nor content are specified for prompt")
}
