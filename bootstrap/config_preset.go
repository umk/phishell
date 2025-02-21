package bootstrap

var presetOpenAI = ConfigFileProfile{
	Concurrency: 5,
	ContextSize: 256 * 1024,
}

var presetOllama = ConfigFileProfile{
	BaseURL:     "http://localhost:11434/v1/",
	Key:         "ollama",
	Concurrency: 1,
	ContextSize: 2 * 1024,
}

var presets = map[string]ConfigFileProfile{
	"openai": presetOpenAI,
	"ollama": presetOllama,
}
