package bootstrap

import (
	"fmt"
	"strings"

	"github.com/umk/phishell/util/termx"
	"github.com/zalando/go-keyring"
)

const keyringAppName = "phishell"

func GetKey(profile string) (string, error) {
	return keyring.Get(keyringAppName, profile)
}

func GetOrReadKey(profile string) (string, error) {
	k, err := keyring.Get(keyringAppName, profile)
	if err != nil {
		if err == keyring.ErrNotFound {
			return ReadKeyAndUpdate(profile, true)
		} else {
			termx.Error.Printf("cannot read from keyring: %s\n", err)
			return ReadKey(profile, true)
		}
	}

	return k, nil
}

func ReadKeyAndUpdate(profile string, force bool) (string, error) {
	k, err := ReadKey(profile, force)
	if err != nil {
		return k, err
	}

	if k != "" {
		if err := keyring.Set(keyringAppName, profile, k); err != nil {
			termx.Error.Printf("cannot write to keyring: %s\n", err)
		}
	}

	return k, nil
}

func ReadKey(profile string, force bool) (string, error) {
	for {
		s, err := termx.ReadSecret(fmt.Sprintf("%s key >>> ", profile))
		if err != nil {
			return "", err
		}

		s = strings.TrimSpace(s)

		if !force || s != "" {
			return s, nil
		}
	}
}
