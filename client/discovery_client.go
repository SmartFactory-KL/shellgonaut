package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
	"github.com/smartfactory-kl/shellgonaut/convert"
)

const DiscoveryPath = "/lookup/shells"

type DiscoveryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

func NewDiscoveryClient(baseURL string, opts ...ClientOption) (*DiscoveryClient, error) {
	discoveryBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, DiscoveryPath)
	if err != nil {
		return nil, fmt.Errorf("invalid url for discovery: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client := &DiscoveryClient{
		httpClient: httpClient,
		baseURL:    discoveryBaseURL,
	}

	return client, nil
}

// ---------------------------------------- Description -----------------------------
// GetDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
func (discClient *DiscoveryClient) GetDiscoveryDescription() ([]byte, error) {
	resp, err := discClient.httpClient.Get(
		discClient.baseURL.JoinPath("/description").String(),
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

// ---------------------------------------- Lookup ----------------------------------

func (discClient *DiscoveryClient) GetNextLookupPage(assetID []types.ISpecificAssetID, cursor string, limit int) (*PagedStringResult, error) {
	targetURL := discClient.baseURL.JoinPath(DiscoveryPath)

	assetBytes, err := convert.JsonableTypeListToBytes(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert assetIDs to bytes: %w", err)
	}

	encodedAssetBytes := toBase64URL(string(assetBytes))

	pagedResult, err := DoPagedGetRequest(
		discClient.httpClient,
		targetURL.String(),
		cursor,
		limit,
		AdditionalHeader{
			Key:   "assetIds",
			Value: encodedAssetBytes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var stringResult PagedStringResult
	stringResult.Metadata = pagedResult.Metadata
	stringResult.Result = make([]string, 0, len(pagedResult.Result))

	for _, rawIn := range pagedResult.Result {
		var strInput string
		if err := json.Unmarshal(rawIn, &strInput); err != nil {
			return nil, fmt.Errorf("failed to parse list of shell IDs: %w", err)
		}
		stringResult.Result = append(stringResult.Result, strInput)
	}

	return &stringResult, nil
}

// ---------------------------------------- ShellID -> AssetID ----------------------
// GetAssetIDListJsonable gets the list of jsonables representing the SpecificAssetIDs that were discovered for shellID
func (discClient *DiscoveryClient) GetAssetIDListJsonable(shellID string) ([]map[string]any, error) {
	targetUrl, err := getEncodedTargetUrl(discClient.baseURL, "/lookup/shells", shellID)
	if err != nil {
		return nil, err
	}

	body, err := DoGetRequest(discClient.httpClient, targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to discover asset ids for shell %s: %w", shellID, err)
	}
	defer body.Close()

	var resultList []map[string]any
	if err := json.NewDecoder(body).Decode(&resultList); err != nil {
		return nil, fmt.Errorf("failed to parse to json: %w", err)
	}

	return resultList, nil
}

// GetAssetIDList returnes a list of parse SpecificAssetIDs that were discovered for shellID
func (discClient *DiscoveryClient) GetAssetIDList(shellID string) ([]types.ISpecificAssetID, error) {
	listJsonable, err := discClient.GetAssetIDListJsonable(shellID)
	if err != nil {
		return nil, err
	}

	if len(listJsonable) == 0 {
		// could be just an empty list, so no error
		return []types.ISpecificAssetID{}, nil
	}

	resultIDList := make([]types.ISpecificAssetID, 0, len(listJsonable))
	for _, itemJsonable := range listJsonable {
		specificAssetID, err := jsonization.SpecificAssetIDFromJsonable(itemJsonable)
		if err != nil {
			// hmm - abort because of one mismatch?
			return nil, fmt.Errorf("failed to convert jsonable into SpecificAssetID: %w", err)
		}

		resultIDList = append(resultIDList, specificAssetID)
	}

	return resultIDList, nil
}

// UploadAssetIDList uploads a new set of SpecificAssetIDs to POST /lookup/shells/:shellID
func (discClient *DiscoveryClient) UploadAssetIDList(shellID string, assetIDList []types.ISpecificAssetID) error {
	if len(assetIDList) == 0 {
		return fmt.Errorf("assetIDList cannot be nil or empty")
	}

	assetListBytes, err := convert.JsonableTypeListToBytes(assetIDList)
	if err != nil {
		return fmt.Errorf("failed to convert asset ids into bytes: %w", err)
	}

	targetUrl, err := getEncodedTargetUrl(discClient.baseURL, "/lookup/shells", shellID)
	if err != nil {
		return err
	}

	body, err := DoPostRequest(discClient.httpClient, targetUrl.String(), assetListBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	// temporary
	return nil
}

// DeleteAssetIDList deletes everything for a shellID by DELETE /lookup/shells/:shellID
func (discClient *DiscoveryClient) DeleteAssetIDList(shellID string) error {
	targetUrl, err := getEncodedTargetUrl(discClient.baseURL, "/lookup/shells", shellID)
	if err != nil {
		return err
	}

	body, err := DoDeleteRequest(discClient.httpClient, targetUrl.String())
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}
