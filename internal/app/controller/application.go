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

// Application controller
func Application(c *gin.Context) {
	cn := c.Param("cn")
	ns := c.Param("ns")
	id := c.Param("id")

	config := model.Configs{}

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

	config, err = cluster.GetConfig(context.TODO(), ns)

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"cluster_name":   cn,
			"namespace_name": ns,
			"error":          err.Error(),
		}).Warn(`Error while fetching beetle configMap`)
	}

	for _, app := range config.Applications {
		if app.ID == id {
			application, err := cluster.GetApplication(
				context.TODO(),
				ns,
				app.ID,
				app.Name,
				app.ImageFormat,
			)

			if err != nil {
				log.WithFields(log.Fields{
					"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
					"application_id": id,
					"cluster_name":   cn,
					"namespace_name": ns,
					"error":          err.Error(),
				}).Warn(`Error while fetching application current version`)
			}

			c.JSON(http.StatusOK, gin.H{
				"ID":         application.ID,
				"Name":       application.Name,
				"Format":     application.Format,
				"Containers": application.Containers,
			})
			return
		}
	}

	log.WithFields(log.Fields{
		"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		"application_id": id,
		"cluster_name":   cn,
		"namespace_name": ns,
	}).Info(`Application not found`)

	c.Status(http.StatusNotFound)
}
