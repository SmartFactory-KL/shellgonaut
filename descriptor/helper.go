package descriptor

// UPDATE: For the first repo only release not needed, therefore not active

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
// 	"github.com/aas-core-works/aas-core3.1-golang/reporting"
// 	"github.com/aas-core-works/aas-core3.1-golang/types"
// )

// func stringFromJsonable(jsonable any) (string, error) {
// 	return simpletonFromJsonable(jsonable, "string", "")
// }

// func boolFromJsonable(jsonable any) (bool, error) {
// 	return simpletonFromJsonable(jsonable, "bool", false)
// }

// func float64FromJsonable(jsonable any) (float64, error) {
// 	return simpletonFromJsonable(jsonable, "number", 0.0)
// }

// func simpletonFromJsonable[T any](jsonable any, expected string, zero T) (T, error) {
// 	if jsonable == nil {
// 		return zero, newDeserializationError(
// 			fmt.Sprintf("expected %s, but got nil", expected),
// 		)
// 	}

// 	result, ok := jsonable.(T)
// 	if ok {
// 		return result, nil
// 	}

// 	return zero, newDeserializationError(
// 		fmt.Sprintf("expected %s, but got %T", expected, jsonable),
// 	)
// }

// func simpletonListFromJsonable[T any](jsonable any, expected string, zero T) ([]T, error) {
// 	if jsonable == nil {
// 		return []T{}, newDeserializationError(
// 			fmt.Sprintf("expected list of %s, but got nil", expected),
// 		)
// 	}

// 	arr, ok := jsonable.([]any)
// 	if !ok {
// 		return []T{}, newDeserializationError(
// 			fmt.Sprintf("expeceted list of %s, but got %T instead", expected, jsonable),
// 		)
// 	}

// 	if len(arr) == 0 {
// 		return []T{}, nil
// 	}

// 	var result []T
// 	for idx, arrItem := range arr {
// 		val, err := simpletonFromJsonable(arrItem, expected, zero)
// 		if err != nil {
// 			return []T{}, newDeserializationError(
// 				fmt.Sprintf("item #%s in array expected as %s but parsing failed: %w", idx, expected, err),
// 			)
// 		}
// 		result = append(result, val)
// 	}

// 	return result, nil
// }

// func jsonizationFromJsonable[T types.IClass](jsonable any, jsonizationFunction func(any) (T, error)) (T, error) {
// 	var zero T

// 	if jsonable == nil {
// 		return zero, newDeserializationError("expected json object, but got nil")
// 	}

// 	jsonMap, ok := jsonable.(map[string]any)
// 	if !ok {
// 		return zero, newDeserializationError(
// 			fmt.Sprintf("expected json object, but got %T instead", jsonable),
// 		)
// 	}

// 	item, err := jsonizationFunction(jsonMap)
// 	if err != nil {
// 		return zero, newDeserializationError(
// 			fmt.Sprintf("failed to parse: %w", err),
// 		)
// 	}

// 	return item, nil
// }

// func jsonizationFromJsonableList[T types.IClass](jsonable any, jsonizationFunction func(any) (T, error)) ([]T, error) {
// 	if jsonable == nil {
// 		return []T{}, newDeserializationError("expected list but got nil")
// 	}

// 	arr, ok := jsonable.([]map[string]any)
// 	if !ok {
// 		return []T{}, newDeserializationError(
// 			fmt.Sprintf("expeceted list, but got %T instead", jsonable),
// 		)
// 	}

// 	if len(arr) == 0 {
// 		return []T{}, nil
// 	}

// 	var result []T
// 	for idx, arrItem := range arr {
// 		val, err := jsonizationFromJsonable(arrItem, jsonizationFunction)
// 		if err != nil {
// 			return []T{}, newDeserializationError(
// 				fmt.Sprintf("item #%s in array parsing failed: %w", idx, err),
// 			)
// 		}
// 		result = append(result, val)
// 	}

// 	return result, nil
// }

// func prependName(err error, name string) error {
// 	var de *jsonization.DeserializationError
// 	if errors.As(err, &de) {
// 		de.Path.PrependName(&reporting.NameSegment{Name: name})
// 	}
// 	return err
// }

// func isValidJsonable(input any) (map[string]any, error) {
// 	if input == nil {
// 		return nil, newDeserializationError("expected a JSON object, but got null")
// 	}

// 	m, ok := input.(map[string]any)
// 	if !ok {
// 		err := newDeserializationError(
// 			fmt.Sprintf(
// 				"Expected a JSON object, but got %T",
// 				input,
// 			),
// 		)
// 		return nil, err
// 	}

// 	return m, nil
// }
