package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli"
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
		termx.Error.Printf("init error: %v\n", err)
		os.Exit(1)
	}

	if err := c.Run(ctx); err != nil {
		termx.Error.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
