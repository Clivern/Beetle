// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"
)

// GetNamespaces Get Namespaces List
func GetNamespaces(ctx context.Context, httpClient *module.HTTPClient, serverURL, cluster, token string) (model.Namespaces, error) {
	var result model.Namespaces

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s/namespace", serverURL, cluster),
		map[string]string{},
		map[string]string{"X-AUTH-TOKEN": token},
	)

	if err != nil {
		return result, err
	}

	statusCode := httpClient.GetStatusCode(response)

	if statusCode != http.StatusOK {
		return result, fmt.Errorf(fmt.Sprintf("Invalid status code %d", statusCode))
	}

	body, err := httpClient.ToString(response)

	if err != nil {
		return result, fmt.Errorf(fmt.Sprintf("Invalid response: %s", err.Error()))
	}

	ok, err := result.LoadFromJSON([]byte(body))

	if err != nil {
		return result, fmt.Errorf(fmt.Sprintf("Invalid response: %s", err.Error()))
	}

	if !ok {
		return result, fmt.Errorf("Invalid response")
	}

	return result, nil
}

// GetNamespace Get Namespace
func GetNamespace(ctx context.Context, httpClient *module.HTTPClient, serverURL, cluster, namespace, token string) (model.Namespace, error) {
	var result model.Namespace

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s/namespace/%s", serverURL, cluster, namespace),
		map[string]string{},
		map[string]string{"X-AUTH-TOKEN": token},
	)

	if err != nil {
		return result, err
	}

	statusCode := httpClient.GetStatusCode(response)

	if statusCode != http.StatusOK {
		return result, fmt.Errorf(fmt.Sprintf("Invalid status code %d", statusCode))
	}

	body, err := httpClient.ToString(response)

	if err != nil {
		return result, fmt.Errorf(fmt.Sprintf("Invalid response: %s", err.Error()))
	}

	ok, err := result.LoadFromJSON([]byte(body))

	if err != nil {
		return result, fmt.Errorf(fmt.Sprintf("Invalid response: %s", err.Error()))
	}

	if !ok {
		return result, fmt.Errorf("Invalid response")
	}

	return result, nil
}
