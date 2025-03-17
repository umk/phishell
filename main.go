package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/umk/phishell/cli"
	"github.com/umk/phishell/client"
	"github.com/umk/phishell/config"
	"github.com/umk/phishell/server"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/termx"
)

//go:embed VERSION
var version string

func main() {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		fmt.Fprintln(os.Stderr, version)
	}

	ctx := context.Background()

	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if !termx.IsInteractive() {
		os.Exit(1)
	}

	if err := client.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error creating clients: %v\n", err)
		os.Exit(1)
	}

	if err := server.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to set up server: %v\n", err)
	} else {
		os.Setenv("PHISHELL_PROFILE", client.Default.Config.Profile)
		defer server.Close(ctx)
	}

	if err := runCli(ctx); err != nil {
		if !errors.Is(err, io.EOF) && !errorsx.IsCanceled(err) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func runCli(ctx context.Context) error {
	c := cli.NewCli()

	log.SetOutput(io.Discard)

	if err := c.Init(ctx); err != nil {
		return err
	}

	if err := c.Run(ctx); err != nil {
		return err
	}

	return nil
}
