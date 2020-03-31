// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	log "github.com/sirupsen/logrus"
)

// HealthCheck controller
func HealthCheck() (bool, error) {
	log.WithFields(log.Fields{
		"CorrelationId": "",
	}).Info(`I am ok`)

	return true, nil
}
