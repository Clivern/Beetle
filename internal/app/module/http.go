// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient struct
type HTTPClient struct{}

// NewHTTPClient creates an instance of http client
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

// Get http call
func (h *HTTPClient) Get(ctx context.Context, endpoint string, parameters map[string]string, headers map[string]string) (*http.Response, error) {

	endpoint, err := h.BuildParameters(endpoint, parameters)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", endpoint, nil)

	req = req.WithContext(ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// Post http call
func (h *HTTPClient) Post(ctx context.Context, endpoint string, data string, parameters map[string]string, headers map[string]string) (*http.Response, error) {

	endpoint, err := h.BuildParameters(endpoint, parameters)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(data)))

	req = req.WithContext(ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// Put http call
func (h *HTTPClient) Put(ctx context.Context, endpoint string, data string, parameters map[string]string, headers map[string]string) (*http.Response, error) {

	endpoint, err := h.BuildParameters(endpoint, parameters)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("PUT", endpoint, bytes.NewBuffer([]byte(data)))

	req = req.WithContext(ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// Delete http call
func (h *HTTPClient) Delete(ctx context.Context, endpoint string, parameters map[string]string, headers map[string]string) (*http.Response, error) {

	endpoint, err := h.BuildParameters(endpoint, parameters)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("DELETE", endpoint, nil)

	req = req.WithContext(ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// BuildParameters add parameters to URL
func (h *HTTPClient) BuildParameters(endpoint string, parameters map[string]string) (string, error) {
	u, err := url.Parse(endpoint)

	if err != nil {
		return "", err
	}

	q := u.Query()

	for k, v := range parameters {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// BuildData build body data
func (h *HTTPClient) BuildData(parameters map[string]string) string {
	var items []string

	for k, v := range parameters {
		items = append(items, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(items, "&")
}

// ToString response body to string
func (h *HTTPClient) ToString(response *http.Response) (string, error) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetStatusCode response status code
func (h *HTTPClient) GetStatusCode(response *http.Response) int {
	return response.StatusCode
}
