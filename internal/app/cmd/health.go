// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/clivern/beetle/internal/app/util"

	log "github.com/sirupsen/logrus"
)

// HealthCheck controller
func HealthCheck() (bool, error) {
	log.WithFields(log.Fields{
		"CorrelationId": util.GenerateUUID4(),
	}).Info(`I am ok`)

	return true, nil
}
