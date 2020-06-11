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
// This method is like running this command
//
// ```
//    	$ kubectl patch deployment toad-deployment --type=json -p '[
// 	  	  	{"op":"replace", "path":"/spec/strategy", "value":{"type":"Recreate"}},
//   		{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"clivern/toad:release-0.2.3"}
//     	]'
// ```
// In order to use it
//
//```
// 		fmt.Println(cluster.RecreateStrategy(model.DeploymentRequest{
//			Cluster:     $$cluster,
//			Namespace:   $$namespace,
//			Application: $$application_id,
//			Version:     $$application_new_version,
//			Strategy:    "recreate",
//		}))
// ```
func (c *Cluster) RecreateStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	result := Application{}
	patch := make(map[string][]model.PatchStringValue)

	config, err := c.GetConfig(context.Background(), deploymentRequest.Namespace)

	if err != nil {
		return false, err
	}

	for _, app := range config.Applications {
		if app.ID == deploymentRequest.Application {
			result, err = c.GetApplication(
				context.Background(),
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
		i += 1
	}

	data := ""
	status := true

	for deployment_name, deployment_patch := range patch {
		data, err = util.ConvertToJSON(deployment_patch)

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
			context.Background(),
			deploymentRequest.Namespace,
			deployment_name,
			data,
		)

		if !status || err != nil {
			return false, err
		}
	}

	// TODO -----> watch deployment status to ensure it is fully succeeded

	return true, nil
}

// RampedStrategy releases a new version on a rolling update fashion, one after the other.
func (c *Cluster) RampedStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	// deploymentRequest.Namespace
	// deploymentRequest.Application
	// deploymentRequest.Version
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

// FetchDeploymentStatus get deployment status
func (c *Cluster) FetchDeploymentStatus() {
	// Wait for the deployment to check the final
	// Status https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#deployment-status
}
