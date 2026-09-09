package client

// TODO: Find out how to handle Descriptors and their jsonization

// const SubmodelRegistryPath = "/submodel-descriptors"

// type SubmodelRegistryClient struct {
// 	httpClient *http.Client
// 	baseURL    *url.URL
// }

// // NewSubmodelRegistryClient creates a new client for a submodel registry
// func NewSubmodelRegistryClient(baseURL string, opts ...ClientOption) (*SubmodelRegistryClient, error) {
// 	submodelRegistryBaseURL, err := EnsureUrlWithoutSuffixOrSlash(baseURL, SubmodelRegistryPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid url for submodel registry: %w", err)
// 	}

// 	httpClient, err := createHttpClientFromOptions(opts...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create http client: %w", err)
// 	}

// 	client := &SubmodelRegistryClient{
// 		httpClient: httpClient,
// 		baseURL:    submodelRegistryBaseURL,
// 	}

// 	return client, nil
// }

// // ---------------------------------------- Description -----------------------------
// // GetSubmodelRegistryDescription calls the /description endpoint. Might be used to check for availability. Returns raw JSON
// func (regClient *SubmodelRegistryClient) GetSubmodelRegistryDescription() ([]byte, error) {
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

// // -------------------------- Submodel Descriptors --------------------------------
// // GetNextSubmodelDescriptorPage requests the next page of submodel descriptors
// func (regClient *SubmodelRegistryClient) GetNextSubmodelDescriptorPage(cursor string, limit int) (*PagedResult[descriptor.ISubmodelDescriptor], error) {
// 	targetURL := regClient.baseURL.JoinPath(SubmodelRegistryPath)

// 	pagedResult, err := DoPagedGetRequest(regClient.httpClient, targetURL.String(), cursor, limit)
// 	if err != nil {
// 		return nil, fmt.Errorf("request failed: %w", err)
// 	}

// 	var typedResult PagedResult[descriptor.ISubmodelDescriptor]
// 	typedResult.Metadata = pagedResult.Metadata
// 	typedResult.Result = make([]descriptor.ISubmodelDescriptor, 0, len(pagedResult.Result))

// 	for _, rawIn := range pagedResult.Result {
// 		submodelDescriptor, err := create.FromBytes(rawIn, descriptor.SubmodelDescriptorFromJsonable)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to parse submodel descriptor: %w", err)
// 		}
// 		typedResult.Result = append(typedResult.Result, submodelDescriptor)
// 	}

// 	return &typedResult, nil
// }
