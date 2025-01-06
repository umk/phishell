package main

import (
	"context"
	_ "embed"

	"github.com/umk/phishell/bootstrap"
)

//go:embed VERSION
var version string

func initContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, bootstrap.CtxVersion, version)

	return ctx
}
