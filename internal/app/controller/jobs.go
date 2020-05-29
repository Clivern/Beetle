// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Jobs controller
func Jobs(c *gin.Context) {
	db := module.Database{}

	err := db.AutoConnect()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		}).Error(fmt.Sprintf(`Error: %s`, err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	c.JSON(http.StatusOK, db.GetJobs())
}
