package bootstrap

// globalCtx aggregates global context for the application.
type globalCtx struct {
	config *Config

	clients []*Profile
}

func newGlobalCtx(config *Config) *globalCtx {
	app := &globalCtx{
		config: config,
	}

	for _, p := range config.Profiles {
		c := NewProfile(p)
		app.clients = append(app.clients, c)
	}

	return app
}
