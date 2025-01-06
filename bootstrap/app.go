package bootstrap

// Aggregates global context for the application.
type App struct {
	Config *Config

	Clients []*ClientRef
}

func NewApp(config *Config) *App {
	app := &App{Config: config}

	for _, p := range config.Services {
		c := NewClientRef(p)
		app.Clients = append(app.Clients, c)
	}

	return app
}

func (a *App) PrimaryClient() *ClientRef {
	if len(a.Clients) == 0 {
		panic("no clients defined for the app")
	}

	return a.Clients[0]
}
