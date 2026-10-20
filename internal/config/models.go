package config

type Config struct {
	SearchPaths    []string `yaml:"search_paths"`
	IgnorePatterns []string `yaml:"ignore_patterns,omitempty"`
	// Layout is the default for every project; a project's .lazytf.yml
	// overrides it. Its keys sit at the top level of the file.
	Layout `yaml:",inline"`
}

func DefaultConfig() Config {
	return Config{
		Layout: DefaultLayout(),
		// Only the working directory: LazyTF never scans home folders the
		// user did not name.
		SearchPaths: []string{"."},
		IgnorePatterns: []string{
			"*/node_modules",
			"*/.git",
			"*/vendor",
			"*/.terraform",
		},
	}
}
