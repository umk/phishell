package bootstrap

import (
	"context"
)

type ContextKey string

const (
	CtxApp     ContextKey = "app"
	CtxVersion ContextKey = "version"
)

func NewContext(config *Config) context.Context {
	ctx := context.Background()

	app := newGlobalCtx(config)

	ctx = context.WithValue(ctx, CtxApp, app)

	return ctx
}

func getApp(ctx context.Context) *globalCtx {
	return ctx.Value(CtxApp).(*globalCtx)
}

func GetVersion(ctx context.Context) string {
	return ctx.Value(CtxVersion).(string)
}

// GetClient gets the default client to use outside of the chat context
// where user can pick the client explicitly.
func GetDefaultClient(ctx context.Context) *ClientRef {
	a := getApp(ctx)

	if len(a.clients) == 0 {
		panic("no clients defined for the app")
	}

	return a.clients[0]
}

// GetClients gets the clients in the order the client profiles are
// specified in command line.
func GetClients(ctx context.Context) []*ClientRef {
	return getApp(ctx).clients
}

// GetConfig gets the application configuration.
func GetConfig(ctx context.Context) *Config {
	return getApp(ctx).config
}

// IsDebug gets a value indicating whether debugging is enabled for this session.
func IsDebug(ctx context.Context) bool {
	return getApp(ctx).config.Debug
}
