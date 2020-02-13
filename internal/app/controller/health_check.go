// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// HealthCheck controller
func HealthCheck(c *gin.Context) {
	logger, _ := module.NewLogger(
		viper.GetString("log.level"),
		viper.GetString("log.format"),
		[]string{viper.GetString("log.output")},
	)

	defer func() {
		_ = logger.Sync()
	}()

	status := "ok"

	logger.Info(fmt.Sprintf(
		`Health Status: %s`,
		status,
	), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
