// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"net/http"

	"github.com/clivern/beetle/internal/app/kubernetes"
	"github.com/clivern/beetle/internal/app/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Clusters controller
func Clusters(c *gin.Context) {
	result := []model.Cluster{}

	clusters, err := kubernetes.GetClusters()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Error(`Error fetching clusters`)

		c.Status(http.StatusInternalServerError)
		return
	}

	var status bool

	for _, cluster := range clusters {
		status, err = cluster.Ping(context.TODO())

		if err != nil {
			log.WithFields(log.Fields{
				"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
				"cluster_name":   cluster.Name,
				"error":          err.Error(),
			}).Error(`Error while ping a cluster`)
		}

		result = append(result, model.Cluster{
			Name:   cluster.Name,
			Health: status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"clusters": result,
	})
}
