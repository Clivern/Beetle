// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/clivern/beetle/internal/app/module"
)

// HealthCheck controller
func HealthCheck() (bool, error) {
	logger, _ := module.NewLogger()

	defer logger.Sync()

	logger.Info("I am ok")

	return true, nil
}
