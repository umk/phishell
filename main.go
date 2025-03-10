package main

import (
	"context"
	_ "embed"
	"errors"
	"io"
	"log"
	"os"

	"github.com/umk/phishell/cli"
	"github.com/umk/phishell/client"
	"github.com/umk/phishell/config"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/termx"
)

//go:embed LOGO
var logo string

//go:embed VERSION
var version string

func main() {
	if err := config.Init(); err != nil {
		termx.Error.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

	if !termx.IsInteractive() {
		os.Exit(1)
	}

	if err := client.Init(); err != nil {
		termx.Error.Printf("error creating clients: %v\n", err)
		os.Exit(1)
	}

	printLogo()

	ctx := context.Background()

	if err := runCli(ctx); err != nil {
		if !errors.Is(err, io.EOF) && !errorsx.IsCanceled(err) {
			termx.Error.Println(err)
			os.Exit(1)
		}
	}
}

func printLogo() {
	termx.Logo.Println(logo)
	termx.Logo.Println(version)
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
