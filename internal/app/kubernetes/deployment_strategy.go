// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

// RecreateStrategy terminates the old version and release the new one.
func (c *Cluster) RecreateStrategy() (bool, error) {
	return true, nil
}

// RampedStrategy releases a new version on a rolling update fashion, one after the other.
func (c *Cluster) RampedStrategy() (bool, error) {
	return true, nil
}

// BlueGreenStrategy releases a new version alongside the old version then switch traffic.
func (c *Cluster) BlueGreenStrategy() (bool, error) {
	return true, nil
}

// CanaryStrategy releases a new version to a subset of users, then proceed to a full rollout.
func (c *Cluster) CanaryStrategy() (bool, error) {
	return true, nil
}
