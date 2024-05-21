package cli

type ApplicationConfig struct {
	Version           string                   `json:"version" yaml:"version"`
	Copyright         string                   `json:"copyright" yaml:"copyright"`
	ReleasedAt        string                   `json:"released_at" yaml:"released_at"`
	CaseSensitive     bool                     `json:"case_sensitive" yaml:"case_sensitive"`
	CliPrefix         string                   `json:"cli_prefix" yaml:"cli_prefix"`
	PreferredLanguage string                   `json:"preferred_language" yaml:"preferred_language"`
	Commands          map[string]CommandConfig `json:"commands" yaml:"commands"`
	Debug             bool                     `json:"debug" yaml:"debug"`
}

type DescriptionConfig struct {
	Language string `json:"lang" yaml:"lang"`
	Name     string `json:"name" yaml:"name"`
	Text     string `json:"text" yaml:"text"`
}

type CommandConfig struct {
	Type         string                   `json:"type" yaml:"type"`
	Handler      string                   `json:"handler" yaml:"handler"`
	Examples     string                   `json:"examples" yaml:"examples"`
	Descriptions []DescriptionConfig      `json:"descriptions" yaml:"descriptions"`
	Commands     map[string]CommandConfig `json:"commands" yaml:"commands"`
}
