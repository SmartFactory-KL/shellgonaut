package client

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/aas-core-works/aas-core3.1-golang/stringification"
	"github.com/aas-core-works/aas-core3.1-golang/types"
)

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
