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

type SubmodelServiceClient struct {
	// Http Client to use for requests
	httpClient *http.Client

	// baseURL to append path to
	baseURL *url.URL

	// if true, no /submodel is needed in path
	isDerived bool
}

// NewSubmodelServiceClient creates a standalone submodel service client, prefixing paths with /submodel
func NewSubmodelServiceClient(baseURL string, opts ...ClientOption) (*SubmodelServiceClient, error) {
	submodelServiceBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, "/submodel")
	if err != nil {
		return nil, fmt.Errorf("invalid url for submodel service: %w", err)
	}

	httpClient, err := createHttpClientFromOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client := &SubmodelServiceClient{
		httpClient: httpClient,
		baseURL:    submodelServiceBaseURL,
		isDerived:  false,
	}

	return client, nil
}

// NewDerivedSubmodelServiceClient creates a submodel service client for usage within a RepositoryClient
func NewDerivedSubmodelServiceClient(httpClient *http.Client, baseURL *url.URL) *SubmodelServiceClient {
	return &SubmodelServiceClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		isDerived:  true,
	}
}

// ---------------------------------------- Description -----------------------------
// GetSubmodelServiceDescription calls the /description endpoint and returns raw json
func (serviceClient *SubmodelServiceClient) GetSubmodelServiceDescription() ([]byte, error) {
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

// ---------------------------------------- Submodel -----------------------------
// GetJsonable returns the Submodel in jsonable format as map[string]any
func (serviceClient *SubmodelServiceClient) GetJsonable() (map[string]any, error) {
	targetURL := serviceClient.getDerivedTargetUrl()

	body, err := DoGetRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get Submodel: %w", err)
	}
	defer body.Close()

	return convert.BodyToJsonable(body)
}

// Get returns the parsed Submodel
func (serviceClient *SubmodelServiceClient) Get() (types.ISubmodel, error) {
	submodelJsonable, err := serviceClient.GetJsonable()
	if err != nil {
		return nil, err
	}

	submodel, err := jsonization.SubmodelFromJsonable(submodelJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse submodel json: %w", err)
	}

	return submodel, nil
}

// GetMetadataJsonable gets the metadata only representation of a submodel as jsonable
func (serviceClient *SubmodelServiceClient) GetMetadataJsonable() (map[string]any, error) {
	targetURL := serviceClient.getDerivedTargetUrl("/$metadata")

	body, err := DoGetRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get Submodel Metadata: %w", err)
	}
	defer body.Close()

	return convert.BodyToJsonable(body)
}

// GetMetadata gets the metadata only representation of a submodel
// Note: Although it has the same return type as Get(), it will be missing some values
func (serviceClient *SubmodelServiceClient) GetMetadata() (types.ISubmodel, error) {
	submodelJsonable, err := serviceClient.GetMetadataJsonable()
	if err != nil {
		return nil, err
	}

	submodel, err := jsonization.SubmodelFromJsonable(submodelJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Submodel Metadata json: %w", err)
	}

	return submodel, nil
}

// GetValueOnly returns the value-only representation of the submodel
// It returns raw json
func (serviceClient *SubmodelServiceClient) GetValueOnly() ([]byte, error) {
	targetURL := serviceClient.getDerivedTargetUrl("/$value")

	body, err := DoGetRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get Submodel ValueOnly: %w", err)
	}
	defer body.Close()

	result, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response for Submodel ValueOnly: %w", err)
	}

	return result, nil
}

