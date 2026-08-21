package create

import (
	"encoding/json"
	"fmt"

	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

// CreateCopy tries to create a copy of src by converting it to bytes and back using jsonizationFunction
func CreateCopy[T types.IClass](src T, jsonizationFunction func(any) (T, error)) (T, error) {
	var zero T

	srcBytes, err := ToBytes(src)
	if err != nil {
		return zero, err
	}

	dst, err := FromBytes(srcBytes, jsonizationFunction)
	if err != nil {
		return zero, err
	}

	return dst, nil
}

// FromBytes tries to create a model element from simple bytes by using the provided jsonization. ...FromJsonable method
func FromBytes[T types.IClass](input []byte, jsonizationFunction func(any) (T, error)) (T, error) {
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

func ListFromBytes[T types.IClass](input []byte, jsonizationFunction func(any) (T, error)) ([]T, error) {
	var inputList []json.RawMessage

	if err := json.Unmarshal(input, &inputList); err != nil {
		return []T{}, fmt.Errorf("failed to unmarshal input to list: %w", err)
	}

	return ListFromJsonRawMessages(inputList, jsonizationFunction)
}

func ListFromJsonRawMessages[T types.IClass](inputList []json.RawMessage, jsonizationFunction func(any) (T, error)) ([]T, error) {
	if len(inputList) == 0 {
		return []T{}, nil
	}

	var resultList []T = []T{}

	for idx, inputItem := range inputList {
		resultItem, err := FromBytes(inputItem, jsonizationFunction)
		if err != nil {
			return []T{}, fmt.Errorf("failure on #%d", idx)
		}

		resultList = append(resultList, resultItem)
	}

	return resultList, nil
}

// ToBytes is a helper to create json []byte from any types.IClass
func ToBytes(input types.IClass) ([]byte, error) {
	jsonable, err := jsonization.ToJsonable(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert element to jsonable: %w", err)
	}

	byteContent, err := json.Marshal(jsonable)
	if err != nil {
		return nil, fmt.Errorf("failed to convert jsonable to bytes: %w", err)
	}

	return byteContent, nil
}

// ListToBytes is ToBytes, but for lists!
func ListToBytes[T types.IClass](input []T) ([]byte, error) {
	result := []json.RawMessage{}

	for idx, entry := range input {
		item, err := ToBytes(entry)
		if err != nil {
			return nil, fmt.Errorf("failed to convert list element %d to bytes: %w", idx, err)
		}
		result = append(result, item)
	}

	content, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert rawMessage list to bytes: %w", err)
	}

	return content, nil
}
