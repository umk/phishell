package bootstrap

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/umk/phishell/util/flagx"
	"github.com/umk/phishell/util/slicesx"
)

type Config struct {
	Dir           string `validate:"required"`
	Debug         bool
	Version       bool
	Profiles      []*ConfigProfile `validate:"dive"`
	OutputBufSize int
}

type ConfigProfile struct {
	Profile string

	BaseURL string `validate:"required,url"`
	Key     string `validate:"required"`
	Model   string `validate:"required"`

	Prompt string

	Concurrency int `validate:"required,min=1"`
	ContextSize int `validate:"required,min=2000"`

	Dir       string `validate:"required"`
	IsPersona bool
}

type ConfigSource int

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
	var personaIds flagx.Strings

	dirFlag := flag.String("dir", "", "base directory (default current directory)")
	flag.Var(&serviceProfIds, "profile", "configuration profile")
	flag.Var(&personaIds, "persona", "configuration profile as persona")
	debugFlag := flag.Bool("debug", false, "debug interactions")
	versionFlag := flag.Bool("v", false, "show version and quit")

	// Parse the flags
	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	personaIds = slicesx.Unique(personaIds)
	serviceProfIds = slicesx.Unique(append(serviceProfIds, personaIds...))

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

		OutputBufSize: 512 * 1024,
	}

	processedIds := make(map[string]bool)

	for _, id := range serviceProfIds {
		if _, ok := processedIds[id]; ok {
			continue
		}

		profile := createDefaultProfile(id)
		profile.IsPersona = slices.Contains(personaIds, id)

		config.Profiles = append(config.Profiles, profile)

		processedIds[id] = true
	}

	for _, p := range config.Profiles {
		if err := setServiceFromConfigFile(p, f); err != nil {
			return nil, err
		}
	}

	loadEnvVars(config)

	v := validator.New()
	if err := v.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setServiceFromConfigFile(target *ConfigProfile, source *ConfigFile) error {
	profile, ok := source.Profiles[target.Profile]
	if !ok {
		return fmt.Errorf("profile not found: %q", target.Profile)
	}

	if err := setServiceFromProfileOrPreset(target, profile); err != nil {
		return err
	}

	return nil
}

func loadEnvVars(c *Config) {
	for _, p := range c.Profiles {
		if p.Key == "" {
			if v, ok := os.LookupEnv("PHI_KEY"); ok {
				p.Key = v

				os.Unsetenv("PHI_KEY")
			}
		}
	}
}

func createDefaultProfile(profile string) *ConfigProfile {
	return &ConfigProfile{
		Profile: profile,

		Concurrency: 1,
	}
}
