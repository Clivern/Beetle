// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthCheck controller
func HealthCheck(c *gin.Context) {
	status := "ok"

	log.WithFields(log.Fields{
		"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
	}).Info(fmt.Sprintf(`Health Status: %s`, status))

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
