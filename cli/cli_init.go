package cli

import (
	"context"
	"os"
	"sync"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/persona"
)

func (c *Cli) Init(ctx context.Context) error {
	config := bootstrap.GetConfig(ctx)

	if err := os.Chdir(config.Dir); err != nil {
		return err
	}

	if err := c.initPersonas(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Cli) initPersonas(ctx context.Context) error {
	personas := make(map[*bootstrap.Profile]*persona.Persona)

	for _, profile := range bootstrap.GetProfiles(ctx) {
		if profile.Config.IsPersona {
			p := persona.New(profile)
			if err := p.Init(ctx); err != nil {
				return err
			}

			personas[profile] = p
		}
	}

	c.personas = personas

	return nil
}

func (c *Cli) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, p := range c.personas {
		wg.Add(1)
		go func(p *persona.Persona) {
			defer wg.Done()
			p.Shutdown(ctx)
		}(p)
	}

	wg.Wait()
	return nil
}
