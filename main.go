package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/termx"
)

//go:embed VERSION
var version string

func main() {
	if !termx.IsInteractive() {
		os.Exit(1)
	}

	if err := bootstrap.InitConfig(); err != nil {
		termx.Error.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

	bootstrap.InitClients()

	if bootstrap.Config.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	ctx := context.Background()

	if err := runCli(ctx); err != nil {
		if !errors.Is(err, io.EOF) && !errorsx.IsCanceled(err) {
			termx.Error.Println(err)
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
