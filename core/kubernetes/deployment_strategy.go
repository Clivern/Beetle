// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/clivern/beetle/core/model"
	"github.com/clivern/beetle/core/util"
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
// $ kubectl patch deployment toad-deployment --type=json -p '[
//     {"op":"replace", "path":"/spec/strategy", "value":{"type":"Recreate"}},
//     {"op":"replace","path":"/spec/template/spec/containers/0/image","value":"clivern/toad:release-0.2.4"}
// ]'
func (c *Cluster) RecreateStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	result := model.Application{}
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
			`[`,
			`[{"op":"replace","path":"/spec/strategy","value":{"type":"Recreate"}},`,
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
// $ kubectl patch deployment toad-deployment --type=json -p '[
//     {"op":"replace", "path":"/spec/strategy", "value":{"type":"RollingUpdate"}},
//     {"op":"replace", "path":"/spec/strategy/rollingUpdate", "value":{"maxSurge":""}},
//	   {"op":"replace", "path":"/spec/strategy/rollingUpdate", "value":{"maxUnavailable":""}},
//     {"op":"replace","path":"/spec/template/spec/containers/0/image","value":"clivern/toad:release-0.2.4"}
// ]'
func (c *Cluster) RampedStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	result := model.Application{}
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

		diff := ""

		if strings.Contains(deploymentRequest.MaxSurge, "%") && strings.Contains(deploymentRequest.MaxUnavailable, "%") {
			diff = fmt.Sprintf(
				`[{"op":"replace","path":"/spec/strategy","value":{"type":"RollingUpdate"}},`+
					`{"op":"replace", "path":"/spec/strategy/rollingUpdate", "value":{"maxSurge":"%s"}},`+
					`{"op":"replace", "path":"/spec/strategy/rollingUpdate", "value":{"maxUnavailable":"%s"}},`,
				deploymentRequest.MaxSurge,
				deploymentRequest.MaxUnavailable,
			)
		} else {
			maxSurge, err := strconv.Atoi(deploymentRequest.MaxSurge)

			if err != nil {
				return false, err
			}

			maxUnavailable, err := strconv.Atoi(deploymentRequest.MaxUnavailable)

			if err != nil {
				return false, err
			}

			diff = fmt.Sprintf(
				`[{"op":"replace","path":"/spec/strategy","value":{"type":"RollingUpdate"}},`+
					`{"op":"replace", "path":"/spec/strategy/rollingUpdate", "value":{"maxSurge":%d}},`+
					`{"op":"replace", "path":"/spec/strategy/rollingUpdate", "value":{"maxUnavailable":%d}},`,
				maxSurge,
				maxUnavailable,
			)
		}

		// Enforce RollingUpdate strategy
		data = strings.Replace(
			data,
			`[`,
			diff,
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
func (c *Cluster) BlueGreenStrategy(_ model.DeploymentRequest) (bool, error) {
	return true, nil
}

// CanaryStrategy releases a new version to a subset of users, then proceed to a full rollout.
func (c *Cluster) CanaryStrategy(_ model.DeploymentRequest) (bool, error) {
	return true, nil
}
