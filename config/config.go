package config

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/umk/phishell/util/flagx"
	"github.com/umk/phishell/util/slicesx"
)

var Config struct {
	Dir           string `validate:"required"`
	Debug         bool
	Profiles      []*Profile `validate:"dive"`
	ChatProfiles  []string
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
		fmt.Fprint(w, "Usage: phishell option...\n")
		fmt.Fprint(w, "Options:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&Config.Dir, "dir", "", "base directory (default current directory)")
	flag.Var((*flagx.Strings)(&Config.ChatProfiles), "profile", "configuration profile to use")

	// Parse the flags
	flag.Parse()

	Config.ChatProfiles = slicesx.Unique(Config.ChatProfiles)

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
	if len(Config.ChatProfiles) == 0 {
		defaultProfId := "default"

		if len(f.Profiles) == 1 {
			for profileId := range f.Profiles {
				defaultProfId = profileId
			}
		} else {
			if f.Default != nil {
				defaultProfId = *f.Default
			}

			if len(f.Profiles) == 0 {
				f.Profiles[defaultProfId] = &ConfigFileProfile{}
			}
		}

		Config.ChatProfiles = append(Config.ChatProfiles, defaultProfId)
	}

	// Initialize Config with default values
	Config.OutputBufSize = 512 * 1024

	for id, p := range f.Profiles {
		profile := &Profile{
			Profile:     id,
			Concurrency: 1,
		}
		if err := setProfileFromFileProfileOrPreset(profile, p); err != nil {
			return err
		}

		Config.Profiles = append(Config.Profiles, profile)
	}

	loadEnvVars()

	if err := checkConfig(); err != nil {
		return err
	}

	return nil
}

func checkConfig() error {
	v := validator.New()

	offset := 0
	for i, profile := range Config.Profiles {
		if err := v.Struct(profile); err != nil {
			if vals, ok := err.(validator.ValidationErrors); ok && len(vals) > 0 {
				err = fmt.Errorf("invalid field %s", vals[0].Namespace())
			}
			if slices.Contains(Config.ChatProfiles, profile.Profile) {
				return fmt.Errorf("validation failed for profile %q: %w", profile.Profile, err)
			}

			fmt.Fprintf(os.Stderr, "Warning: profile %q is ignored: %v\n", profile.Profile, err)
			offset++
		} else {
			Config.Profiles[i-offset] = profile
		}
	}

	Config.Profiles = Config.Profiles[:len(Config.Profiles)-offset]

	if err := v.Struct(Config); err != nil {
		return err
	}

	return nil
}

func loadEnvVars() {
	_, Config.Debug = os.LookupEnv("DEBUG")

	if v, ok := os.LookupEnv("PHISHELL_KEY"); ok {
		for _, p := range Config.Profiles {
			if p.Key == "" {
				p.Key = v
			}
		}

		os.Unsetenv("PHISHELL_KEY")
	}
}
