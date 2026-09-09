package traverse

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aas-core-works/aas-core3.1-golang/stringification"
	"github.com/aas-core-works/aas-core3.1-golang/types"
	id_short_path "github.com/smartfactory-kl/shellgonaut/path"
)

// FindSubmodelByIdShort looks for a submodel within a list of ISubmodel by checking idShort.
func FindSubmodelByIdShort(submodelList []types.ISubmodel, idShort string) (types.ISubmodel, error) {
	if len(submodelList) == 0 {
		return nil, fmt.Errorf("list of submodels is empty")
	}

	submodelIdx := slices.IndexFunc(submodelList, func(el types.ISubmodel) bool {
		return el != nil && el.IDShort() != nil && *el.IDShort() == idShort
	})

	if submodelIdx == -1 {
		return nil, fmt.Errorf("no submodel %q in list", idShort)
	}

	submodel := submodelList[submodelIdx]
	if submodel == nil {
		return nil, fmt.Errorf("submodel found for %q but it was nil", idShort)
	}

	return submodel, nil
}

func GetFirstSubmodelElementOfTypeInListable[T types.ISubmodelElement](root types.IClass, submodelElementModelType types.ModelType) (T, error) {
	var zero T

	submodelElementList, listableDescription, err := GetSubmodelElementListFromListable(root)
	if err != nil {
		return zero, fmt.Errorf("failed to get list of submodel elements for listable: %w", err)
	}

	elementIdx := slices.IndexFunc(submodelElementList, func(element types.ISubmodelElement) bool {
		return element != nil && element.ModelType() == submodelElementModelType
	})

	if elementIdx == -1 {
		return zero, fmt.Errorf("no modelType %q found in listable %q", stringification.MustModelTypeToString(submodelElementModelType), listableDescription)
	}

	resultElement, ok := submodelElementList[elementIdx].(T)
	if !ok {
		return zero, fmt.Errorf("element %q found in %q but failed type assertion", stringification.MustModelTypeToString(submodelElementModelType), listableDescription)
	}

	return resultElement, nil
}

func GetSubmodelElementsOfTypeInListable[T types.ISubmodelElement](root types.IClass, submodelElementModelType types.ModelType) ([]T, error) {
	var list []T

	submodelElementList, _, err := GetSubmodelElementListFromListable(root)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of submodel elements for listable: %w", err)
	}

	for idx, submodelElement := range submodelElementList {
		if submodelElement != nil && submodelElement.ModelType() == submodelElementModelType {
			resultElement, ok := submodelElementList[idx].(T)
			if ok {
				list = append(list, resultElement)
			}
		}
	}

	return list, nil
}

// GetSubmodelElementOfTypeByIdShortPath gets you an element of type T on idShortPath. If anything does not match, an error is returned
func GetSubmodelElementOfTypeByIdShortPath[T types.ISubmodelElement](root types.IClass, idShortPath string, submodelElementModelType types.ModelType) (T, error) {
	var zero T

	element, err := GetSubmodelElementByIdShortPath(root, idShortPath)
	if err != nil {
		return zero, fmt.Errorf("failed to find element: %w", err)
	}

	if element == nil {
		return zero, fmt.Errorf("element cannot be nil")
	}

	if element.ModelType() != submodelElementModelType {
		return zero, fmt.Errorf("failed to find element: expected model type %s but got %s", stringification.MustModelTypeToString(submodelElementModelType), stringification.MustModelTypeToString(element.ModelType()))
	}

	assertedElement, ok := element.(T)
	if !ok {
		return zero, fmt.Errorf("failed to assert element type: %w", err)
	}

	return assertedElement, nil
}

// GetSubmodelElementByIdShortPath looks for a generic SubmodelElement by stepping through idShortPath
func GetSubmodelElementByIdShortPath(root types.IClass, idShortPath string) (types.ISubmodelElement, error) {
	if root == nil {
		return nil, fmt.Errorf("root cannot be nil")
	}

	// gather steps from idShortPath
	stepList, err := id_short_path.GatherStepsFromIdShortPath(idShortPath)
	if err != nil {
		return nil, fmt.Errorf("invalid idShortPath: %w", err)
	}
	// walk all steps
	currentRoot := root
	for _, stepItem := range stepList {
		// nil check
		if currentRoot == nil {
			return nil, fmt.Errorf("cannot go step: element cannot be nil")
		}

		// check if step basics are still valid
		switch stepItem.StepType {
		case id_short_path.StepTypeIDShort:
			// for this, the next item has to be a listable (including SubmodelElementList)
			var err error
			currentRoot, err = GetElementInListable(currentRoot, stepItem.NextIdShort)
			if err != nil {
				return nil, fmt.Errorf("cannot go step: StepTypeIDShort: %w", err)
			}
		case id_short_path.StepTypeIndex:
			// for this, the next item has to be a submodel element list
			if currentRoot.ModelType() != types.ModelTypeSubmodelElementList {
				return nil, fmt.Errorf("cannot go step: StepTypeIndex expected SubmodelElementList, got %s instead for path %s", stringification.MustModelTypeToString(currentRoot.ModelType()), idShortPath)
			}

			listElement, ok := currentRoot.(types.ISubmodelElementList)
			if !ok {
				return nil, fmt.Errorf("cannot go step: StepTypeIndex got modelType SubmodelElementList, but the type assertion failed for path %s.", idShortPath)
			}

			if stepItem.NextIndex < 0 || stepItem.NextIndex >= len(listElement.Value()) {
				return nil, fmt.Errorf("cannot go step: StepTypeIndex: idx %d out of range %d in %s.", stepItem.NextIndex, len(listElement.Value()), idShortPath)
			}

			if listElement.Value() == nil {
				return nil, fmt.Errorf("cannot go step: StepTypeIndex: list has nil value")
			}

			currentRoot = listElement.Value()[stepItem.NextIndex]
		default:
			return nil, fmt.Errorf("invalid step type: %d", stepItem.StepType)
		}
	}

	if currentRoot == nil {
		return nil, fmt.Errorf("failed to find element: final element was nil")
	}

	lastElement, ok := currentRoot.(types.ISubmodelElement)
	if !ok {
		return nil, fmt.Errorf("failed type assertion, final element is not a submodel element")
	}

	// return the item
	return lastElement, nil
}

