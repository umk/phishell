package bootstrap

// globalCtx aggregates global context for the application.
type globalCtx struct {
	config *Config

	clients []*ClientRef
}

func newGlobalCtx(config *Config) *globalCtx {
	app := &globalCtx{
		config: config,
	}

	for _, p := range config.Services {
		c := NewClientRef(p)
		app.clients = append(app.clients, c)
	}

	return app
}
