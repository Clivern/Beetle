// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"net/http"

	"github.com/clivern/beetle/core/kubernetes"
	"github.com/clivern/beetle/core/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Cluster controller
func Cluster(c *gin.Context) {
	cn := c.Param("cn")
	result := model.Cluster{}

	cluster, err := kubernetes.GetCluster(cn)

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"cluster_name":   cn,
			"error":          err.Error(),
		}).Info(`Cluster not found`)

		c.Status(http.StatusNotFound)
		return
	}

	status, err := cluster.Ping(context.TODO())

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"cluster_name":   cn,
			"error":          err.Error(),
		}).Error(`Error ping a cluster`)
	}

	result.Name = cluster.Name
	result.Health = status

	c.JSON(http.StatusOK, gin.H{
		"name":   result.Name,
		"health": result.Health,
	})
}
