package convert

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

// CreateCopy tries to create a copy of src by converting it to bytes and back using jsonizationFunction
func CloneJsonableType[T types.IClass](src T, jsonizationFunction func(any) (T, error)) (T, error) {
	var zero T

	srcBytes, err := JsonableTypeToBytes(src)
	if err != nil {
		return zero, err
	}

	dst, err := JsonableTypeFromBytes(srcBytes, jsonizationFunction)
	if err != nil {
		return zero, err
	}

	return dst, nil
}

// -------------------------------- HELPER --------------------------------------------------
// bodyToJsonable takes a io.ReadCloser and tries to read a "jsonable" map[string]any from it.
// The caller has to close the body io.ReadCloser.
func BodyToJsonable(body io.ReadCloser) (map[string]any, error) {
	var jsonBody map[string]any
	if err := json.NewDecoder(body).Decode(&jsonBody); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return jsonBody, nil
}

// bodyToJsonableList takes a io.ReadCloser and tries to read a "jsonable" []map[string]any from it for listing responses
// The caller has to close the body io.ReadCloser.
func BodyToJsonableList(body io.ReadCloser) ([]map[string]any, error) {
	var jsonListBody []map[string]any
	if err := json.NewDecoder(body).Decode(&jsonListBody); err != nil {
		return nil, fmt.Errorf("failed to parse json list: %w", err)
	}

	return jsonListBody, nil
}

// JsonableToBytes takes any AAS entity and converts it to JSON represented via []byte
func JsonableTypeToBytes(entity types.IClass) ([]byte, error) {
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

// FromBytes tries to create a model element from simple bytes by using the provided jsonization. ...FromJsonable method
func JsonableTypeFromBytes[T types.IClass](input []byte, jsonizationFunction func(any) (T, error)) (T, error) {
	var zero T

	var jsonable map[string]any
	if err := json.Unmarshal(input, &jsonable); err != nil {
		return zero, fmt.Errorf("failed to parse json from bytes: %w", err)
	}

	element, err := jsonizationFunction(jsonable)
	if err != nil {
		return zero, fmt.Errorf("failed to parse jsonable into element: %w", err)
	}

	return element, nil
}

// JsonableTypeListToBytes converts all entities to json raw and then marshals them to a json list represented as []byte
func JsonableTypeListToBytes[T types.IClass](entityList []T) ([]byte, error) {
	result, err := JsonableTypeListToRawJsonMessages(entityList)
	if err != nil {
		return nil, err
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal list: %w", err)
	}

	return resultJson, nil
}

func JsonableTypeListFromBytes[T types.IClass](input []byte, jsonizationFunction func(any) (T, error)) ([]T, error) {
	var inputList []json.RawMessage

	if err := json.Unmarshal(input, &inputList); err != nil {
		return []T{}, fmt.Errorf("failed to unmarshal input to list: %w", err)
	}

	return JsonableTypeListFromRawJsonMessages(inputList, jsonizationFunction)
}

// JsonableTypeListToRawJsonMessage converts a list of entities into a list of raw json messages representing them
func JsonableTypeListToRawJsonMessages[T types.IClass](entityList []T) ([]json.RawMessage, error) {
	if len(entityList) == 0 {
		return []json.RawMessage{}, nil
	}

	result := make([]json.RawMessage, 0, len(entityList))

	for _, entityItem := range entityList {
		item, err := JsonableTypeToBytes(entityItem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert list item: %w", err)
		}
		result = append(result, json.RawMessage(item))
	}

	return result, nil
}

func JsonableTypeListFromRawJsonMessages[T types.IClass](inputList []json.RawMessage, jsonizationFunction func(any) (T, error)) ([]T, error) {
	if len(inputList) == 0 {
		return []T{}, nil
	}

	var resultList []T = []T{}

	for idx, inputItem := range inputList {
		resultItem, err := JsonableTypeFromBytes(inputItem, jsonizationFunction)
		if err != nil {
			return []T{}, fmt.Errorf("failure on #%d", idx)
		}

		resultList = append(resultList, resultItem)
	}

	return resultList, nil
}
