// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/model"
)

// GetApplications Get Applications List
func (c *Client) GetApplications(ctx context.Context, cluster, namespace string) (model.Applications, error) {
	var result model.Applications

	response, err := c.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s/namespace/%s/app", c.APIURL, cluster, namespace),
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

// GetApplication Get Application
func (c *Client) GetApplication(ctx context.Context, cluster, namespace, application string) (model.Application, error) {
	var result model.Application

	response, err := c.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s/namespace/%s/app/%s", c.APIURL, cluster, namespace, application),
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
