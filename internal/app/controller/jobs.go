// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Jobs controller
func Jobs(c *gin.Context) {
	logger, _ := module.NewLogger()

	defer logger.Sync()

	db := module.Database{}

	err := db.AutoConnect()

	if err != nil {
		logger.Error(fmt.Sprintf(
			`Error: %s`,
			err.Error(),
		), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	c.JSON(http.StatusOK, db.GetJobs())
}
