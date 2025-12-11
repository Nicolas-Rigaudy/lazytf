package config

type Config struct {
	SearchPaths    []string `yaml:"search_paths"`
	IgnorePatterns []string `yaml:"ignore_patterns,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		SearchPaths: []string{
			".",
			"~/Projects",
			"~/Documents",
		},
		IgnorePatterns: []string{
			"*/node_modules",
			"*/.git",
			"*/vendor",
			"*/.terraform",
		},
	}
}
