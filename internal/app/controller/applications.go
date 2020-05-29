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

// Applications controller
func Applications(c *gin.Context) {
	cn := c.Param("cn")
	ns := c.Param("ns")
	config := model.Configs{}

	cluster, err := kubernetes.GetCluster(cn)

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		}).Info(fmt.Sprintf(`Cluster not found %s: %s`, cn, err.Error()))

		c.Status(http.StatusNotFound)
		return
	}

	config, err = cluster.GetConfig(context.Background(), ns)

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		}).Warn(fmt.Sprintf(`Error while fetching beetle configMap for cluster %s namespace %s: %s`, cn, ns, err.Error()))
	}

	applications := []kubernetes.Application{}

	for _, app := range config.Applications {
		application, err := cluster.GetApplication(
			context.Background(),
			ns,
			app.ID,
			app.Name,
			app.ImageFormat,
		)

		if err != nil {
			log.WithFields(log.Fields{
				"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			}).Warn(fmt.Sprintf(`Error while fetching application %s current version cluster %s namespace %s: %s`, app.ID, cn, ns, err.Error()))
			continue
		}

		applications = append(applications, application)
	}

	c.JSON(http.StatusOK, gin.H{
		"apps": applications,
	})
}
