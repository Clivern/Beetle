// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// App controller
func App(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// AppDeploy controller
func AppDeploy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// AppRollback controller
func AppRollback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
