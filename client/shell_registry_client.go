package client

// TODO: Find out how to handle Descriptors and their jsonization

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// )

// const ShellRegistryPath = "/shell-descriptors"

// type ShellRegistryClient struct {
// 	httpClient *http.Client
// 	baseURL    *url.URL
// }

// // NewShellRegistryClient creates a new client for a shell registry
// func NewShellRegistryClient(baseURL string, opts ...ClientOption) (*ShellRegistryClient, error) {
// 	shellRegistryBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, ShellRegistryPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid url for shell registry: %w", err)
// 	}

// 	httpClient, err := createHttpClientFromOptions(opts...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create http client: %w", err)
// 	}

// 	client := &ShellRegistryClient{
// 		httpClient: httpClient,
// 		baseURL:    shellRegistryBaseURL,
// 	}

// 	return client, nil
// }

// // ---------------------------------------- Description -----------------------------
// // GetShellRegistryDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
// func (regClient *ShellRegistryClient) GetShellRegistryDescription() ([]byte, error) {
// 	resp, err := regClient.httpClient.Get(
// 		regClient.baseURL.JoinPath("/description").String(),
// 	)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to GET /description: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("non-200 response for GET /description: %s", resp.Status)
// 	}

// 	var result []byte
// 	result, err = io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response bytes for GET /description: %w", err)
// 	}

// 	return result, nil
// }
