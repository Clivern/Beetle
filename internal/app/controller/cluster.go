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
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Cluster controller
func Cluster(c *gin.Context) {
	result := []kubernetes.ClusterModel{}

	logger, _ := module.NewLogger(
		viper.GetString("log.level"),
		viper.GetString("log.format"),
		[]string{viper.GetString("log.output")},
	)

	defer func() {
		_ = logger.Sync()
	}()

	clusters, err := kubernetes.GetClusters()

	if err != nil {
		logger.Info(fmt.Sprintf(
			`Error! %s`,
			err.Error(),
		), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))

		c.JSON(http.StatusOK, gin.H{
			"clusters": result,
		})

		return
	}

	var status bool

	for _, cluster := range clusters {
		status, err = cluster.Ping()

		if err != nil {
			logger.Info(fmt.Sprintf(
				`Error! %s`,
				err.Error(),
			), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))
		}

		result = append(result, kubernetes.ClusterModel{
			Name:   cluster.Name,
			Health: status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"clusters": result,
	})
}
