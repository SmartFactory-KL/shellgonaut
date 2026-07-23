package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
)

func DoGetRequest(httpClient *http.Client, targetUrl string) (io.ReadCloser, error) {
	return doRequest(httpClient, http.MethodGet, targetUrl, nil)
}

func DoPostRequest(httpClient *http.Client, targetUrl string, content []byte) (io.ReadCloser, error) {
	return doRequest(httpClient, http.MethodPost, targetUrl, bytes.NewReader(content))
}

func DoPatchRequest(httpClient *http.Client, targetUrl string, content []byte) (io.ReadCloser, error) {
	return doRequest(httpClient, http.MethodPatch, targetUrl, bytes.NewReader(content))
}

func DoPutRequest(httpClient *http.Client, targetUrl string, content []byte) (io.ReadCloser, error) {
	return doRequest(httpClient, http.MethodPut, targetUrl, bytes.NewReader(content))
}

func DoDeleteRequest(httpClient *http.Client, targetUrl string) (io.ReadCloser, error) {
	return doRequest(httpClient, http.MethodDelete, targetUrl, nil)
}

func doRequest(httpClient *http.Client, method string, targetUrl string, body io.Reader) (io.ReadCloser, error) {
	request, err := http.NewRequest(method, targetUrl, body)
	if err != nil {
		return nil, fmt.Errorf("failed to craft %s request: %w", method, err)
	}

	if body != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}

	return readResponseBody(response)
}

func readResponseBody(response *http.Response) (io.ReadCloser, error) {
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return response.Body, nil
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusConflict:
		return nil, ErrAlreadyExists
	default:
		return nil, fmt.Errorf("request returned unexpected status %s", response.Status)
	}
}
