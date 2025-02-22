package persona

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/util/fsx"
)

type Persona struct {
	profile *bootstrap.Profile
}

func New(profile *bootstrap.Profile) *Persona {
	return &Persona{profile: profile}
}

func (p *Persona) Init(ctx context.Context) error {
	norm := fsx.Normalize(p.profile.Config.Profile)
	pdir := filepath.Join(p.profile.Config.Dir, norm)

	if err := os.MkdirAll(pdir, 0644); err != nil {
		return fmt.Errorf("cannot create persona directory: %w", err)
	}

	if err := createPersonaFiles(pdir); err != nil {
		return err
	}

	return nil
}

func (p *Persona) Post(ctx context.Context, content string) error {
	return nil
}
