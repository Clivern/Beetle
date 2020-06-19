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

// GetApplications Get Applications List
func GetApplications(ctx context.Context, httpClient *module.HTTPClient, serverURL, cluster, namespace, apiKey string) (model.Applications, error) {
	var result model.Applications

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s/namespace/%s/app", serverURL, cluster, namespace),
		map[string]string{},
		map[string]string{"X-API-KEY": apiKey},
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

// GetApplication Get Application
func GetApplication(ctx context.Context, httpClient *module.HTTPClient, serverURL, cluster, namespace, application, apiKey string) (model.Application, error) {
	var result model.Application

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/cluster/%s/namespace/%s/app/%s", serverURL, cluster, namespace, application),
		map[string]string{},
		map[string]string{"X-API-KEY": apiKey},
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
