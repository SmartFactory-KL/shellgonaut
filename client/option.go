package client

// This file contains the options and helpers to deal with creating http clients

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type ClientOption func(*ClientConfiguration) error

type ClientConfiguration struct {
	timeout     time.Duration
	tokenSource oauth2.TokenSource
}

// WithTokenSource adds a OAuth2 based auth to the http client used
func WithTokenSource(tokenSource oauth2.TokenSource) ClientOption {
	return func(clientCfg *ClientConfiguration) error {
		if tokenSource == nil {
			return fmt.Errorf("token source cannot be nil")
		}

		clientCfg.tokenSource = tokenSource
		return nil
	}
}

// WithTimeout sets a general timeout to all requests of the http client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(clientCfg *ClientConfiguration) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be greater than zero")
		}

		clientCfg.timeout = timeout
		return nil
	}
}

// createHttpClientFromOptions creates a HttpClient based on default values and all applied opts
func createHttpClientFromOptions(opts ...ClientOption) (*http.Client, error) {
	clientCfg := ClientConfiguration{
		timeout:     2 * time.Minute,
		tokenSource: nil,
	}

	for _, opt := range opts {
		if err := opt(&clientCfg); err != nil {
			return nil, fmt.Errorf("failed to apply http client option: %w", err)
		}
	}

	var httpRoundTripper http.RoundTripper = http.DefaultTransport.(*http.Transport).Clone()

	if clientCfg.tokenSource != nil {
		httpRoundTripper = &oauth2.Transport{
			Source: clientCfg.tokenSource,
			Base:   httpRoundTripper,
		}
	}

	httpClient := &http.Client{
		Timeout:   clientCfg.timeout,
		Transport: httpRoundTripper,
	}

	return httpClient, nil
}
