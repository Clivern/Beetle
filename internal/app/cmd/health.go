// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/internal/app/util"

	log "github.com/sirupsen/logrus"
)

// HealthCheck controller
func HealthCheck() {
	status := "ok"

	db := module.Database{}

	err := db.AutoConnect()

	if err != nil {
		status = "down"

		log.WithFields(log.Fields{
			"CorrelationId": util.GenerateUUID4(),
		}).Error(fmt.Sprintf(`Error: %s`, err.Error()))

		log.WithFields(log.Fields{
			"CorrelationId": util.GenerateUUID4(),
		}).Info(fmt.Sprintf(`Health Status: %s`, status))

		panic(err.Error())
	}

	err = db.Ping()

	if err != nil {
		status = "down"

		log.WithFields(log.Fields{
			"CorrelationId": util.GenerateUUID4(),
		}).Error(fmt.Sprintf(`Error: %s`, err.Error()))

		log.WithFields(log.Fields{
			"CorrelationId": util.GenerateUUID4(),
		}).Info(fmt.Sprintf(`Health Status: %s`, status))

		panic(err.Error())
	}

	defer db.Close()

	log.WithFields(log.Fields{
		"CorrelationId": util.GenerateUUID4(),
	}).Info(fmt.Sprintf(`Health Status: %s`, status))
}
