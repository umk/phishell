package client

import (
	"errors"

	"github.com/umk/phishell/config"
)

var ChatProfiles []*Ref
var Profiles = make(map[string]*Ref)

var Default *Ref

func Init() error {
	if len(config.Config.ChatProfiles) == 0 {
		return errors.New("no chat profiles defined")
	}

	for _, p := range config.Config.Profiles {
		Profiles[p.Profile] = NewRef(p)
	}

	for _, id := range config.Config.ChatProfiles {
		ChatProfiles = append(ChatProfiles, Profiles[id])
	}

	Default = ChatProfiles[0]
	return nil
}
