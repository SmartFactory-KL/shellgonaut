package client

// This file contains the options and helpers to deal with creating http clients

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type ClientOption func(*ClientConfiguration) error

type ClientAuthStyle string

const (
	ClientAuthTokenSource ClientAuthStyle = "token"
	ClientAuthApiKey      ClientAuthStyle = "apiKey"
)

type ClientConfiguration struct {
	timeout time.Duration

	authStyle *ClientAuthStyle

	tokenSource oauth2.TokenSource
	apiKey      string
}

// WithTokenSource adds a OAuth2 based auth to the http client used
func WithTokenSource(tokenSource oauth2.TokenSource) ClientOption {
	return func(clientCfg *ClientConfiguration) error {
		if clientCfg.authStyle != nil {
			return fmt.Errorf("only one auth style can be selected, %s was already in place", *clientCfg.authStyle)
		}

		if tokenSource == nil {
			return fmt.Errorf("token source cannot be nil")
		}

		clientCfg.authStyle = new(ClientAuthTokenSource)
		clientCfg.tokenSource = tokenSource
		return nil
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return func(clientCfg *ClientConfiguration) error {
		if clientCfg.authStyle != nil {
			return fmt.Errorf("only one auth style can be selected, %s was already in place", *clientCfg.authStyle)
		}

		if len(apiKey) == 0 {
			return fmt.Errorf("API key cannot be empty")
		}

		clientCfg.authStyle = new(ClientAuthApiKey)
		clientCfg.apiKey = apiKey

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

	if clientCfg.authStyle != nil {
		switch *clientCfg.authStyle {
		case ClientAuthTokenSource:
			httpRoundTripper = &oauth2.Transport{
				Source: clientCfg.tokenSource,
				Base:   httpRoundTripper,
			}
		case ClientAuthApiKey:
			httpRoundTripper = &apiKeyTransport{
				apiKey: clientCfg.apiKey,
				base:   httpRoundTripper,
			}
		default:
			return nil, fmt.Errorf("unkown auth style, abort")
			// no auth, nothing to do
		}
	}

	httpClient := &http.Client{
		Timeout:   clientCfg.timeout,
		Transport: httpRoundTripper,
	}

	return httpClient, nil
}

// API Key RoundTripper
type apiKeyTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("X-Api-Key", t.apiKey)

	return t.base.RoundTrip(req)
}
