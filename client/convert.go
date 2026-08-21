package client

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

// -------------------------------- HELPER --------------------------------------------------
// bodyToJsonable takes a io.ReadCloser and tries to read a "jsonable" map[string]any from it.
// The caller has to close the body io.ReadCloser.
func bodyToJsonable(body io.ReadCloser) (map[string]any, error) {
	var jsonBody map[string]any
	if err := json.NewDecoder(body).Decode(&jsonBody); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return jsonBody, nil
}

// bodyToJsonableList takes a io.ReadCloser and tries to read a "jsonable" []map[string]any from it for listing responses
// The caller has to close the body io.ReadCloser.
func bodyToJsonableList(body io.ReadCloser) ([]map[string]any, error) {
	var jsonListBody []map[string]any
	if err := json.NewDecoder(body).Decode(&jsonListBody); err != nil {
		return nil, fmt.Errorf("failed to parse json list: %w", err)
	}

	return jsonListBody, nil
}

// aasEntityToBytes takes any AAS entity and converts it to JSON represented via []byte
func aasEntityToBytes(entity types.IClass) ([]byte, error) {
	if entity == nil {
		return nil, fmt.Errorf("failed to convert entity to jsonable: entity cannot be nil")
	}

	entityJsonable, err := jsonization.ToJsonable(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to convert entity to jsonable: %w", err)
	}

	entityJson, err := json.Marshal(entityJsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal jsonable: %w", err)
	}

	return entityJson, nil
}

// aasEntityListToBytes converts all entities to json raw and then marshals them to a json list represented as []byte
func aasEntityListToBytes[T types.IClass](entityList []T) ([]byte, error) {
	result, err := aasEntityListToJsonRaw(entityList)
	if err != nil {
		return nil, err
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal list: %w", err)
	}

	return resultJson, nil
}

// aasEntityListToJsonRaw converts a list of entities into a list of raw json messages representing them
func aasEntityListToJsonRaw[T types.IClass](entityList []T) ([]json.RawMessage, error) {
	if len(entityList) == 0 {
		return []json.RawMessage{}, nil
	}

	result := make([]json.RawMessage, 0, len(entityList))

	for _, entityItem := range entityList {
		item, err := aasEntityToBytes(entityItem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert list item: %w", err)
		}
		result = append(result, json.RawMessage(item))
	}

	return result, nil
}
