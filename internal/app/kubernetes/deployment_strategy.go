// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"fmt"

	"github.com/clivern/beetle/internal/app/model"
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
func (c *Cluster) RecreateStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
	return true, nil
}

// RampedStrategy releases a new version on a rolling update fashion, one after the other.
func (c *Cluster) RampedStrategy(deploymentRequest model.DeploymentRequest) (bool, error) {
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
