// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/clivern/beetle/core/module"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Jobs controller
func Jobs(c *gin.Context) {
	db := module.Database{}

	err := db.AutoConnect()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Error(`Failure while connecting database`)

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	c.JSON(http.StatusOK, gin.H{
		"jobs": db.GetJobs(),
	})
}
