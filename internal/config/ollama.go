package config

const (
	DefaultOllamaURL   = "http://localhost:11434"
	DefaultOllamaModel = "gpt-5-mini"
)

type Ollama struct {
	URL   string `yaml:"url"`
	Model string `yaml:"model"`
}

func (o *Ollama) apply() {
	if o.URL == "" {
		o.URL = DefaultOllamaURL
	}

	if o.Model == "" {
		o.Model = DefaultOllamaModel
	}
}
