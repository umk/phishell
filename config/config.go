package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/umk/phishell/util/flagx"
)

var Config struct {
	Dir           string `validate:"required"`
	Debug         bool
	Profiles      []*Profile `validate:"dive"`
	OutputBufSize int
}

type Profile struct {
	Profile string

	BaseURL string `validate:"required,url"`
	Key     string `validate:"required"`
	Model   string `validate:"required"`

	Prompt string

	Concurrency int `validate:"required,min=1"`
	ContextSize int `validate:"required,min=2000"`
}

// Init reads configuration from flags, environment variables, and config files.
func Init() error {
	// Define command-line flags
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "Usage: hackkd option...\n")
		fmt.Fprint(w, "Options:\n")
		flag.PrintDefaults()
	}

	var serviceProfIds flagx.Strings

	flag.StringVar(&Config.Dir, "dir", "", "base directory (default current directory)")
	flag.Var(&serviceProfIds, "profile", "configuration profile to use")

	// Parse the flags
	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Determine the current directory
	if Config.Dir == "" {
		// Fallback to current working directory
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get working directory: %w", err)
		}
		Config.Dir = wd
	}

	// Loading config files
	f, err := loadConfigFiles(Config.Dir)
	if err != nil {
		return err
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
	Config.OutputBufSize = 512 * 1024
	_, Config.Debug = os.LookupEnv("DEBUG")

	processedIds := make(map[string]bool)

	for _, id := range serviceProfIds {
		if _, ok := processedIds[id]; ok {
			continue
		}

		profile := createDefaultService(id)
		Config.Profiles = append(Config.Profiles, profile)

		processedIds[id] = true
	}

	for _, p := range Config.Profiles {
		if err := setServiceFromConfigFile(p, f); err != nil {
			return err
		}
	}

	loadEnvVars()

	v := validator.New()
	if err := v.Struct(Config); err != nil {
		return err
	}

	return nil
}

func setServiceFromConfigFile(target *Profile, source *ConfigFile) error {
	profile, ok := source.Profiles[target.Profile]
	if !ok {
		return fmt.Errorf("profile not found: %q", target.Profile)
	}

	if err := setServiceFromProfileOrPreset(target, profile); err != nil {
		return err
	}

	return nil
}

func loadEnvVars() {
	if v, ok := os.LookupEnv("PHI_SHELL_KEY"); ok {
		for _, p := range Config.Profiles {
			if p.Key == "" {
				p.Key = v
			}
		}

		os.Unsetenv("PHI_SHELL_KEY")
	}
}

func createDefaultService(profile string) *Profile {
	return &Profile{
		Profile: profile,

		Concurrency: 1,
	}
}