// UpdateValueOnly simply sends input as it is to PATCH /$value
func (serviceClient *SubmodelServiceClient) UpdateValueOnly(input []byte) error {
	if input == nil {
		return fmt.Errorf("input for valueOnly cannot be nil")
	}

	targetURL := serviceClient.getDerivedTargetUrl("/$value")

	body, err := DoPatchRequest(serviceClient.httpClient, targetURL.String(), input)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// Update updates the submodel with PUT
func (serviceClient *SubmodelServiceClient) Update(submodel types.ISubmodel) error {
	if submodel == nil {
		return fmt.Errorf("submodel cannot be nil")
	}

	submodelBytes, err := convert.JsonableTypeToBytes(submodel)
	if err != nil {
		return fmt.Errorf("failed to convert submodel to bytes: %w", err)
	}

	targetURL := serviceClient.getDerivedTargetUrl()

	body, err := DoPutRequest(serviceClient.httpClient, targetURL.String(), submodelBytes)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// ---------------------------------------- Submodel Elements -----------------------------
// GetNextSubmodelElementsPage gets the next page of submodel elements for the submodel
// limit = 0 means no limit (however: BaSyx has internal limit anyway)
func (serviceClient *SubmodelServiceClient) GetNextSubmodelElementsPage(cursor string, limit int) (*PagedResult[types.ISubmodelElement], error) {
	targetURL := serviceClient.getDerivedTargetUrl("/submodel-elements")

	pagedResult, err := DoPagedGetRequest(serviceClient.httpClient, targetURL.String(), cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var typedResult PagedResult[types.ISubmodelElement]
	typedResult.Metadata = pagedResult.Metadata
	typedResult.Result = make([]types.ISubmodelElement, 0, len(pagedResult.Result))

	for _, rawIn := range pagedResult.Result {
		submodelElement, err := convert.JsonableTypeFromBytes(rawIn, jsonization.SubmodelElementFromJsonable)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SubmodelElement: %w", err)
		}
		typedResult.Result = append(typedResult.Result, submodelElement)
	}

	return &typedResult, nil
}

// GetSubmodelElementJsonable gets a submodel element at idShortPath in the jsonable format
func (serviceClient *SubmodelServiceClient) GetSubmodelElementJsonable(idShortPath string) (map[string]any, error) {
	targetURL := serviceClient.getDerivedTargetUrl("submodel-elements", idShortPath)

	body, err := DoGetRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return convert.BodyToJsonable(body)
}

// GetSubmodelElement gets a submodel element at idShortPath
func (serviceClient *SubmodelServiceClient) GetSubmodelElement(idShortPath string) (types.ISubmodelElement, error) {
	submodelElementJsonable, err := serviceClient.GetSubmodelElementJsonable(idShortPath)
	if err != nil {
		return nil, err
	}

	submodelElement, err := jsonization.SubmodelElementFromJsonable(submodelElementJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SubmodelElement JSON: %w", err)
	}

	return submodelElement, err
}

// GetSubmodelElementValueOnly gets the ValueOnly representation of a SubmodelElement in raw json
func (serviceClient *SubmodelServiceClient) GetSubmodelElementValueOnly(idShortPath string) ([]byte, error) {
	targetURL := serviceClient.getDerivedTargetUrl("submodel-elements", idShortPath, "$value")

	body, err := DoGetRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	result, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response bytes: %w", err)
	}

	return result, nil
}

// UploadSubmodelElement uploads a submodel element with POST /submodel-elements/:idShortPath (or /submodel-elements if len(idShortPath) == 0)
func (serviceClient *SubmodelServiceClient) UploadSubmodelElement(idShortPath string, submodelElement types.ISubmodelElement) error {
	if submodelElement == nil {
		return fmt.Errorf("SubmodelElement cannot be nil")
	}

	submodelElementBytes, err := convert.JsonableTypeToBytes(submodelElement)
	if err != nil {
		return fmt.Errorf("failed to convert SubmodelElement to bytes: %w", err)
	}

	var targetURL *url.URL

	if len(idShortPath) > 0 {
		targetURL = serviceClient.getDerivedTargetUrl("/submodel-elements", idShortPath)
	} else {
		targetURL = serviceClient.getDerivedTargetUrl("/submodel-elements")
	}

	body, err := DoPostRequest(serviceClient.httpClient, targetURL.String(), submodelElementBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateSubmodelElement updates a submodel element with PUT /submodel-elements/:idShortPath
func (serviceClient *SubmodelServiceClient) UpdateSubmodelElement(idShortPath string, submodelElement types.ISubmodelElement) error {
	if submodelElement == nil {
		return fmt.Errorf("SubmodelElement cannot be nil")
	}

	submodelElementBytes, err := convert.JsonableTypeToBytes(submodelElement)
	if err != nil {
		return fmt.Errorf("failed to convert SubmodelElement to bytes: %w", err)
	}

	targetURL := serviceClient.getDerivedTargetUrl("/submodel-elements", idShortPath)

	body, err := DoPutRequest(serviceClient.httpClient, targetURL.String(), submodelElementBytes)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// UpdateSubmodelElementValueOnly updates a submodel element in ValueOnly mode with PATCH /submodel-elements/:idShortPath/$value
func (serviceClient *SubmodelServiceClient) UpdateSubmodelElementValueOnly(idShortPath string, input []byte) error {
	if input == nil {
		return fmt.Errorf("input for valueOnly cannot be nil")
	}

	targetURL := serviceClient.getDerivedTargetUrl("submodel-elements", idShortPath, "$value")

	body, err := DoPatchRequest(serviceClient.httpClient, targetURL.String(), input)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// DeleteSubmodelElement deletes a submodel element with DELETE /submodel-elements/:idShortPath
func (serviceClient *SubmodelServiceClient) DeleteSubmodelElement(idShortPath string) error {
	if len(idShortPath) == 0 {
		return fmt.Errorf("idShortPath cannot be empty")
	}

	targetURL := serviceClient.getDerivedTargetUrl("submodel-elements", idShortPath)

	body, err := DoDeleteRequest(serviceClient.httpClient, targetURL.String())
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	return nil
}

// ---------------------------------------- Operations -----------------------------
// InvokeOperation invokes a synchronous operation on idShortPath using operationRequest as input
func (serviceClient *SubmodelServiceClient) InvokeOperation(idShortPath string, operationRequest *OperationRequest) (*OperationResult, error) {
	if operationRequest == nil {
		return nil, fmt.Errorf("operationRequest cannot be nil")
	}

	// create target url to /invoke endpoint
	targetURL := serviceClient.getDerivedTargetUrl("submodel-elements", idShortPath, "invoke")

	// convert OperationRequest into JSON
	// this is more work since OperationRequest is not part of core works and
	// therefore not within jsonization - but its elements are.
	jsonRequest := OperationRequestJSON{}

	inputArguments, err := convert.JsonableTypeListToRawJsonMessages(operationRequest.InputArguments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inputArguments: %w", err)
	}

	inoutputArguments, err := convert.JsonableTypeListToRawJsonMessages(operationRequest.InoutputArguments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inoutputArguments: %w", err)
	}

	jsonRequest.InputArguments = inputArguments
	jsonRequest.InoutputArguments = inoutputArguments

	jsonBytes, err := json.Marshal(jsonRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to convert OperationRequest to json: %w", err)
	}

	body, err := DoPostRequest(serviceClient.httpClient, targetURL.String(), jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer body.Close()

	var operationResultJSON OperationResultJSON
	if err := json.NewDecoder(body).Decode(&operationResultJSON); err != nil {
		return nil, fmt.Errorf("failed to read response as OperationResult: %w", err)
	}

	operationResult := OperationResult{
		Success: operationResultJSON.Success,
	}

	operationResult.OutputArguments, err = convert.JsonableTypeListFromRawJsonMessages(operationResultJSON.OutputArguments, jsonization.OperationVariableFromJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result output arguments: %w", err)
	}

	operationResult.InoutputArguments, err = convert.JsonableTypeListFromRawJsonMessages(operationResultJSON.InoutputArguments, jsonization.OperationVariableFromJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result inoutput arguments: %w", err)
	}

	return &operationResult, nil

}

// ----------------------------------- Attachments ---------------------------------
// TODO: Add!

// ---------------------------------------- Util -----------------------------
// getDerivedTargetUrl returns the url to use, adding /submodel if isDerived is true
func (serviceClient *SubmodelServiceClient) getDerivedTargetUrl(pathItems ...string) *url.URL {
	toAppend := make([]string, 0, len(pathItems)+1)

	if !serviceClient.isDerived {
		toAppend = append(toAppend, "submodel")
	}

	toAppend = append(toAppend, pathItems...)

	for idx, pathItem := range toAppend {
		toAppend[idx] = url.PathEscape(pathItem)
	}

	targetURL := serviceClient.baseURL.JoinPath(toAppend...)
	return targetURL
}
