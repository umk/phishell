package bootstrap

var presetOpenAI = ConfigFileProfile{
	Concurrency:    5,
	CompactionToks: 5000,
}

var presetOllama = ConfigFileProfile{
	BaseURL:        "http://localhost:11434/v1/",
	Key:            "ollama",
	Concurrency:    1,
	CompactionToks: 1800,
}

var presets = map[string]ConfigFileProfile{
	"openai": presetOpenAI,
	"ollama": presetOllama,
}
