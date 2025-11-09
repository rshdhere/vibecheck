// Package config handles the configuration storage for vibecheck
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultProvider string `json:"default_provider"`
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".vibecheck.json"), nil
}

// Load reads the configuration from disk
func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return default config
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &Config{DefaultProvider: "openai"}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func Save(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetDefaultProvider returns the default provider from config
func GetDefaultProvider() string {
	cfg, err := Load()
	if err != nil {
		return "openai" // fallback to openai if error
	}
	return cfg.DefaultProvider
}

// SetDefaultProvider saves the default provider to config
func SetDefaultProvider(provider string) error {
	cfg := &Config{DefaultProvider: provider}
	return Save(cfg)
}
