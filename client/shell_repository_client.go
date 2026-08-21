package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/SmartFactory-KL/shellgonaut/create"
	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

const ShellRepositoryPath = "/shells"

type ShellRepositoryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewShellRepositoryClient creates a new Shell Repository Client while validating the baseURL
func NewShellRepositoryClient(baseURL string, opts ...ClientOption) (*ShellRepositoryClient, error) {
	shellRepoBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, ShellRepositoryPath)
	if err != nil {
		return nil, fmt.Errorf("invalid url for shell repository: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client := &ShellRepositoryClient{
		httpClient: httpClient,
		baseURL:    shellRepoBaseURL,
	}

	return client, nil
}

// ---------------------------------------- Description -----------------------------
// GetShellRepositoryDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
func (repoClient *ShellRepositoryClient) GetShellRepositoryDescription() ([]byte, error) {
	resp, err := repoClient.httpClient.Get(
		repoClient.baseURL.JoinPath("/description").String(),
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

// ---------------------------------------- Shell Pages ----------------------------
// GetNextShellPage requests the next page starting from cursor. empty cursor starts from the beginning, limit = 0 means no limit
// however: basyx usually has a limit anyway.
func (repoClient *ShellRepositoryClient) GetNextShellPage(cursor string, limit int) (*PagedResult[types.IAssetAdministrationShell], error) {
	targetURL := repoClient.baseURL.JoinPath(ShellRepositoryPath)

	pagedResult, err := DoPagedGetRequest(repoClient.httpClient, targetURL.String(), cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var typedResult PagedResult[types.IAssetAdministrationShell]
	typedResult.Metadata = pagedResult.Metadata
	typedResult.Result = make([]types.IAssetAdministrationShell, 0, len(pagedResult.Result))

	for _, rawIn := range pagedResult.Result {
		shell, err := create.FromBytes(rawIn, jsonization.AssetAdministrationShellFromJsonable)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Shell: %w", err)
		}
		typedResult.Result = append(typedResult.Result, shell)
	}

	return &typedResult, nil
}

// ---------------------------------------- Shells  --------------------------------
// GetSubmodelJsonable gets a shell in the "jsonable" format (as map[string]any - fit for jsonization)
func (repoClient *ShellRepositoryClient) GetShellJsonable(shellID string) (map[string]any, error) {
	return repoClient.derived(shellID).GetJsonable()
}

// GetShell returns a parsed asset administration shell. errors if shell cannot be parsed
func (repoClient *ShellRepositoryClient) GetShell(shellID string) (types.IAssetAdministrationShell, error) {
	return repoClient.derived(shellID).Get()
}

// UploadShell uploads a shell to the /shells endpoint
func (repoClient *ShellRepositoryClient) UploadShell(shell types.IAssetAdministrationShell) error {
	if shell == nil {
		return fmt.Errorf("shell cannot be nil")
	}

	shellBytes, err := aasEntityToBytes(shell)
	if err != nil {
		return fmt.Errorf("failed to convert shell to bytes: %w", err)
	}

	targetUrl := repoClient.baseURL.JoinPath(ShellRepositoryPath)

	body, err := DoPostRequest(repoClient.httpClient, targetUrl.String(), shellBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateShell updates a shell with PUT /shells/:shellID
func (repoClient *ShellRepositoryClient) UpdateShell(shellID string, shell types.IAssetAdministrationShell) error {
	return repoClient.derived(shellID).Update(shell)
}

// DeleteShellByID deletes a shell using DELETE /shells/:shellID
func (repoClient *ShellRepositoryClient) DeleteShell(shellID string) error {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, ShellRepositoryPath, shellID)
	if err != nil {
		return err
	}

	body, err := DoDeleteRequest(repoClient.httpClient, targetUrl.String())
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// GetNextSubmodelReferencePage gets the next page of submodel references for shell
func (repoClient *ShellRepositoryClient) GetNextSubmodelReferencePage(shellID string, cursor string, limit int) (*PagedResult[types.IReference], error) {
	return repoClient.derived(shellID).GetNextSubmodelReferencePage(cursor, limit)
}

// UploadShellSubmodelReference creates a Reference for the given submodel and uploads it to POST  /shells/:shellID/submodel-refs
func (repoClient *ShellRepositoryClient) UploadShellSubmodelReference(shellID string, submodelReference types.IReference) error {
	return repoClient.derived(shellID).UploadSubmodelReference(submodelReference)
}

// DeleteShellSubmodelReference removes a single reference via DELETE /shells/:shellID/submodel-refs/:submodelID
func (repoClient *ShellRepositoryClient) DeleteShellSubmodelReference(shellID string, submodelID string) error {
	return repoClient.derived(shellID).DeleteSubmodelReference(submodelID)
}

// --------------------- Util ------------------------
func (repoClient *ShellRepositoryClient) derived(shellID string) *ShellServiceClient {
	encodedShellID := toBase64URL(shellID)
	targetURL := repoClient.baseURL.JoinPath(ShellRepositoryPath, encodedShellID)
	return NewDerivedShellServiceClient(repoClient.httpClient, targetURL)
}
