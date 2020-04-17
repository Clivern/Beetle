// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthCheck controller
func HealthCheck(c *gin.Context) {
	status := "ok"

	db := module.Database{}

	err := db.AutoConnect()

	if err != nil {
		status = "down"

		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Error(fmt.Sprintf(`Error: %s`, err.Error()))

		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Info(fmt.Sprintf(`Health Status: %s`, status))

		c.Status(http.StatusInternalServerError)
		return
	}

	err = db.Ping()

	if err != nil {
		status = "down"

		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Error(fmt.Sprintf(`Error: %s`, err.Error()))

		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Info(fmt.Sprintf(`Health Status: %s`, status))

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	log.WithFields(log.Fields{
		"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
	}).Info(fmt.Sprintf(`Health Status: %s`, status))

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
