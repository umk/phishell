package config

var presetOpenAI = ConfigFileProfile{
	BaseURL:     "https://api.openai.com/v1/",
	Concurrency: 5,
	ContextSize: 256 * 1024,

	Indexing: ConfigFileIndexing{
		ChunkToks:   1000,
		OverlapToks: 200,
	},
}

var presetOllama = ConfigFileProfile{
	BaseURL:     "http://localhost:11434/v1/",
	Key:         "ollama",
	Concurrency: 1,
	ContextSize: 2 * 1024,

	Indexing: ConfigFileIndexing{
		ChunkToks:   200,
		OverlapToks: 40,
	},
}

var presets = map[string]ConfigFileProfile{
	"openai": presetOpenAI,
	"ollama": presetOllama,
}
