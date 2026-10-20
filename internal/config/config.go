// Package config provides configuration management for LazyTF.
// It handles loading, saving, and validating application configuration
// from ~/.config/lazytf/config.yml following the XDG Base Directory specification.
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigPath returns the path to the config file
// Following XDG Base Directory specification: ~/.config/lazytf/config.yml
func ConfigPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homePath, ".config", "lazytf", "config.yml")
	return configPath, nil
}

// Load reads and parses the config file
// Returns DefaultConfig() if file doesn't exist
func Load() (Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}
	return loadFrom(configPath)
}

// LoadFile is Load for an explicit path (the --config flag). Unlike the
// default location, a file that was asked for by name must exist.
func LoadFile(path string) (Config, error) {
	if _, err := os.Stat(path); err != nil {
		return DefaultConfig(), err
	}
	return loadFrom(path)
}

// loadFrom reads the config at path over the defaults: keys the file does not
// set keep their default value, so a file with only search_paths still gets
// the default ignore patterns.
func loadFrom(configPath string) (Config, error) {
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return DefaultConfig(), err
	}
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

// Save writes the config to disk
// Creates the directory if it doesn't exist
func Save(cfg Config) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(configPath)
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if the config file exists
func Exists() bool {
	configPath, err := ConfigPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(configPath)
	return err == nil
}
