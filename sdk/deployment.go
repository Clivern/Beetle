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

// CreateDeployment Get Application
func CreateDeployment(ctx context.Context, httpClient *module.HTTPClient, serverURL string, request model.DeploymentRequest, apiKey string) (model.Job, error) {
	var result model.Job

	requestBody, err := request.ConvertToJSON()

	if err != nil {
		return result, err
	}

	response, err := httpClient.Post(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s/namespace/%s/app/%s", serverURL, request.Cluster, request.Namespace, request.Application),
		requestBody,
		map[string]string{},
		map[string]string{"X-API-KEY": apiKey},
	)

	if err != nil {
		return result, err
	}

	statusCode := httpClient.GetStatusCode(response)

	if statusCode != http.StatusAccepted {
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
