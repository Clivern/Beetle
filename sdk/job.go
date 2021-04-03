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

// GetJobs Get Jobs List
func (c *Client) GetJobs(ctx context.Context) (model.Jobs, error) {
	var result model.Jobs

	response, err := c.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/job", c.APIURL),
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

// GetJob Get Job
func (c *Client) GetJob(ctx context.Context, uuid string) (model.Job, error) {
	var result model.Job

	response, err := c.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/api/v1/job/%s", c.APIURL, uuid),
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

// DeleteJob Delete Job
func (c *Client) DeleteJob(ctx context.Context, uuid string) (bool, error) {
	response, err := c.HTTPClient.Delete(
		ctx,
		fmt.Sprintf("%s/api/v1/job/%s", c.APIURL, uuid),
		map[string]string{},
		map[string]string{"X-API-KEY": c.APIKey},
	)

	if err != nil {
		return false, err
	}

	statusCode := c.HTTPClient.GetStatusCode(response)

	if statusCode != http.StatusNoContent {
		return false, fmt.Errorf(fmt.Sprintf("Invalid status code %d", statusCode))
	}

	return true, nil
}
