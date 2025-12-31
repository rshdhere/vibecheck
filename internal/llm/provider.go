// Package llm is for providing user a choice to choose an llm of their choice
package llm

import (
	"context"
	"errors"
	"maps"
	"slices"
)

type Provider interface {
	GenerateCommitMessage(
		ctx context.Context,
		diff string,
		additionalContext string,
	) (string, error)
}

var providers map[string]Provider

func init() {
	providers = map[string]Provider{}
}

func Register(name string, provider Provider) {
	providers[name] = provider
}

func GetRegisteredNames() []string {
	return slices.Collect(maps.Keys(providers))
}

var ErrNoProvider = errors.New("no provider for name")

func GetProvider(name string) (Provider, error) {
	x, exists := providers[name]
	if !exists {
		return nil, ErrNoProvider
	}
	return x, nil
}
