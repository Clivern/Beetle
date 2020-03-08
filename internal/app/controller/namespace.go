// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/kubernetes"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Namespace controller
func Namespace(c *gin.Context) {
	cn := c.Param("cn")
	ns := c.Param("ns")

	result := kubernetes.NamespaceModel{}

	logger, _ := module.NewLogger()

	defer logger.Sync()

	clusters, err := kubernetes.GetClusters()

	if err != nil {
		logger.Info(fmt.Sprintf(
			`Error! %s`,
			err.Error(),
		), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))

		c.Status(http.StatusInternalServerError)
		return
	}

	for _, cluster := range clusters {
		if cn != cluster.Name {
			continue
		}

		result, err = cluster.GetNamespace(ns)

		if err != nil {
			logger.Info(fmt.Sprintf(
				`Error! %s`,
				err.Error(),
			), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))
		}
	}

	if result.Name == "" {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":   result.Name,
		"uid":    result.UID,
		"status": result.Status,
	})
}
