package traverse

import (
	"fmt"
	"slices"

	"github.com/aas-core-works/aas-core3.1-golang/types"
)

// SetPropertyValueWithinListable is a shorthand to set the value of a Property within a listable. Usually used to set values that stem from a template json
func SetPropertyValueWithinListable(listable types.IClass, idShort string, value string) error {
	var submodelElementList []types.ISubmodelElement

	submodelElementList, listableDescription, err := GetSubmodelElementListFromListable(listable)
	if err != nil {
		return fmt.Errorf("failed to get list of submodel elements for listable: %w", err)
	}

	propertyIdx := slices.IndexFunc(submodelElementList, func(el types.ISubmodelElement) bool {
		return el != nil && types.IsProperty(el) && el.IDShort() != nil && *el.IDShort() == idShort
	})

	if propertyIdx == -1 {
		return fmt.Errorf("idShort %q not found in listable %q", idShort, listableDescription)
	}

	property, ok := submodelElementList[propertyIdx].(types.IProperty)
	if !ok {
		return fmt.Errorf("typing mismatch: property is not of type IProperty for idShort %q in listable of type %q", idShort, listableDescription)
	}

	property.SetValue(&value)

	return nil
}

// SetMultiLanguagePropertyValueWithinListable is a shorthand to set language/value combo. Usually used to set template values that stem from a json template
func SetMultiLanguagePropertyValueWithinListable(listable types.IClass, idShort string, language string, value string) error {
	var submodelElementList []types.ISubmodelElement
	var listableDescription = ""

	submodelElementList, listableDescription, err := GetSubmodelElementListFromListable(listable)
	if err != nil {
		return fmt.Errorf("failed to get list of submodel elements for listable: %w", err)
	}

	propertyIdx := slices.IndexFunc(submodelElementList, func(el types.ISubmodelElement) bool {
		return el != nil && types.IsProperty(el) && el.IDShort() != nil && *el.IDShort() == idShort
	})

	if propertyIdx == -1 {
		return fmt.Errorf("idShort %q not found in listable %q", idShort, listableDescription)
	}

	property, ok := submodelElementList[propertyIdx].(types.IProperty)
	if !ok {
		return fmt.Errorf("typing mismatch: property is not of type IProperty for idShort %q in listable of type %q", idShort, listableDescription)
	}

	property.SetValue(&value)

	return nil
}
