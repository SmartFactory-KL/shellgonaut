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

const ConceptDescriptionRepositoryPath = "/concept-descriptions"

type ConceptDescriptionRepositoryClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

func NewConceptDescriptionRepositoryClient(baseURL string, opts ...ClientOption) (*ConceptDescriptionRepositoryClient, error) {
	cdRepoBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, ConceptDescriptionRepositoryPath)
	if err != nil {
		return nil, fmt.Errorf("invalid url for concept description repository: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client for concept description repository: %w", err)
	}

	client := &ConceptDescriptionRepositoryClient{
		httpClient: httpClient,
		baseURL:    cdRepoBaseURL,
	}

	return client, nil
}

// ---------------------------------------- Desription -----------------------------
// GetConceptDescriptionRepositoryDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
func (repoClient *ConceptDescriptionRepositoryClient) GetConceptDescriptionRepositoryDescription() ([]byte, error) {
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

// ------------------------------ Concept Descriptions -----------------------------------
// GetNextConceptDescriptionPage requests the next page of concept descriptions.
// limit = 0 means no limit (Note: BaSyx might have interal limit anyway)
func (repoClient *ConceptDescriptionRepositoryClient) GetNextConceptDescriptionPage(cursor string, limit int) (*PagedResult[types.IConceptDescription], error) {
	targetURL := repoClient.baseURL.JoinPath(ConceptDescriptionRepositoryPath)

	pagedResult, err := DoPagedGetRequest(repoClient.httpClient, targetURL.String(), cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var typedResult PagedResult[types.IConceptDescription]
	typedResult.Metadata = pagedResult.Metadata
	typedResult.Result = make([]types.IConceptDescription, 0, len(pagedResult.Result))

	for _, rawIn := range pagedResult.Result {
		cdItem, err := create.FromBytes(rawIn, jsonization.ConceptDescriptionFromJsonable)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Concept Description: %w", err)
		}
		typedResult.Result = append(typedResult.Result, cdItem)
	}

	return &typedResult, nil
}

// GetConceptDescriptionJsonable gets a concept description in its jsonable format
func (repoClient *ConceptDescriptionRepositoryClient) GetConceptDescriptionJsonable(cdID string) (map[string]any, error) {
	targetURL, err := getEncodedTargetUrl(repoClient.baseURL, ConceptDescriptionRepositoryPath, cdID)
	if err != nil {
		return nil, fmt.Errorf("failed to create target url: %w", err)
	}

	body, err := DoGetRequest(repoClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get concept description: %w", err)
	}
	defer body.Close()

	return bodyToJsonable(body)
}

// GetConceptDescription gets a parsed concept description
func (repoClient *ConceptDescriptionRepositoryClient) GetConceptDescription(cdID string) (types.IConceptDescription, error) {
	cdJsonable, err := repoClient.GetConceptDescriptionJsonable(cdID)
	if err != nil {
		return nil, err
	}

	conceptDescription, err := jsonization.ConceptDescriptionFromJsonable(cdJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse concept description json: %w", err)
	}

	return conceptDescription, nil
}

// UploadConceptDescription uploads a new concept description
func (repoClient *ConceptDescriptionRepositoryClient) UploadConceptDescription(conceptDescription types.IConceptDescription) error {
	if conceptDescription == nil {
		return fmt.Errorf("conceptDescription cannot be nil")
	}

	cdBytes, err := aasEntityToBytes(conceptDescription)
	if err != nil {
		return fmt.Errorf("failed to convert conceptDescription to bytes: %w", err)
	}

	targetUrl := repoClient.baseURL.JoinPath(ConceptDescriptionRepositoryPath)

	body, err := DoPostRequest(repoClient.httpClient, targetUrl.String(), cdBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

func (repoClient *ConceptDescriptionRepositoryClient) UpdateConceptDescription(cdID string, conceptDescription types.IConceptDescription) error {
	if conceptDescription == nil {
		return fmt.Errorf("concept description cannot be nil")
	}

	shellBytes, err := aasEntityToBytes(conceptDescription)
	if err != nil {
		return fmt.Errorf("failed to convert shell to bytes: %w", err)
	}

	targetURL, err := getEncodedTargetUrl(repoClient.baseURL, ConceptDescriptionRepositoryPath, cdID)
	if err != nil {
		return fmt.Errorf("failed to create target url: %w", err)
	}

	body, err := DoPutRequest(repoClient.httpClient, targetURL.String(), shellBytes)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

func (repoClient *ConceptDescriptionRepositoryClient) DeleteConceptDescription(cdID string) error {
	targetUrl, err := getEncodedTargetUrl(repoClient.baseURL, ConceptDescriptionRepositoryPath, cdID)
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
