package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

type ShellRepositoryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewShellRepositoryClient creates a new Shell Repository Client while validating the baseURL
func NewShellRepositoryClient(baseURL string) (*ShellRepositoryClient, error) {
	shellRepoBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, "/shells")
	if err != nil {
		return nil, fmt.Errorf("invalid url for shell repository: %w", err)
	}

	client := &ShellRepositoryClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},

		baseURL: shellRepoBaseURL,
	}

	return client, nil
}

// ---------------------------------------- Desription -----------------------------
// GetDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
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
// TODO: add GetAllShellPages() and GetNextShellPage(cursor)

// ---------------------------------------- Shells  --------------------------------
// GetSubmodelJsonable gets a shell in the "jsonable" format (as map[string]any - fit for jsonization)
func (repoClient *ShellRepositoryClient) GetShellJsonable(shellID string) (map[string]any, error) {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/shells", shellID)
	if err != nil {
		return nil, err
	}

	body, err := DoGetRequest(repoClient.httpClient, targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get shell: %w", err)
	}
	defer body.Close()

	return bodyToJsonable(body)
}

// GetShell returns a parsed asset administration shell. errors if shell cannot be parsed
func (repoClient *ShellRepositoryClient) GetShell(shellID string) (types.IAssetAdministrationShell, error) {
	shellJsonable, err := repoClient.GetShellJsonable(shellID)
	if err != nil {
		return nil, err
	}

	shell, err := jsonization.AssetAdministrationShellFromJsonable(shellJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shell json: %w", err)
	}

	return shell, nil
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

	targetUrl := repoClient.baseURL.JoinPath("/shells")

	body, err := DoPostRequest(repoClient.httpClient, targetUrl.String(), shellBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateShell updates a shell with PUT /shells/:shellID
func (repoClient *ShellRepositoryClient) UpdateShell(shell types.IAssetAdministrationShell) error {
	if shell == nil {
		return fmt.Errorf("shell cannot be nil")
	}

	shellBytes, err := aasEntityToBytes(shell)
	if err != nil {
		return fmt.Errorf("failed to convert shell to bytes: %w", err)
	}

	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/shells", shell.ID())
	if err != nil {
		return err
	}

	body, err := DoPutRequest(repoClient.httpClient, targetUrl.String(), shellBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// DeleteShell is a wrapper for DeleteShellByID while using shell.ID() as shellID
func (repoClient *ShellRepositoryClient) DeleteShell(shell types.IAssetAdministrationShell) error {
	if shell == nil {
		return fmt.Errorf("shell cannot be nil")
	}

	return repoClient.DeleteShellByID(shell.ID())
}

// DeleteShellByID deletes a shell using DELETE /shells/:shellID
func (repoClient *ShellRepositoryClient) DeleteShellByID(shellID string) error {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/shells", shellID)
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

// GetShellSubmodelReferencesJsonable calls the /submodel-refs endpoint and returns a list of jsonables for a shell
func (repoClient *ShellRepositoryClient) GetShellSubmodelReferencesJsonable(shellID string) ([]map[string]any, error) {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/shells", shellID, "/submodel-refs")
	if err != nil {
		return nil, err
	}

	body, err := DoGetRequest(repoClient.httpClient, targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get shell: %w", err)
	}
	defer body.Close()

	return bodyToJsonableList(body)
}

// GetShellSubmodelReferences calls the /submodel-refs endpoint but converts the received jsonables to []types.IReference
// fails if any reference is non-jsonizable
func (repoClient *ShellRepositoryClient) GetShellSubmodelReferences(shellID string) ([]types.IReference, error) {
	referenceJsonableList, err := repoClient.GetShellSubmodelReferencesJsonable(shellID)
	if err != nil {
		return nil, err
	}

	list := []types.IReference{}

	for idx, entry := range referenceJsonableList {
		reference, err := jsonization.ReferenceFromJsonable(entry)
		if err != nil {
			return nil, fmt.Errorf("submodel reference #%d of shellID %s failed to jsonize: %w", idx, shellID, err)
		}
		list = append(list, reference)
	}

	return list, nil
}

// UploadShellSubmodelReference creates a Reference for the given submodel and uploads it to POST  /shells/:shellID/submodel-refs
func (repoClient *ShellRepositoryClient) UploadShellSubmodelReference(shellID string, submodel types.ISubmodel) error {
	if submodel == nil {
		return fmt.Errorf("submodel cannot be nil")
	}

	submodelReference, err := createReferenceForSubmodel(submodel)
	if err != nil {
		return fmt.Errorf("failed to create submodel reference for submodel: %w", err)
	}

	referenceBytes, err := aasEntityToBytes(submodelReference)
	if err != nil {
		return fmt.Errorf("failed to convert submodel reference to bytes: %w", err)
	}

	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/shells", shellID, "/submodel-refs")
	if err != nil {
		return err
	}

	body, err := DoPostRequest(repoClient.httpClient, targetUrl.String(), referenceBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// DeleteShellSubmodelReference removes a single reference via DELETE /shells/:shellID/submodel-refs/:submodelID
func (repoClient *ShellRepositoryClient) DeleteShellSubmodelReference(shellID string, submodelID string) error {
	encodedSubmodelID := toBase64URL(submodelID)
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/shells", shellID, "/submodel-refs", encodedSubmodelID)
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
