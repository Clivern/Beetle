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

// Application controller
func Application(c *gin.Context) {
	cn := c.Param("cn")
	ns := c.Param("ns")
	id := c.Param("id")

	config := model.Configs{}

	cluster, err := kubernetes.GetCluster(cn)

	if err != nil {
		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Info(fmt.Sprintf(`Cluster not found %s: %s`, cn, err.Error()))

		c.Status(http.StatusNotFound)
		return
	}

	config, err = cluster.GetConfig(context.Background(), ns)

	if err != nil {
		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Warn(fmt.Sprintf(`Error while fetching beetle configMap for cluster %s namespace %s: %s`, cn, ns, err.Error()))
	}

	for _, app := range config.Applications {
		if app.ID == id {
			application, err := cluster.GetApplication(
				context.Background(),
				ns,
				app.ID,
				app.Name,
				app.ImageFormat,
			)

			if err != nil {
				log.WithFields(log.Fields{
					"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
				}).Warn(fmt.Sprintf(`Error while fetching application %s current version cluster %s namespace %s: %s`, id, cn, ns, err.Error()))
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

	fmt.Println(config)

	log.WithFields(log.Fields{
		"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
	}).Info(fmt.Sprintf(`Application %s not found for cluster %s namespace %s`, id, cn, ns))

	c.Status(http.StatusNotFound)
}
