// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/clivern/beetle/internal/app/model"

	"github.com/gin-gonic/gin"
)

// CreateDeployment controller
func CreateDeployment(c *gin.Context, messages chan<- string) {
	messageObj := model.Message{
		UUID: c.Request.Header.Get("X-Correlation-ID"),
	}

	message, _ := messageObj.ConvertToJSON()

	messages <- message
	c.Status(http.StatusOK)
}

// GetDeployment controller
func GetDeployment(c *gin.Context) {
	c.Status(http.StatusOK)
}
