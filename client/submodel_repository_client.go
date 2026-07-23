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

type SubmodelRepositoryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewSubmodelRepositoryClient creates a new Submodel Repository Client while validating baseURL
func NewSubmodelRepositoryClient(baseURL string) (*SubmodelRepositoryClient, error) {
	submodelRepoBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, "/submodels")
	if err != nil {
		return nil, fmt.Errorf("invalid url for submodel repository: %w", err)
	}

	client := &SubmodelRepositoryClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},

		baseURL: submodelRepoBaseURL,
	}

	return client, nil
}

// ---------------------------------------- Desription -----------------------------
// GetDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
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
// TODO: add GetAllSubmodelPages() and GetNextSubmodelPage(cursor)

// ---------------------------------------- Submodels --------------------------------
// GetSubmodelJsonable gets a submodel in the "jsonable" format (as map[string]any - fit for jsonization)
func (repoClient *SubmodelRepositoryClient) GetSubmodelJsonable(submodelID string) (map[string]any, error) {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID)
	if err != nil {
		return nil, err
	}

	body, err := DoGetRequest(repoClient.httpClient, targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get submodel: %w", err)
	}
	defer body.Close()

	return bodyToJsonable(body)
}

// GetSubmodel returns a parsed submodel. errors if submodel cannot be parsed.
func (repoClient *SubmodelRepositoryClient) GetSubmodel(submodelID string) (types.ISubmodel, error) {
	submodelJsonable, err := repoClient.GetSubmodelJsonable(submodelID)
	if err != nil {
		return nil, err
	}

	submodel, err := jsonization.SubmodelFromJsonable(submodelJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse submodel json: %w", err)
	}

	return submodel, nil
}

// GetSubmodelMetadataJsonable gets the metadata representation of a submodel in jsonable form
func (repoClient *SubmodelRepositoryClient) GetSubmodelMetadataJsonable(submodelID string) (map[string]any, error) {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID, "/$metadata")
	if err != nil {
		return nil, err
	}

	body, err := DoGetRequest(repoClient.httpClient, targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get submodel metadata: %w", err)
	}
	defer body.Close()

	return bodyToJsonable(body)
}

