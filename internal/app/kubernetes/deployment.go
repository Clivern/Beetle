// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

// Deployment struct
type Deployment struct {
	Cluster     string
	Namespace   string
	Application string
	Current     string
	Target      string
}

// Run runs the deployment process
func (d *Deployment) Run() (bool, error) {
	return true, nil
}
