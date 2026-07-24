package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/SmartFactory-KL/shellgonaut/create"
	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

type ShellRegistryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewShellRegistryClient creates a new Shell Registry Client while validating the baseURL
func NewShellRegistryClient(baseURL string) (*ShellRegistryClient, error) {
	shellRegistryBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, "/shell-descriptors")
	if err != nil {
		return nil, fmt.Errorf("invalid url for shell repository: %w", err)
	}

	client := &ShellRegistryClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},

		baseURL: shellRegistryBaseURL,
	}

	return client, nil
}

// ---------------------------------------- Desription -----------------------------
// GetShellRegistryDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
func (regClient *ShellRegistryClient) GetShellRegistryDescription() ([]byte, error) {
	resp, err := regClient.httpClient.Get(
		regClient.baseURL.JoinPath("/description").String(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to GET /description: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response for GET /description: %s", resp.Status)
	}

	var result []byte
	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response bytes for GET /description: %w", err)
	}

	return result, nil
}

// ---------------------------------------- Shell Descriptor Pages ----------------------------
// GetNextShellPage requests the next page starting from cursor. empty cursor starts from the beginning, limit = 0 means no limit
// however: basyx usually has a limit anyway.
func (repoClient *ShellRepositoryClient) GetNextShellDescriptorPage(cursor string, limit int) (*BasyxPagedResult[types.IAssetAdministrationShell], error) {
	targetUrl := repoClient.baseURL.JoinPath("/shell-descriptors")

	request, err := http.NewRequest(http.MethodGet, targetUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request: %w", err)
	}

	if len(cursor) > 0 {
		request.Header.Add("cursor", cursor)
	}

	if limit > 0 {
		request.Header.Add("limit", strconv.FormatInt(int64(limit), 10))
	}

	resp, err := repoClient.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// read page
		var basyxResultRaw BasyxPagedResultRaw
		if err := json.NewDecoder(resp.Body).Decode(&basyxResultRaw); err != nil {
			return nil, fmt.Errorf("request returned %s but decoding failed: %w", resp.Status, err)
		}

		// copy metadata
		var basyxResult BasyxPagedResult[types.IAssetAdministrationShell]
		basyxResult.Metadata = basyxResultRaw.Metadata
		basyxResult.Result = make([]types.IAssetAdministrationShell, 0, len(basyxResultRaw.Result))

		// return early for empty or nil result
		if len(basyxResultRaw.Result) == 0 {
			return &basyxResult, nil
		}

		// convert raw messages to AAS
		for _, rawJsonInput := range basyxResultRaw.Result {
			shell, err := create.FromBytes(rawJsonInput, jsonization.AssetAdministrationShellFromJsonable)
			if err != nil {
				return nil, fmt.Errorf("failed to parse AAS: %w", err)
			}

			basyxResult.Result = append(basyxResult.Result, shell)
		}

		return &basyxResult, nil
	} else {
		// try to read errors
		var errorResult BasyxErrorResult
		if err := json.NewDecoder(resp.Body).Decode(&errorResult); err != nil {
			return nil, fmt.Errorf("request failed with status %s but error could not be decoded: %w", resp.Status, err)
		}

		return nil, errorResult
	}
}
