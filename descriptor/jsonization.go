package descriptor

// UPDATE: For the first repo only release not needed, therefore not active

// import (
// 	"fmt"
// )

// /**
// IMPORTANT NOTE:
// These are only aimed at reading shell descriptor and submodel descriptor endpoints.
// They will not work with the "ToJsonable" method from aas-core-works!
// Maybe that will get added here in the future.
// */

// func AssetAdministrationShellDescriptorFromJsonable(jsonable any) (IAssetAdministrationShellDescriptor, error) {
// 	m, err := isValidJsonable(jsonable)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := shellDescriptorFromMap(m)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func shellDescriptorFromMap(m map[string]any) (result IAssetAdministrationShellDescriptor, err error) {
// 	// var theDescription types.ILangStringTextType
// 	// var theDisplayName types.ILangStringNameType
// 	// var theExtensions []types.IExtension
// 	// var theAdministration types.IAdministrativeInformation
// 	// var theAssetKind *types.AssetKind
// 	// var theAssetType *string
// 	// var theEndpoints []IEndpoint
// 	// var theGlobalAssetID *string
// 	// var theIDShort *string
// 	// var theID string
// 	// var theSpecificAssetIDs []types.ISpecificAssetID
// 	// var theSubmodelDescriptors []ISubmodelDescriptor

// 	// TODO: Finish
// 	return nil, nil
// }

// func SubmodelDescriptorFromJsonable(jsonable any) (ISubmodelDescriptor, error) {
// 	m, err := isValidJsonable(jsonable)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := submodelDescriptorFromMap(m)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func submodelDescriptorFromMap(m map[string]any) (result ISubmodelDescriptor, err error) {
// 	// TODO: Finish
// 	return nil, nil
// }

// func EndpointFromJsonable(jsonable any) (IEndpoint, error) {
// 	m, err := isValidJsonable(jsonable)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := endpointFromMap(m)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func endpointFromMap(m map[string]any) (result IEndpoint, err error) {
// 	var theInterface string
// 	var theProtocolInformation IProtocolInformation

// 	for k, v := range m {
// 		switch k {
// 		case "interface":
// 			parsedString, err := stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theInterface = parsedString
// 		case "protocolInformation":
// 			parsedProtocolInformation, err := ProtocolInformationFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theProtocolInformation = parsedProtocolInformation
// 		}
// 	}

// 	item := NewEndpoint(theInterface, theProtocolInformation)

// 	return item, nil
// }

// func ProtocolInformationFromJsonable(jsonable any) (IProtocolInformation, error) {
// 	m, err := isValidJsonable(jsonable)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := protocolInformationFromMap(m)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func protocolInformationFromMap(m map[string]any) (IProtocolInformation, error) {
// 	var theHref string
// 	var theEndpointProtocol *string
// 	var theEndpointProtocolVersion []string
// 	var theSubprotocol *string
// 	var theSubprotocolBody *string
// 	var theSubprotocolBodyEncoding *string
// 	var theSecurityAttributes []IProtocolInformationSecurityAttribute

// 	for k, v := range m {
// 		switch k {
// 		case "href":
// 			parsedString, err := stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theHref = parsedString

// 		case "endpointProtocol":
// 			parsedString, err := stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theEndpointProtocol = &parsedString

// 		case "endpointProtocolVersion":
// 			parsedArr, err := simpletonListFromJsonable(v, "string", "")
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theEndpointProtocolVersion = parsedArr

// 		case "subprotocol":
// 			parsedString, err := stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theSubprotocol = &parsedString

// 		case "subprotocolBody":
// 			parsedString, err := stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theSubprotocolBody = &parsedString

// 		case "subprotocolBodyEncoding":
// 			parsedString, err := stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theSubprotocolBodyEncoding = &parsedString
// 		case "securityAttributes":
// 			parsedArr, err := jsonizationFromJsonableList(v, ProtocolInformationSecurityAttributeFromJsonable)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 			theSecurityAttributes = parsedArr
// 		default:
// 			return nil, newDeserializationError(
// 				fmt.Sprintf("Unexpected property: %s", k),
// 			)
// 		}
// 	}

// 	item := NewProtocolInformation(theHref)

// 	item.SetEndpointProtocol(theEndpointProtocol)
// 	item.SetEndpointProtocolVersion(theEndpointProtocolVersion)
// 	item.SetSubprotocol(theSubprotocol)
// 	item.SetSubprotocolBody(theSubprotocolBody)
// 	item.SetSubprotocolBodyEncoding(theSubprotocolBodyEncoding)
// 	item.SetSecurityAttributes(theSecurityAttributes)

// 	return item, nil
// }

// func ProtocolInformationSecurityAttributeFromJsonable(jsonable any) (IProtocolInformationSecurityAttribute, error) {
// 	m, err := isValidJsonable(jsonable)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := protocolInformationSecurityInformationFromMap(m)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func protocolInformationSecurityInformationFromMap(m map[string]any) (result IProtocolInformationSecurityAttribute, err error) {
// 	var theType string
// 	var theKey string
// 	var theValue string

// 	for k, v := range m {
// 		switch k {
// 		case "type":
// 			theType, err = stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 		case "key":
// 			theKey, err = stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 		case "value":
// 			theValue, err = stringFromJsonable(v)
// 			if err != nil {
// 				return nil, prependName(err, k)
// 			}
// 		default:
// 			return nil, newDeserializationError(
// 				fmt.Sprintf("Unexpected property: %s", k),
// 			)
// 		}
// 	}

// 	return NewProtocolInformationSecurityAttribute(
// 		theType, theKey, theValue,
// 	), nil
// }
