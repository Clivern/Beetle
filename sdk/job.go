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

// GetJobs Get Jobs List
func GetJobs(ctx context.Context, httpClient *module.HTTPClient, serverURL, apiKey string) (model.Jobs, error) {
	var result model.Jobs

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/job", serverURL),
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

// GetJob Get Job
func GetJob(ctx context.Context, httpClient *module.HTTPClient, serverURL, uuid, apiKey string) (model.Job, error) {
	var result model.Job

	response, err := httpClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/job/%s", serverURL, uuid),
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

// DeleteJob Delete Job
func DeleteJob(ctx context.Context, httpClient *module.HTTPClient, serverURL, uuid, apiKey string) (bool, error) {
	response, err := httpClient.Delete(
		ctx,
		fmt.Sprintf("%s/api/v1/job/%s", serverURL, uuid),
		map[string]string{},
		map[string]string{"X-API-KEY": apiKey},
	)

	if err != nil {
		return false, err
	}

	statusCode := httpClient.GetStatusCode(response)

	if statusCode != http.StatusNoContent {
		return false, fmt.Errorf(fmt.Sprintf("Invalid status code %d", statusCode))
	}

	return true, nil
}
