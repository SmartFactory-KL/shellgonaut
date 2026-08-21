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

const SubmodelRepositoryPath = "/submodels"

type SubmodelRepositoryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewSubmodelRepositoryClient creates a new Submodel Repository Client while validating baseURL
func NewSubmodelRepositoryClient(baseURL string, opts ...ClientOption) (*SubmodelRepositoryClient, error) {
	submodelRepoBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, SubmodelRepositoryPath)
	if err != nil {
		return nil, fmt.Errorf("invalid url for submodel repository: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client := &SubmodelRepositoryClient{
		httpClient: httpClient,
		baseURL:    submodelRepoBaseURL,
	}

	return client, nil
}

// ---------------------------------------- Description -----------------------------
// GetSubmodelRepositoryDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
func (repoClient *SubmodelRepositoryClient) GetSubmodelRepositoryDescription() ([]byte, error) {
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

// ---------------------------------------- Submodel Pages ---------------------------
// GetNextSubmodelPage requests the next page starting from cursor. empty cursor starts from the beginning, limit = 0 means no limit
// however: basyx usually has a limit anyway.
func (repoClient *SubmodelRepositoryClient) GetNextSubmodelPage(cursor string, limit int) (*PagedResult[types.ISubmodel], error) {
	targetURL := repoClient.baseURL.JoinPath(SubmodelRepositoryPath)

	pagedResult, err := DoPagedGetRequest(repoClient.httpClient, targetURL.String(), cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var typedResult PagedResult[types.ISubmodel]
	typedResult.Metadata = pagedResult.Metadata
	typedResult.Result = make([]types.ISubmodel, 0, len(pagedResult.Result))

	for _, rawIn := range pagedResult.Result {
		submodel, err := create.FromBytes(rawIn, jsonization.SubmodelFromJsonable)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Shell: %w", err)
		}
		typedResult.Result = append(typedResult.Result, submodel)
	}

	return &typedResult, nil
}

// ---------------------------------------- Submodels --------------------------------
// GetSubmodelJsonable gets a submodel in the "jsonable" format (as map[string]any - fit for jsonization)
func (repoClient *SubmodelRepositoryClient) GetSubmodelJsonable(submodelID string) (map[string]any, error) {
	return repoClient.derived(submodelID).GetJsonable()
}

// GetSubmodel returns a parsed submodel. errors if submodel cannot be parsed.
func (repoClient *SubmodelRepositoryClient) GetSubmodel(submodelID string) (types.ISubmodel, error) {
	return repoClient.derived(submodelID).Get()
}

// GetSubmodelMetadataJsonable gets the metadata representation of a submodel in jsonable form
func (repoClient *SubmodelRepositoryClient) GetSubmodelMetadataJsonable(submodelID string) (map[string]any, error) {
	return repoClient.derived(submodelID).GetMetadataJsonable()
}

// GetSubmodelMetadata gets the metadata only representation of a submodel
func (repoClient *SubmodelRepositoryClient) GetSubmodelMetadata(submodelID string) (types.ISubmodel, error) {
	return repoClient.derived(submodelID).GetMetadata()
}

// GetSubmodelValueOnly returns ValueOnly representation of submodel
func (repoClient *SubmodelRepositoryClient) GetSubmodelValueOnly(submodelID string) ([]byte, error) {
	return repoClient.derived(submodelID).GetValueOnly()
}

// UploadSubmodel uploads a submodel to the /submodels endpoint
func (repoClient *SubmodelRepositoryClient) UploadSubmodel(submodel types.ISubmodel) error {
	if submodel == nil {
		return fmt.Errorf("submodel cannot be nil")
	}

	submodelBytes, err := aasEntityToBytes(submodel)
	if err != nil {
		return fmt.Errorf("failed to convert submodel to bytes: %w", err)
	}

	targetUrl := repoClient.baseURL.JoinPath(SubmodelRepositoryPath)

	body, err := DoPostRequest(repoClient.httpClient, targetUrl.String(), submodelBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateSubmodel updates a submodel with PUT /submodels/:submodelID
func (repoClient *SubmodelRepositoryClient) UpdateSubmodel(submodelID string, submodel types.ISubmodel) error {
	return repoClient.derived(submodelID).Update(submodel)
}

// UpdateSubmodelValueOnly updates a Submodel using ValueOnly mode
func (repoClient *SubmodelRepositoryClient) UpdateSubmodelValueOnly(submodelID string, content []byte) error {
	return repoClient.derived(submodelID).UpdateValueOnly(content)
}

// DeleteSubmodelByID deletes a submodel using DELETE /submodels/:submodelID
func (repoClient *SubmodelRepositoryClient) DeleteSubmodel(submodelID string) error {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, SubmodelRepositoryPath, submodelID)
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

// ---------------------------------------- Submodel Elements --------------------------------
func (repoClient *SubmodelRepositoryClient) GetNextSubmodelElementsPage(submodelID string, cursor string, limit int) (*PagedResult[types.ISubmodelElement], error) {
	return repoClient.derived(submodelID).GetNextSubmodelElementsPage(cursor, limit)
}

// GetSubmodelElementJsonable gets a submodel element from submodelID with idShortPath in the jsonable format
func (repoClient *SubmodelRepositoryClient) GetSubmodelElementJsonable(submodelID string, idShortPath string) (map[string]any, error) {
	return repoClient.derived(submodelID).GetSubmodelElementJsonable(idShortPath)
}

// GetSubmodelElement returns a parse submodel element. errors if submodel element cannot be parsed.
// Returns only a generic submodelElement. Use EnsureSubmodelElementType() on the result to get types
func (repoClient *SubmodelRepositoryClient) GetSubmodelElement(submodelID string, idShortPath string) (types.ISubmodelElement, error) {
	return repoClient.derived(submodelID).GetSubmodelElement(idShortPath)
}

// GetSubmodelElementValue returns the valueOnly representation of a submodel element. returns only a string, no typed elements.
func (repoClient *SubmodelRepositoryClient) GetSubmodelElementValue(submodelID string, idShortPath string) ([]byte, error) {
	return repoClient.derived(submodelID).GetSubmodelElementValueOnly(idShortPath)
}

// UploadSubmodelElement uploads a submodel with POST /submodels/:submodelID/submodel-elements/:idShortPath
func (repoClient *SubmodelRepositoryClient) UploadSubmodelElement(submodelID string, idShortPath string, element types.ISubmodelElement) error {
	return repoClient.derived(submodelID).UploadSubmodelElement(idShortPath, element)
}

// UpdateSubmodelElement updates a submodel element with PUT /submodels/:submodelID/submodel-elements/:idShortPath
func (repoClient *SubmodelRepositoryClient) UpdateSubmodelElement(submodelID string, idShortPath string, element types.ISubmodelElement) error {
	return repoClient.derived(submodelID).UpdateSubmodelElement(idShortPath, element)
}

// UpdateSubmodelElementValue updates the value of a submodel element using PATCH /submodels/:submodelID/submodel-elements/:idShortPath/$value
func (repoClient *SubmodelRepositoryClient) UpdateSubmodelElementValueOnly(submodelID string, idShortPath string, content []byte) error {
	return repoClient.derived(submodelID).UpdateSubmodelElementValueOnly(idShortPath, content)
}

// DeleteSubmodelElement removes a submodel element with DELETE /submodels/:submodelID/submodel-elements/:idShortPath
func (repoClient *SubmodelRepositoryClient) DeleteSubmodelElement(submodelID string, idShortPath string) error {
	return repoClient.derived(submodelID).DeleteSubmodelElement(idShortPath)
}

// InvokeOperation invokes a synchronous operation on idShortPath using operationRequest as input
func (repoClient *SubmodelRepositoryClient) InvokeOperation(submodelID string, idShortPath string, operationRequest *OperationRequest) (*OperationResult, error) {
	return repoClient.derived(submodelID).InvokeOperation(idShortPath, operationRequest)
}

// --------------------- Util ------------------------
func (repoClient *SubmodelRepositoryClient) derived(submodelID string) *SubmodelServiceClient {
	encodedSubmodelID := toBase64URL(submodelID)
	targetURL := repoClient.baseURL.JoinPath(SubmodelRepositoryPath, encodedSubmodelID)
	return NewDerivedSubmodelServiceClient(repoClient.httpClient, targetURL)
}
