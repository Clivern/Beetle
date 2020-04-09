// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/kubernetes"
	"github.com/clivern/beetle/internal/app/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Cluster controller
func Cluster(c *gin.Context) {
	cn := c.Param("cn")
	result := model.Cluster{}

	clusters, err := kubernetes.GetClusters()

	if err != nil {
		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Error(fmt.Sprintf(`Error: %s`, err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	var status bool

	for _, cluster := range clusters {
		if cn != cluster.Name {
			continue
		}

		status, err = cluster.Ping(context.Background())

		if err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
			}).Error(fmt.Sprintf(`Error: %s`, err.Error()))
		}

		result.Name = cluster.Name
		result.Health = status
	}

	if result.Name == "" {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":   result.Name,
		"health": result.Health,
	})
}
