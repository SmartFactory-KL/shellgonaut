package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
	"github.com/smartfactory-kl/shellgonaut/convert"
)

type ShellServiceClient struct {
	// Http Client to use for requests
	httpClient *http.Client

	// baseURL to append path to
	baseURL *url.URL

	// if true, no /aas is needed in path
	isDerived bool
}

// NewShellServiceClient creates a standalone shell service client, prefixing paths with /aas
func NewShellServiceClient(baseURL string, opts ...ClientOption) (*ShellServiceClient, error) {
	shellServiceBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, "/aas")
	if err != nil {
		return nil, fmt.Errorf("invalid url for shell service: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client := &ShellServiceClient{
		httpClient: httpClient,
		baseURL:    shellServiceBaseURL,
		isDerived:  false,
	}

	return client, nil
}

// NewDerivedShellServiceClient creates a new shell service client for usage within a RepositoryClient
func NewDerivedShellServiceClient(httpClient *http.Client, baseURL *url.URL) *ShellServiceClient {
	return &ShellServiceClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		isDerived:  true,
	}
}

// ---------------------------------------- Description -----------------------------
// GetShellServiceDescription calls the /description endpoint and returns raw json
func (serviceClient *ShellServiceClient) GetShellServiceDescription() ([]byte, error) {
	resp, err := serviceClient.httpClient.Get(
		serviceClient.baseURL.JoinPath("/description").String(),
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

// ------------------------------ Shell ---------------------------------------
// GetJsonable returns the shell in jsonable format as map[string]any
func (serviceClient *ShellServiceClient) GetJsonable() (map[string]any, error) {
	targetURL := serviceClient.getDerivedTargetUrl()

	body, err := DoGetRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get Shell: %w", err)
	}
	defer body.Close()

	return convert.BodyToJsonable(body)
}

// Get returns the parsed shell
func (serviceClient *ShellServiceClient) Get() (types.IAssetAdministrationShell, error) {
	shellJsonable, err := serviceClient.GetJsonable()
	if err != nil {
		return nil, err
	}

	submodel, err := jsonization.AssetAdministrationShellFromJsonable(shellJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shell json: %w", err)
	}

	return submodel, nil
}

// Update updates the shell with PUT
func (serviceClient *ShellServiceClient) Update(shell types.IAssetAdministrationShell) error {
	if shell == nil {
		return fmt.Errorf("shell cannot be nil")
	}

	shellBytes, err := convert.JsonableTypeToBytes(shell)
	if err != nil {
		return fmt.Errorf("failed to convert shell to bytes: %w", err)
	}

	targetURL := serviceClient.getDerivedTargetUrl()

	body, err := DoPutRequest(serviceClient.httpClient, targetURL.String(), shellBytes)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// ----------------------- Asset Information -----------------------------------
// GetAssetInformation returns only AssetInformation of the shell
func (serviceClient *ShellServiceClient) GetAssetInformation() (types.IAssetInformation, error) {
	targetURL := serviceClient.getDerivedTargetUrl("asset-information")

	body, err := DoGetRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get Shell: %w", err)
	}
	defer body.Close()

	assetInfoJsonable, err := convert.BodyToJsonable(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset information json: %w", err)
	}

	assetInfo, err := jsonization.AssetInformationFromJsonable(assetInfoJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse asset information json: %w", err)
	}

	return assetInfo, nil
}

// UpdateAssetInformation updates asset information of shell using PUT /asset-information
func (serviceClient *ShellServiceClient) UpdateAssetInformation(assetInfo types.IAssetInformation) error {
	if assetInfo == nil {
		return fmt.Errorf("asset information cannot be nil")
	}

	shellBytes, err := convert.JsonableTypeToBytes(assetInfo)
	if err != nil {
		return fmt.Errorf("failed to convert asset information to bytes: %w", err)
	}

	targetURL := serviceClient.getDerivedTargetUrl("asset-information")

	body, err := DoPutRequest(serviceClient.httpClient, targetURL.String(), shellBytes)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// --------------------- Shell Submodel Reference -----------------------------
// GetNextSubmodelReferencePage gets the next page of submodel references for the shell
// limit = 0 means no limit (however: BaSyx has internal limit anyway)
func (serviceClient *ShellServiceClient) GetNextSubmodelReferencePage(cursor string, limit int) (*PagedResult[types.IReference], error) {
	targetURL := serviceClient.getDerivedTargetUrl("/submodel-refs")

	pagedResult, err := DoPagedGetRequest(serviceClient.httpClient, targetURL.String(), cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var typedResult PagedResult[types.IReference]
	typedResult.Metadata = pagedResult.Metadata
	typedResult.Result = make([]types.IReference, 0, len(pagedResult.Result))

	for _, rawIn := range pagedResult.Result {
		submodelReference, err := convert.JsonableTypeFromBytes(rawIn, jsonization.ReferenceFromJsonable)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Submodel Reference: %w", err)
		}
		typedResult.Result = append(typedResult.Result, submodelReference)
	}

	return &typedResult, nil
}

// UploadSubmodelReference uploads a new submodel reference to the shell
func (serviceClient *ShellServiceClient) UploadSubmodelReference(reference types.IReference) error {
	if reference == nil {
		return fmt.Errorf("reference cannot be nil")
	}

	referenceBytes, err := convert.JsonableTypeToBytes(reference)
	if err != nil {
		return fmt.Errorf("failed to convert reference to bytes: %w", err)
	}

	targetURL := serviceClient.getDerivedTargetUrl("submodel-refs")

	body, err := DoPostRequest(serviceClient.httpClient, targetURL.String(), referenceBytes)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// DeleteSubmodelReference deletes the submodel reference to submodelID from the shell
func (serviceClient *ShellServiceClient) DeleteSubmodelReference(submodelID string) error {
	if len(submodelID) == 0 {
		return fmt.Errorf("submodelID cannot be empty")
	}

	encodedSubmodelID := toBase64URL(submodelID)
	targetURL := serviceClient.getDerivedTargetUrl("submodel-refs", encodedSubmodelID)

	body, err := DoDeleteRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// ---------------------------------------- Util -----------------------------
// getDerivedTargetUrl returns the url to use, adding /submodel if isDerived is true
func (serviceClient *ShellServiceClient) getDerivedTargetUrl(pathItems ...string) *url.URL {
	toAppend := make([]string, 0, len(pathItems)+1)

	if !serviceClient.isDerived {
		toAppend = append(toAppend, "aas")
	}

	toAppend = append(toAppend, pathItems...)

	for idx, pathItem := range toAppend {
		toAppend[idx] = url.PathEscape(pathItem)
	}

	targetURL := serviceClient.baseURL.JoinPath(toAppend...)
	return targetURL
}
