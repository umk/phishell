package bootstrap

import (
	"flag"
	"fmt"
	"os"

	"github.com/umk/phishell/util/flagx"
)

type Config struct {
	Dir     string
	Debug   bool
	Version bool

	Services []*ConfigService
}

type ConfigService struct {
	Profile string

	BaseURL string
	Key     ConfigValue[string]
	Model   string

	Prompt string

	Retries        int
	Concurrency    int
	CompactionToks int
}

type ConfigValue[V any] struct {
	Value  V
	Source ConfigSource
}

type ConfigSource int

const (
	CfNone ConfigSource = iota
	CfConfigFile
	CfKeyChain
	CfEnvironment
)

// LoadConfig reads configuration from flags, environment variables, and config files.
func LoadConfig() (*Config, error) {
	// Define command-line flags
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "Usage: phishell option...\n")
		fmt.Fprint(w, "Options:\n")
		flag.PrintDefaults()
	}

	var serviceProfIds flagx.Strings

	dirFlag := flag.String("dir", "", "base directory (default current directory)")
	flag.Var(&serviceProfIds, "profile", "configuration profile")
	debugFlag := flag.Bool("debug", false, "debug interactions")
	versionFlag := flag.Bool("v", false, "show version and quit")

	// Parse the flags
	flag.Parse()

	// Determine the current directory
	var currentDir string
	if *dirFlag != "" {
		currentDir = *dirFlag
	} else {
		// Fallback to current working directory
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("cannot get working directory: %w", err)
		}
		currentDir = wd
	}

	// Loading config files
	f, err := loadConfigFiles(currentDir)
	if err != nil {
		return nil, err
	}

	// Getting profile for service client
	if len(serviceProfIds) == 0 {
		defaultProfId := "default"

		if f.Default != nil {
			defaultProfId = *f.Default
		}

		serviceProfIds = append(serviceProfIds, defaultProfId)
	}

	// Initialize Config with default values
	config := &Config{
		Dir:     currentDir,
		Debug:   *debugFlag,
		Version: *versionFlag,
	}

	processedIds := make(map[string]bool)

	for _, id := range serviceProfIds {
		if _, ok := processedIds[id]; ok {
			continue
		}

		service := createDefaultService(id)
		config.Services = append(config.Services, service)

		processedIds[id] = true
	}

	loadEnvVars(config)

	for _, p := range config.Services {
		if err := setServiceFromConfigFile(p, f); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func setServiceFromConfigFile(target *ConfigService, source *ConfigFile) error {
	profile, ok := source.Profiles[target.Profile]
	if !ok {
		return fmt.Errorf("profile not found: %s", target.Profile)
	}

	if err := setServiceFromProfileOrPreset(target, profile); err != nil {
		return err
	}

	if target.Key.Value == "" {
		k, err := getServiceKey(target.Profile)
		if err != nil {
			return err
		}

		target.Key = k
	}

	return nil
}

func loadEnvVars(c *Config) {
	for _, p := range c.Services {
		if p.Key.Value == "" {
			if v, ok := os.LookupEnv("PHI_KEY"); ok {
				p.Key.Value = v
				p.Key.Source = CfEnvironment

				os.Unsetenv("PHI_KEY")
			}
		}
	}
}

func getServiceKey(profile string) (ConfigValue[string], error) {
	k, err := GetOrReadKey(profile)
	if err != nil {
		return ConfigValue[string]{}, fmt.Errorf("failed to read secret: %w", err)
	}

	return ConfigValue[string]{
		Source: CfKeyChain,
		Value:  k,
	}, nil
}

func createDefaultService(profile string) *ConfigService {
	return &ConfigService{
		Profile: profile,

		Retries:        5,
		Concurrency:    1,
		CompactionToks: 2000,
	}
}
