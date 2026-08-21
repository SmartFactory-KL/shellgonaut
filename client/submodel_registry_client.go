package client

import (
	"fmt"
	"net/http"
	"net/url"
)

const SubmodelRegistryPath = "/submodel-descriptors"

type SubmodelRegistryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewSubmodelRegistryClient creates a new client for a submodel registry
func NewSubmodelRegistryClient(baseURL string, opts ...ClientOption) (*SubmodelRegistryClient, error) {
	submodelRegistryBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, SubmodelRegistryPath)
	if err != nil {
		return nil, fmt.Errorf("invalid url for submodel registry: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client := &SubmodelRegistryClient{
		httpClient: httpClient,
		baseURL:    submodelRegistryBaseURL,
	}

	return client, nil
}
