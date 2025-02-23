package persona

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/umk/phishell/api"
	"github.com/umk/phishell/tool/host/provider"
	"github.com/umk/phishell/util/execx"
	"github.com/umk/phishell/util/fsx"
)

func (p *Persona) Init(ctx context.Context) error {
	norm := fsx.Normalize(p.profile.Config.Profile)
	jsDir := filepath.Join(p.profile.Config.Dir, norm)

	if err := os.MkdirAll(jsDir, 0644); err != nil {
		return fmt.Errorf("cannot create persona directory: %w", err)
	}

	if err := ensurePersonaFiles(jsDir); err != nil {
		return err
	}

	sock, server, err := api.Serve(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	p.server = server

	provider, err := initProvider(jsDir, sock)
	if err != nil {
		server.Close()
		return fmt.Errorf("failed to start provider: %w", err)
	}
	p.provider = provider

	return nil
}

func initProvider(jsDir string, sock string) (*provider.Provider, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("unable to determine location of executable: %w", err)
	}

	corePath := filepath.Join(filepath.Dir(execPath), "persona.js")
	return provider.Start(&execx.Cmd{
		Cmd:  "node",
		Args: []string{corePath, "--", jsDir},
		Env:  append(os.Environ(), "PHI_SHELL=1", "OPENAI_API_BASE="+sock, "OPENAI_API_KEY=1"),
	})
}

func (p *Persona) Shutdown(ctx context.Context) error {
	p.provider.Terminate(provider.PsCompleted, "closed")
	done := make(chan struct{})
	go func() {
		p.provider.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}

	p.server.Shutdown(ctx)

	if p.provider.Info.Status == provider.PsFailed {
		return fmt.Errorf("provider shutdown failed: %s", p.provider.Info.StatusMessage)
	}
	if p.provider.Info.Status != provider.PsCompleted {
		return errors.New("failed to shut down provider in time")
	}

	return nil
}
