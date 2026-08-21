package client

import (
	"fmt"
	"net/http"
	"net/url"
)

const ShellRegistryPath = "/shell-descriptors"

type ShellRegistryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewShellRegistryClient creates a new client for a shell registry
func NewShellRegistryClient(baseURL string, opts ...ClientOption) (*ShellRegistryClient, error) {
	shellRegistryBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, ShellRegistryPath)
	if err != nil {
		return nil, fmt.Errorf("invalid url for shell registry: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client := &ShellRegistryClient{
		httpClient: httpClient,
		baseURL:    shellRegistryBaseURL,
	}

	return client, nil
}
