// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clivern/beetle/core/model"
)

// GetClusters Get Clusters List
func (c *Client) GetClusters(ctx context.Context) (model.Clusters, error) {
	var result model.Clusters

	response, err := c.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster", c.APIURL),
		map[string]string{},
		map[string]string{"X-API-KEY": c.APIKey},
	)

	if err != nil {
		return result, err
	}

	statusCode := c.HTTPClient.GetStatusCode(response)

	if statusCode != http.StatusOK {
		return result, fmt.Errorf(fmt.Sprintf("Invalid status code %d", statusCode))
	}

	body, err := c.HTTPClient.ToString(response)

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
func (c *Client) GetCluster(ctx context.Context, cluster string) (model.Cluster, error) {
	var result model.Cluster

	response, err := c.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s", c.APIURL, cluster),
		map[string]string{},
		map[string]string{"X-API-KEY": c.APIKey},
	)

	if err != nil {
		return result, err
	}

	statusCode := c.HTTPClient.GetStatusCode(response)

	if statusCode != http.StatusOK {
		return result, fmt.Errorf(fmt.Sprintf("Invalid status code %d", statusCode))
	}

	body, err := c.HTTPClient.ToString(response)

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
