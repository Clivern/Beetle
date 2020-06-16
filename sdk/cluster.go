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

// GetClusters Get Clusters List
func GetClusters(ctx context.Context, httpClient *module.HTTPClient, serverURL, token string) (model.Clusters, error) {
	var result model.Clusters

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster", serverURL),
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

// GetCluster Get Cluster
func GetCluster(ctx context.Context, httpClient *module.HTTPClient, serverURL, cluster, token string) (model.Cluster, error) {
	var result model.Cluster

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s", serverURL, cluster),
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
