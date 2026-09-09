package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
)

type QueryItem struct {
	Key   string
	Value string
}

func DoPagedGetRequest(httpClient *http.Client, targetUrl string, cursor string, limit int, additionalQueryItems ...QueryItem) (*PagedResultRaw, error) {
	// add possible query params
	parsedURL, err := url.Parse(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse targed URL: %w", err)
	}

	query := parsedURL.Query()
	if len(cursor) > 0 {
		query.Set("cursor", cursor)
	}

	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	if len(additionalQueryItems) > 0 {
		for _, queryItem := range additionalQueryItems {
			query.Set(queryItem.Key, queryItem.Value)
		}
	}

	parsedURL.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to craft Paged GET request: %w", err)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Paged GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// read raw result
		var pagedResultRaw PagedResultRaw
		if err := json.NewDecoder(resp.Body).Decode(&pagedResultRaw); err != nil {
			return nil, fmt.Errorf("request returned %s but decoding JSON failed: %w", resp.Status, err)
		}

		return &pagedResultRaw, nil
	} else {
		// read error
		var errorResult ErrorResult
		result, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("request failed with status %s and reading response also failed: %w", resp.Status, err)
		}

		if err := json.Unmarshal(result, &errorResult); err != nil {
			return nil, fmt.Errorf(
				"request failed with status %s but decoding error json also failed with %w for response %s",
				resp.Status,
				err,
				string(result),
			)
		}

		return nil, errorResult
	}
}

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
