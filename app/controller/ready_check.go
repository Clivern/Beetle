// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/clivern/beetle/app/module"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ReadyCheck controller
func ReadyCheck(c *gin.Context) {
	status := "ok"

	db := module.Database{}

	err := db.AutoConnect()

	if err != nil {
		status = "down"

		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"status":         status,
			"error":          err.Error(),
		}).Error(`Failed ready check`)

		c.Status(http.StatusInternalServerError)
		return
	}

	err = db.Ping()

	if err != nil {
		status = "down"

		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"status":         status,
			"error":          err.Error(),
		}).Error(`Failed ready check`)

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	log.WithFields(log.Fields{
		"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		"status":         status,
	}).Info(`Passed ready check`)

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
