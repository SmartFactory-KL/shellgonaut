package client

import (
	"fmt"

	"github.com/aas-core-works/aas-core3.1-golang/stringification"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

// ExtractModelTypeFromJsonable tries to peek into the map to get the modelType so it can be used
// to switch the correct jsonization method
func ExtractModelTypeFromJsonable(input map[string]any) (types.ModelType, error) {
	stringInput, ok := input["modelType"]
	if !ok {
		return 0, fmt.Errorf("Jsonable does not contain modelType or its not a string")
	}

	modelTypeString, ok := stringInput.(string)
	if !ok {
		return 0, fmt.Errorf("Jsonable does not contain modelType or its not a string")
	}

	modelType, ok := stringification.ModelTypeFromString(modelTypeString)
	if !ok {
		return 0, fmt.Errorf("Jsonable contains invalid modelType %s", modelTypeString)
	}

	return modelType, nil
}
