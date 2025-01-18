package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/termx"
)

func main() {
	if !termx.IsInteractive() {
		os.Exit(1)
	}

	config, err := bootstrap.LoadConfig()
	if err != nil {
		termx.Error.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

	if config.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	ctx := bootstrap.NewContext(config)
	ctx = initContext(ctx)

	c := cli.NewCli(bootstrap.IsDebug(ctx))

	if err := c.Init(ctx); err != nil {
		if errors.Is(err, io.EOF) || errorsx.IsCanceled(err) {
			return
		}

		termx.Error.Printf("init error: %v\n", err)
		os.Exit(1)
	}

	if err := c.Run(ctx); err != nil {
		if errors.Is(err, io.EOF) || errorsx.IsCanceled(err) {
			return
		}

		termx.Error.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
