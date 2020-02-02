// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"github.com/nbio/st"
	"net/http"
	"strings"
	"testing"
)

// TestHttpGet test cases
func TestHttpGet(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Get(
		"https://httpbin.org/get",
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, true, strings.Contains(body, "value1"))
	st.Expect(t, true, strings.Contains(body, "arg1"))
	st.Expect(t, true, strings.Contains(body, "arg1=value1"))
	st.Expect(t, true, strings.Contains(body, "X-Auth"))
	st.Expect(t, true, strings.Contains(body, "hipp-123"))
	st.Expect(t, nil, error)
}

// TestHttpDelete test cases
func TestHttpDelete(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Delete(
		"https://httpbin.org/delete",
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, true, strings.Contains(body, "value1"))
	st.Expect(t, true, strings.Contains(body, "arg1"))
	st.Expect(t, true, strings.Contains(body, "arg1=value1"))
	st.Expect(t, true, strings.Contains(body, "X-Auth"))
	st.Expect(t, true, strings.Contains(body, "hipp-123"))
	st.Expect(t, nil, error)
}

// TestHttpPost test cases
func TestHttpPost(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Post(
		"https://httpbin.org/post",
		`{"Username":"admin", "Password":"12345"}`,
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, true, strings.Contains(body, `"12345"`))
	st.Expect(t, true, strings.Contains(body, `"Username"`))
	st.Expect(t, true, strings.Contains(body, `"admin"`))
	st.Expect(t, true, strings.Contains(body, `"Password"`))
	st.Expect(t, true, strings.Contains(body, "value1"))
	st.Expect(t, true, strings.Contains(body, "arg1"))
	st.Expect(t, true, strings.Contains(body, "arg1=value1"))
	st.Expect(t, true, strings.Contains(body, "X-Auth"))
	st.Expect(t, true, strings.Contains(body, "hipp-123"))
	st.Expect(t, nil, error)
}

// TestHttpPut test cases
func TestHttpPut(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Put(
		"https://httpbin.org/put",
		`{"Username":"admin", "Password":"12345"}`,
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, true, strings.Contains(body, `"12345"`))
	st.Expect(t, true, strings.Contains(body, `"Username"`))
	st.Expect(t, true, strings.Contains(body, `"admin"`))
	st.Expect(t, true, strings.Contains(body, `"Password"`))
	st.Expect(t, true, strings.Contains(body, "value1"))
	st.Expect(t, true, strings.Contains(body, "arg1"))
	st.Expect(t, true, strings.Contains(body, "arg1=value1"))
	st.Expect(t, true, strings.Contains(body, "X-Auth"))
	st.Expect(t, true, strings.Contains(body, "hipp-123"))
	st.Expect(t, nil, error)
}

// TestHttpGetStatusCode1 test cases
func TestHttpGetStatusCode1(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Get(
		"https://httpbin.org/status/200",
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, "", body)
	st.Expect(t, nil, error)
}

// TestHttpGetStatusCode2 test cases
func TestHttpGetStatusCode2(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Get(
		"https://httpbin.org/status/500",
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusInternalServerError, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, "", body)
	st.Expect(t, nil, error)
}

// TestHttpGetStatusCode3 test cases
func TestHttpGetStatusCode3(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Get(
		"https://httpbin.org/status/404",
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusNotFound, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, "", body)
	st.Expect(t, nil, error)
}

// TestHttpGetStatusCode4 test cases
func TestHttpGetStatusCode4(t *testing.T) {
	httpClient := NewHTTPClient()
	response, error := httpClient.Get(
		"https://httpbin.org/status/201",
		map[string]string{"arg1": "value1"},
		map[string]string{"X-Auth": "hipp-123"},
	)
	t.Log(httpClient.GetStatusCode(response))
	st.Expect(t, http.StatusCreated, httpClient.GetStatusCode(response))
	st.Expect(t, nil, error)

	body, error := httpClient.ToString(response)
	t.Log(body)
	st.Expect(t, "", body)
	st.Expect(t, nil, error)
}

// TestBuildParameters test cases
func TestBuildParameters(t *testing.T) {
	httpClient := NewHTTPClient()
	url, error := httpClient.BuildParameters("http://127.0.0.1", map[string]string{"arg1": "value1"})
	t.Log(url)
	st.Expect(t, "http://127.0.0.1?arg1=value1", url)
	st.Expect(t, nil, error)
}

// TestBuildData test cases
func TestBuildData(t *testing.T) {
	httpClient := NewHTTPClient()
	st.Expect(t, httpClient.BuildData(map[string]string{}), "")
	st.Expect(t, httpClient.BuildData(map[string]string{"arg1": "value1"}), "arg1=value1")
}
