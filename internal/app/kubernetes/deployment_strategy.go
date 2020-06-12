// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/util"
)

// Deploy deploys an application
func (c *Cluster) Deploy(deploymentRequest model.DeploymentRequest) (bool, error) {

	switch strategy := deploymentRequest.Strategy; strategy {

	case model.RecreateStrategy:
		return c.RecreateStrategy(deploymentRequest)

	case model.RampedStrategy:
		return c.RampedStrategy(deploymentRequest)

	case model.CanaryStrategy:
		return c.CanaryStrategy(deploymentRequest)

	case model.BlueGreenStrategy:
		return c.BlueGreenStrategy(deploymentRequest)

	default:
		return false, fmt.Errorf("Invalid deployment strategy %s", strategy)
	}
}

// RecreateStrategy terminates the old version and release the new one.
//
// This method is like running this command
//
// ```
//      $ kubectl patch deployment toad-deployment --type=json -p '[
//          {"op":"replace", "path":"/spec/strategy", "value":{"type":"Recreate"}},
//          {"op":"replace","path":"/spec/template/spec/containers/0/image","value":"clivern/toad:release-0.2.4"}
//      ]'
// ```
//
// In order to use it
//
//```
//      cluster.RecreateStrategy(model.DeploymentRequest{
//      	Cluster:     "~~cluster~~",
//          Namespace:   "~~namespace",
//          Application: "~~application id~~",
//          Version:     ~~application new version~~",
//          Strategy:    "recreate",
//      })
// ```
func (c *Cluster) RecreateStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	result := Application{}
	patch := make(map[string][]model.PatchStringValue)

	config, err := c.GetConfig(context.TODO(), deploymentRequest.Namespace)

	if err != nil {
		return false, err
	}

	for _, app := range config.Applications {
		if app.ID == deploymentRequest.Application {
			result, err = c.GetApplication(
				context.TODO(),
				deploymentRequest.Namespace,
				app.ID,
				app.Name,
				app.ImageFormat,
			)
			if err != nil {
				return false, err
			}
			break
		}
	}

	i := 0
	for _, container := range result.Containers {
		if _, ok := patch[container.Deployment.Name]; !ok {
			patch[container.Deployment.Name] = []model.PatchStringValue{}
		}
		patch[container.Deployment.Name] = append(patch[container.Deployment.Name], model.PatchStringValue{
			Op:    "replace",
			Path:  fmt.Sprintf("/spec/template/spec/containers/%d/image", i),
			Value: strings.Replace(container.Image, container.Version, deploymentRequest.Version, -1),
		})
		i++
	}

	data := ""
	status := true

	for deploymentName, deploymentPatch := range patch {
		data, err = util.ConvertToJSON(deploymentPatch)

		if err != nil {
			return false, err
		}

		// Enforce Recreate strategy
		data = strings.Replace(
			data,
			`[`, `[{"op":"replace","path":"/spec/strategy","value":{"type":"Recreate"}},`,
			-1,
		)

		status, err = c.PatchDeployment(
			context.TODO(),
			deploymentRequest.Namespace,
			deploymentName,
			data,
		)

		if !status || err != nil {
			return false, err
		}
	}

	for deploymentName := range patch {
		status, err = c.FetchDeploymentStatus(context.TODO(), deploymentRequest.Namespace, deploymentName, 600)

		if !status || err != nil {
			return false, err
		}
	}

	return true, nil
}

// RampedStrategy releases a new version on a rolling update fashion, one after the other.
//
// it will set maxSurge as 25% and maxUnavailable as 25%
//
// This method is like running this command
//
// ```
//      $ kubectl patch deployment toad-deployment --type=json -p '[
//          {"op":"replace", "path":"/spec/strategy", "value":{"type":"RollingUpdate"}},
//          {"op":"replace","path":"/spec/template/spec/containers/0/image","value":"clivern/toad:release-0.2.4"}
//      ]'
// ```
//
// In order to use it
//
//```
//      cluster.RampedStrategy(model.DeploymentRequest{
//      	Cluster:     "~~cluster~~",
//          Namespace:   "~~namespace",
//          Application: "~~application id~~",
//          Version:     ~~application new version~~",
//          Strategy:    "ramped",
//      })
// ```
func (c *Cluster) RampedStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	result := Application{}
	patch := make(map[string][]model.PatchStringValue)

	config, err := c.GetConfig(context.TODO(), deploymentRequest.Namespace)

	if err != nil {
		return false, err
	}

	for _, app := range config.Applications {
		if app.ID == deploymentRequest.Application {
			result, err = c.GetApplication(
				context.TODO(),
				deploymentRequest.Namespace,
				app.ID,
				app.Name,
				app.ImageFormat,
			)
			if err != nil {
				return false, err
			}
			break
		}
	}

	i := 0
	for _, container := range result.Containers {
		if _, ok := patch[container.Deployment.Name]; !ok {
			patch[container.Deployment.Name] = []model.PatchStringValue{}
		}
		patch[container.Deployment.Name] = append(patch[container.Deployment.Name], model.PatchStringValue{
			Op:    "replace",
			Path:  fmt.Sprintf("/spec/template/spec/containers/%d/image", i),
			Value: strings.Replace(container.Image, container.Version, deploymentRequest.Version, -1),
		})
		i++
	}

	data := ""
	status := true

	for deploymentName, deploymentPatch := range patch {
		data, err = util.ConvertToJSON(deploymentPatch)

		if err != nil {
			return false, err
		}

		// Enforce RollingUpdate strategy
		data = strings.Replace(
			data,
			`[`, `[{"op":"replace","path":"/spec/strategy","value":{"type":"RollingUpdate"}},`,
			-1,
		)

		status, err = c.PatchDeployment(
			context.TODO(),
			deploymentRequest.Namespace,
			deploymentName,
			data,
		)

		if !status || err != nil {
			return false, err
		}

	}

	for deploymentName := range patch {
		status, err = c.FetchDeploymentStatus(context.TODO(), deploymentRequest.Namespace, deploymentName, 1000)

		if !status || err != nil {
			return false, err
		}
	}

	return true, nil
}

// BlueGreenStrategy releases a new version alongside the old version then switch traffic.
func (c *Cluster) BlueGreenStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	return true, nil
}

// CanaryStrategy releases a new version to a subset of users, then proceed to a full rollout.
func (c *Cluster) CanaryStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	return true, nil
}
