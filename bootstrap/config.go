package bootstrap

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
	Version       bool
	Services      []*ConfigService `validate:"dive"`
	OutputBufSize int
}

type ConfigService struct {
	Profile string

	BaseURL string `validate:"required,url"`
	Key     string `validate:"required"`
	Model   string `validate:"required"`

	Prompt string

	Concurrency int `validate:"required,min=1"`
	ContextSize int `validate:"required,min=2000"`
}

type ConfigSource int

// InitConfig reads configuration from flags, environment variables, and config files.
func InitConfig() error {
	// Define command-line flags
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "Usage: hackkd option...\n")
		fmt.Fprint(w, "Options:\n")
		flag.PrintDefaults()
	}

	var serviceProfIds flagx.Strings

	flag.StringVar(&Config.Dir, "dir", "", "base directory (default current directory)")
	flag.Var(&serviceProfIds, "profile", "configuration profile")
	flag.BoolVar(&Config.Debug, "debug", false, "debug interactions")
	flag.BoolVar(&Config.Version, "v", false, "show version and quit")

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

	processedIds := make(map[string]bool)

	for _, id := range serviceProfIds {
		if _, ok := processedIds[id]; ok {
			continue
		}

		service := createDefaultService(id)
		Config.Services = append(Config.Services, service)

		processedIds[id] = true
	}

	for _, p := range Config.Services {
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

func setServiceFromConfigFile(target *ConfigService, source *ConfigFile) error {
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
	for _, p := range Config.Services {
		if p.Key == "" {
			if v, ok := os.LookupEnv("PHI_KEY"); ok {
				p.Key = v

				os.Unsetenv("PHI_KEY")
			}
		}
	}
}

func createDefaultService(profile string) *ConfigService {
	return &ConfigService{
		Profile: profile,

		Concurrency: 1,
	}
}
