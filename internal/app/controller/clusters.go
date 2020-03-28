// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/kubernetes"
	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Clusters controller
func Clusters(c *gin.Context) {
	result := []model.Cluster{}

	logger, _ := module.NewLogger()

	defer logger.Sync()

	clusters, err := kubernetes.GetClusters()

	if err != nil {
		logger.Error(fmt.Sprintf(
			`Error: %s`,
			err.Error(),
		), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))

		c.Status(http.StatusInternalServerError)
		return
	}

	var status bool

	for _, cluster := range clusters {
		status, err = cluster.Ping()

		if err != nil {
			logger.Error(fmt.Sprintf(
				`Error: %s`,
				err.Error(),
			), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))
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
