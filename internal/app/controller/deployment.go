// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateDeployment controller
func CreateDeployment(c *gin.Context) {
	c.Status(http.StatusOK)
}

// GetDeployment controller
func GetDeployment(c *gin.Context) {
	c.Status(http.StatusOK)
}
