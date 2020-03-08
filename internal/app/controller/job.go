// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/gin-gonic/gin"
)

// Job controller
func Job(c *gin.Context) {
	id := c.Param("id")

	// Init DB Connection
	db := module.Database{}
	err := db.AutoConnect()

	if err != nil {
		panic(err.Error())
	}

	// Migrate Database
	success := db.Migrate()

	if !success {
		panic("Error! Unable to migrate database tables.")
	}

	defer db.Close()

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