// GetSubmodelMetadata gets the metadata only representation of a submodel
func (repoClient *SubmodelRepositoryClient) GetSubmodelMetadata(submodelID string) (types.ISubmodel, error) {
	submodelJsonable, err := repoClient.GetSubmodelMetadataJsonable(submodelID)
	if err != nil {
		return nil, err
	}

	submodel, err := jsonization.SubmodelFromJsonable(submodelJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse submodel metadata json: %w", err)
	}

	return submodel, nil
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

	targetUrl := repoClient.baseURL.JoinPath("/submodels")

	body, err := DoPostRequest(repoClient.httpClient, targetUrl.String(), submodelBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateSubmodel updates a submodel with PUT /submodels/:submodelID
func (repoClient *SubmodelRepositoryClient) UpdateSubmodel(submodel types.ISubmodel) error {
	if submodel == nil {
		return fmt.Errorf("submodel cannot be nil")
	}

	submodelBytes, err := aasEntityToBytes(submodel)
	if err != nil {
		return fmt.Errorf("failed to convert submodel to bytes: %w", err)
	}

	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodel.ID())
	if err != nil {
		return err
	}

	body, err := DoPutRequest(repoClient.httpClient, targetUrl.String(), submodelBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// DeleteSubmodel is a wrapper for DeleteSubmodelByID while using submodel.ID() as submodelID
func (repoClient *SubmodelRepositoryClient) DeleteSubmodel(submodel types.ISubmodel) error {
	if submodel == nil {
		return fmt.Errorf("submodel cannot be nil")
	}
	return repoClient.DeleteSubmodelByID(submodel.ID())
}

// DeleteSubmodelByID deletes a submodel using DELETE /submodels/:submodelID
func (repoClient *SubmodelRepositoryClient) DeleteSubmodelByID(submodelID string) error {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID)
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
// GetSubmodelElementJsonable gets a submodel element from submodelID with idShortPath in the jsonable format
func (repoClient *SubmodelRepositoryClient) GetSubmodelElementJsonable(submodelID string, idShortPath string) (map[string]any, error) {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID, "/submodel-elements", idShortPath)
	if err != nil {
		return nil, err
	}

	body, err := DoGetRequest(repoClient.httpClient, targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get submodel element: %w", err)
	}
	defer body.Close()

	return bodyToJsonable(body)
}

// GetSubmodelElement returns a parse submodel element. errors if submodel element cannot be parsed.
// Returns only a generic submodelElement. Use EnsureSubmodelElementType() on the result to get types
func (repoClient *SubmodelRepositoryClient) GetSubmodelElement(submodelID string, idShortPath string) (types.ISubmodelElement, error) {
	submodelElementJsonable, err := repoClient.GetSubmodelElementJsonable(submodelID, idShortPath)
	if err != nil {
		return nil, err
	}

	submodelElement, err := jsonization.SubmodelElementFromJsonable(submodelElementJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse submodel element json: %w", err)
	}

	return submodelElement, nil
}

// GetSubmodelElementValue returns the valueOnly representation of a submodel element. returns only a string, no typed elements.
func (repoClient *SubmodelRepositoryClient) GetSubmodelElementValue(submodelID string, idShortPath string) (string, error) {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID, "/submodel-elements", idShortPath, "$value")
	if err != nil {
		return "", err
	}

	body, err := DoGetRequest(repoClient.httpClient, targetUrl.String())
	if err != nil {
		return "", fmt.Errorf("failed to get submodel element value: %w", err)
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of submodel element value: %w", err)
	}

	return string(content), nil
}

// UploadSubmodelElement uploads a submodel with POST /submodels/:submodelID/submodel-elements/:idShortPath
func (repoClient *SubmodelRepositoryClient) UploadSubmodelElement(submodelID string, idShortPath string, element types.ISubmodelElement) error {
	if element == nil {
		return fmt.Errorf("element cannot be nil")
	}

	submodelElementBytes, err := aasEntityToBytes(element)
	if err != nil {
		return fmt.Errorf("failed to convert submodel element to bytes: %w", err)
	}

	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID, "/submodel-elements", idShortPath)
	if err != nil {
		return err
	}

	body, err := DoPostRequest(repoClient.httpClient, targetUrl.String(), submodelElementBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateSubmodelElement updates a submodel element with PUT /submodels/:submodelID/submodel-elements/:idShortPath
func (repoClient *SubmodelRepositoryClient) UpdateSubmodelElement(submodelID string, idShortPath string, element types.ISubmodelElement) error {
	if element == nil {
		return fmt.Errorf("element cannot be nil")
	}

	submodelElementBytes, err := aasEntityToBytes(element)
	if err != nil {
		return fmt.Errorf("failed to convert submodel element to bytes: %w", err)
	}

	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID, "/submodel-elements", idShortPath)
	if err != nil {
		return err
	}

	body, err := DoPutRequest(repoClient.httpClient, targetUrl.String(), submodelElementBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateSubmodelElementValue updates the value of a submodel element using PATCH /submodels/:submodelID/submodel-elements/:idShortPath/$value
func (repoClient *SubmodelRepositoryClient) UpdateSubmodelElementValue(submodelID string, idShortPath string, content []byte) error {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID, "/submodel-elements", idShortPath, "$value")
	if err != nil {
		return err
	}

	body, err := DoPatchRequest(repoClient.httpClient, targetUrl.String(), content)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// DeleteSubmodelElement removes a submodel element with DELETE /submodels/:submodelID/submodel-elements/:idShortPath
func (repoClient *SubmodelRepositoryClient) DeleteSubmodelElement(submodelID string, idShortPath string) error {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, "/submodels", submodelID, "/submodel-elements", idShortPath)
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
