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

func main() {
	if !termx.IsInteractive() {
		os.Exit(1)
	}

	logCleanup, err := setupLogging()
	if err != nil {
		termx.Error.Printf("unable to create log file: %v\n", err)
	}
	defer logCleanup()

	config, err := bootstrap.LoadConfig()
	if err != nil {
		termx.Error.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

	if config.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	ctx := setupContext(config)

	if err := runCli(ctx); err != nil {
		if !errors.Is(err, io.EOF) && !errorsx.IsCanceled(err) {
			termx.Error.Println(err)
			os.Exit(1)
		}
	}
}

func runCli(ctx context.Context) error {
	c := cli.NewCli(bootstrap.IsDebug(ctx))

	if err := c.Init(ctx); err != nil {
		return err
	}

	if err := c.Run(ctx); err != nil {
		return err
	}

	return nil
}

func setupLogging() (func(), error) {
	logName := fmt.Sprintf("phishell.%d.log", os.Getpid())
	f, err := os.CreateTemp("", logName)
	if err != nil {
		return func() {}, err
	}

	log.SetOutput(f)
	cleanup := func() {
		f.Close()
		os.Remove(f.Name())
	}

	return cleanup, nil
}

func setupContext(config *bootstrap.Config) context.Context {
	ctx := bootstrap.NewContext(config)
	return initContext(ctx)
}
