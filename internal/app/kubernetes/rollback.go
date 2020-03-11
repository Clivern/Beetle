// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

// Rollback struct
type Rollback struct {
	Cluster     string
	Namespace   string
	Application string
	Current     string
	Target      string
}

// Run runs the rollback process
func (d *Rollback) Run() (bool, error) {
	return true, nil
}