// GetSubmodelElementListFromListable returns a []ISubmodelElement of a listable, a brief description of it for logging and an error if the modelType is not listable
func GetSubmodelElementListFromListable(listable types.IClass) ([]types.ISubmodelElement, string, error) {
	var submodelElementList []types.ISubmodelElement
	listableDescription := ""

	if listable == nil {
		return nil, "", fmt.Errorf("listabel cannot be nil")
	}

	switch listable.ModelType() {
	case types.ModelTypeSubmodel:
		submodelElementList = listable.(types.ISubmodel).SubmodelElements()
	case types.ModelTypeSubmodelElementCollection:
		submodelElementList = listable.(types.ISubmodelElementCollection).Value()
	case types.ModelTypeSubmodelElementList:
		submodelElementList = listable.(types.ISubmodelElementList).Value()
	case types.ModelTypeEntity:
		submodelElementList = listable.(types.IEntity).Statements()
	case types.ModelTypeAnnotatedRelationshipElement:
		list := listable.(types.IAnnotatedRelationshipElement).Annotations()
		submodelElementList = make([]types.ISubmodelElement, 0, len(list))

		for _, item := range list {
			submodelElementList = append(submodelElementList, item)
		}
	default:
		return nil, "", fmt.Errorf("provided modelType %q is not a listable", stringification.MustModelTypeToString(listable.ModelType()))
	}

	idShort := listable.(types.IReferable).IDShort()
	if idShort == nil {
		idShort = new("")
	}

	listableDescription = fmt.Sprintf("%q [%s]", *idShort, stringification.MustModelTypeToString(listable.ModelType()))

	return submodelElementList, listableDescription, nil
}

// GetElementInListable simply looks for a idShort within a listable container
func GetElementInListable(listable types.IClass, idShort string) (types.ISubmodelElement, error) {
	submodelElementList, listableDescription, err := GetSubmodelElementListFromListable(listable)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of submodel elements for listable: %w", err)
	}

	elementIdx := slices.IndexFunc(submodelElementList, func(el types.ISubmodelElement) bool {
		return el != nil && el.IDShort() != nil && *el.IDShort() == idShort
	})

	if elementIdx == -1 {
		return nil, fmt.Errorf("idShort %q not found in listable %q", idShort, listableDescription)
	}

	return submodelElementList[elementIdx], nil
}

func GetFirstSubmodelIDForSemanticID(submodels []types.IReference, semanticID string) (string, error) {
	if len(strings.TrimSpace(semanticID)) == 0 {
		return "", fmt.Errorf("semanticID cannot be empty")
	}

	if len(submodels) == 0 {
		return "", fmt.Errorf("no submodels to check")
	}

	var firstMatchReference types.IReference
	for _, submodelReference := range submodels {
		if submodelReference == nil {
			continue
		}

		refID := submodelReference.ReferredSemanticID()
		if refID == nil {
			continue
		}

		if len(refID.Keys()) == 0 {
			continue
		}

		for _, keyEntry := range refID.Keys() {
			if keyEntry == nil {
				continue
			}

			if strings.EqualFold(keyEntry.Value(), semanticID) {
				// found matching semantics, now use reference keys to get submodel ID
				firstMatchReference = submodelReference
			}
		}
	}

	if firstMatchReference != nil {
		smKeys := firstMatchReference.Keys()
		if len(smKeys) == 0 {
			return "", fmt.Errorf("found semantics but keys were empty for semanticID %s", semanticID)
		}

		for _, smKeyEntry := range smKeys {
			if smKeyEntry == nil {
				continue
			}

			if smKeyEntry.Type() == types.KeyTypesSubmodel {
				return smKeyEntry.Value(), nil
			}
		}
	}

	return "", fmt.Errorf("no submodel found for semanticID %s", semanticID)
}
