package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/SmartFactory-KL/shellgonaut/create"
	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
	"github.com/aas-core-works/aas-core3.1-golang/stringification"
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

// aasEntityListToBytes applies aasEntityToBytes to all items within entityList
// while still returning a single []byte
// containing the json list
func aasEntityListToBytes[T types.IClass](entityList []T) ([]byte, error) {
	if len(entityList) == 0 {
		return nil, fmt.Errorf("list cannot be nil or empty")
	}

	result := make([]json.RawMessage, 0, len(entityList))

	for _, entityItem := range entityList {
		item, err := aasEntityToBytes(entityItem)
		if err != nil {
			// TODO: Similar to other places: Should a single invalid item really
			// TODO: stop the whole operation?
			return nil, fmt.Errorf("failed to convert list item: %w", err)
		}

		result = append(result, item)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal list: %w", err)
	}

	return resultJson, nil
}

// createReferenceForSubmodel creates a reference that can be attached to the "submodels" property of a shell
func createReferenceForSubmodel(submodel types.ISubmodel) (types.IReference, error) {
	if submodel == nil {
		return nil, fmt.Errorf("submodel cannot be nil")
	}

	submodelID := submodel.ID()
	if submodelID == "" {
		return nil, fmt.Errorf("submodel ID cannot be empty")
	}

	reference := types.NewReference(
		types.ReferenceTypesModelReference,
		[]types.IKey{
			types.NewKey(types.KeyTypesSubmodel, submodelID),
		},
	)

	if semanticID := submodel.SemanticID(); semanticID != nil {
		refCopy, err := create.CreateCopy(semanticID, jsonization.ReferenceFromJsonable)
		if err != nil {
			return nil, fmt.Errorf("failed to create reference copy: %w", err)
		}
		reference.SetReferredSemanticID(refCopy)
	}

	return reference, nil
}

// getEncodedTargetUrl creates a url by appending new path elements to baseURL, encoding only idToEncode to base64URL
func getEncodedTargetUrl(baseUrl *url.URL, preIDPath string, idToEncode string, postIDPath ...string) (*url.URL, error) {
	if len(strings.TrimSpace(idToEncode)) == 0 {
		return nil, fmt.Errorf("id cannot be empty")
	}

	encodedID := toBase64URL(idToEncode)

	parts := append([]string{preIDPath, encodedID}, postIDPath...)
	targetUrl := baseUrl.JoinPath(parts...)
	return targetUrl, nil
}

// toBase64URL encodes []byte(input) into base64URL string
func toBase64URL(input string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(input))
}

// EnsureSubmodelElementType takes any ISubmodelElement and checks wether it fits the new type T as well as the modelType.
// This returns a more usable subtype than only ISubmodelElement
func EnsureSubmodelElementType[T types.IClass](element types.ISubmodelElement, targetModelType types.ModelType) (T, error) {
	var zero T

	if element == nil {
		return zero, fmt.Errorf("element cannot be nil")
	}

	if element.ModelType() != targetModelType {
		return zero, fmt.Errorf(
			"expected modelType %s but got %s",
			stringification.MustModelTypeToString(targetModelType),
			stringification.MustModelTypeToString(element.ModelType()),
		)
	}

	typedElement, ok := element.(T)
	if !ok {
		return zero, fmt.Errorf("type assertion failed")
	}

	return typedElement, nil
}

// EnsureUrlWithoutSuffixOrSlash normalizes and validates an HTTP(S) URL.
//
// It trims surrounding whitespace and removes trailing slashes. If urlInput ends
// with suffix, it removes the suffix and any resulting trailing slashes. It
// returns an error if the resulting URL is empty, invalid, lacks a scheme or
// host, or does not use the http or https scheme.
func EnsureUrlWithoutSuffixOrSlash(urlInput string, suffix string) (*url.URL, error) {
	urlInput = strings.TrimSpace(urlInput)
	urlInput = strings.TrimRight(urlInput, "/")

	if before, ok := strings.CutSuffix(urlInput, suffix); ok {
		urlInput = before
	}

	urlInput = strings.TrimRight(urlInput, "/")

	if len(urlInput) == 0 {
		return nil, fmt.Errorf("url cannot be empty")
	}

	urlResult, err := url.Parse(urlInput)
	if err != nil {
		return nil, fmt.Errorf("url parsing failed: %w", err)
	}

	if len(urlResult.Host) == 0 || len(urlResult.Scheme) == 0 {
		return nil, fmt.Errorf("url needs scheme and host, one is missing: %q", urlResult.String())
	}

	if urlResult.Scheme != "http" && urlResult.Scheme != "https" {
		return nil, fmt.Errorf("only http/https supported as scheme")
	}

	return urlResult, nil
}
