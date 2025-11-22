// Package keys handles the storage and retrieval of API keys for vibecheck
package keys

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Keys represents all stored API keys
type Keys struct {
	OpenAI     string `json:"openai,omitempty"`
	Gemini     string `json:"gemini,omitempty"`
	Anthropic  string `json:"anthropic,omitempty"`
	Groq       string `json:"groq,omitempty"`
	Grok       string `json:"grok,omitempty"`
	Kimi       string `json:"kimi,omitempty"`
	Qwen       string `json:"qwen,omitempty"`
	DeepSeek   string `json:"deepseek,omitempty"`
	Perplexity string `json:"perplexity,omitempty"`
	OllamaHost string `json:"ollama_host,omitempty"`
}

// ProviderToKeyField maps provider names to their key field names in the Keys struct
var ProviderToKeyField = map[string]string{
	"openai":     "openai",
	"gemini":     "gemini",
	"anthropic":  "anthropic",
	"groq":       "groq",
	"grok":       "grok",
	"kimi":       "kimi",
	"qwen":       "qwen",
	"deepseek":   "deepseek",
	"perplexity": "perplexity",
	"ollama":     "ollama_host",
}

// ProviderToEnvVar maps provider names to their environment variable names
var ProviderToEnvVar = map[string]string{
	"openai":     "OPENAI_API_KEY",
	"gemini":     "GEMINI_API_KEY",
	"anthropic":  "ANTHROPIC_API_KEY",
	"groq":       "GROQ_API_KEY",
	"grok":       "XAI_API_KEY",
	"kimi":       "MOONSHOT_API_KEY",
	"qwen":       "QWEN_API_KEY",
	"deepseek":   "DEEPSEEK_API_KEY",
	"perplexity": "PERPLEXITY_API_KEY",
	"ollama":     "OLLAMA_HOST",
}

// getKeysPath returns the path to the keys file
func getKeysPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".vibecheck_keys.json"), nil
}

// Load reads the keys from disk
func Load() (*Keys, error) {
	path, err := getKeysPath()
	if err != nil {
		return nil, err
	}

	// If keys file doesn't exist, return empty keys
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &Keys{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var keys Keys
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}

	return &keys, nil
}

// Save writes the keys to disk
func Save(keys *Keys) error {
	path, err := getKeysPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600) // 0600 = read/write for owner only (security)
}

// GetAPIKey retrieves the API key for a provider with priority:
// 1. .env files (loaded by godotenv into environment variables)
// 2. Environment variables (export commands)
// 3. vibecheck keys file
func GetAPIKey(provider string) (string, bool) {
	// First check environment variables (includes .env files loaded by godotenv and export commands)
	envVar, ok := ProviderToEnvVar[provider]
	if ok {
		key, exists := os.LookupEnv(envVar)
		if exists && key != "" {
			return key, true
		}
	}

	// Fall back to keys file
	keys, err := Load()
	if err == nil {
		keyField, ok := ProviderToKeyField[provider]
		if ok {
			var key string
			switch keyField {
			case "openai":
				key = keys.OpenAI
			case "gemini":
				key = keys.Gemini
			case "anthropic":
				key = keys.Anthropic
			case "groq":
				key = keys.Groq
			case "grok":
				key = keys.Grok
			case "kimi":
				key = keys.Kimi
			case "qwen":
				key = keys.Qwen
			case "deepseek":
				key = keys.DeepSeek
			case "perplexity":
				key = keys.Perplexity
			case "ollama_host":
				key = keys.OllamaHost
			}
			if key != "" {
				return key, true
			}
		}
	}

	return "", false
}

// SetAPIKey sets the API key for a provider in the keys file
func SetAPIKey(provider, key string) error {
	keys, err := Load()
	if err != nil {
		keys = &Keys{}
	}

	keyField, ok := ProviderToKeyField[provider]
	if !ok {
		return errors.New("unknown provider")
	}

	switch keyField {
	case "openai":
		keys.OpenAI = key
	case "gemini":
		keys.Gemini = key
	case "anthropic":
		keys.Anthropic = key
	case "groq":
		keys.Groq = key
	case "grok":
		keys.Grok = key
	case "kimi":
		keys.Kimi = key
	case "qwen":
		keys.Qwen = key
	case "deepseek":
		keys.DeepSeek = key
	case "perplexity":
		keys.Perplexity = key
	case "ollama_host":
		keys.OllamaHost = key
	}

	return Save(keys)
}

// GetAllKeys returns all stored keys (for display purposes, keys will be masked)
func GetAllKeys() (map[string]string, error) {
	keys, err := Load()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	if keys.OpenAI != "" {
		result["openai"] = maskKey(keys.OpenAI)
	}
	if keys.Gemini != "" {
		result["gemini"] = maskKey(keys.Gemini)
	}
	if keys.Anthropic != "" {
		result["anthropic"] = maskKey(keys.Anthropic)
	}
	if keys.Groq != "" {
		result["groq"] = maskKey(keys.Groq)
	}
	if keys.Grok != "" {
		result["grok"] = maskKey(keys.Grok)
	}
	if keys.Kimi != "" {
		result["kimi"] = maskKey(keys.Kimi)
	}
	if keys.Qwen != "" {
		result["qwen"] = maskKey(keys.Qwen)
	}
	if keys.DeepSeek != "" {
		result["deepseek"] = maskKey(keys.DeepSeek)
	}
	if keys.Perplexity != "" {
		result["perplexity"] = maskKey(keys.Perplexity)
	}
	if keys.OllamaHost != "" {
		result["ollama"] = keys.OllamaHost // Don't mask host
	}

	return result, nil
}

// maskKey masks most of the key, showing only first 4 and last 4 characters
func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
