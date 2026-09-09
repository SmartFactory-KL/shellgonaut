package traverse

import (
	"fmt"
	"strconv"

	"github.com/aas-core-works/aas-core3.1-golang/types"
)

func ReadPropertyValueAsString(property types.IProperty) (string, error) {
	if property == nil {
		return "", fmt.Errorf("Property was nil")
	}

	propValue := property.Value()
	if propValue == nil {
		return "", fmt.Errorf("failed to read property value of %q: value was nil", IDShortOrEmpty(property.IDShort()))
	}

	return *propValue, nil
}

func ReadPropertyValueAsInt64(property types.IProperty) (int64, error) {
	if property == nil {
		return 0, fmt.Errorf("Property was nil")
	}

	propValue := property.Value()
	if propValue == nil {
		return 0, fmt.Errorf("failed to read property value of %q: value was nil", IDShortOrEmpty(property.IDShort()))
	}

	return strconv.ParseInt(*propValue, 10, 64)
}

func ReadPropertyValueAsFloat64(property types.IProperty) (float64, error) {
	if property == nil {
		return 0, fmt.Errorf("Property was nil")
	}

	propValue := property.Value()
	if propValue == nil {
		return 0, fmt.Errorf("failed to read property value of %q: value was nil", IDShortOrEmpty(property.IDShort()))
	}

	return strconv.ParseFloat(*propValue, 64)
}

func ReadFirstSubmodelIDFromReference(reference types.IReference) (string, error) {
	if reference == nil {
		return "", fmt.Errorf("reference cannot be nil")
	}

	for _, key := range reference.Keys() {
		if key.Type() == types.KeyTypesSubmodel {
			if len(key.Value()) > 0 {
				return key.Value(), nil
			}
		}
	}

	return "", fmt.Errorf("reference contained no non-empty submodel key")
}

func IDShortOrEmpty(idShort *string) string {
	if idShort == nil {
		return ""
	}

	return *idShort
}
