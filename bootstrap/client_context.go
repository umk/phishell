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

	app := NewApp(config)

	ctx = context.WithValue(ctx, CtxApp, app)

	return ctx
}

func GetApp(ctx context.Context) *App {
	return ctx.Value(CtxApp).(*App)
}

func GetVersion(ctx context.Context) string {
	return ctx.Value(CtxVersion).(string)
}

// IsDebug gets a value indicating whether debugging is enabled for this session.
func IsDebug(ctx context.Context) bool {
	return GetApp(ctx).Config.Debug
}
