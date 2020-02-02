// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/clivern/beetle/internal/app/config"
	"github.com/clivern/beetle/internal/app/controller"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// API struct
type API struct {
	Config *config.Config
}

// NewAPI create a new instance
func NewAPI(config *config.Config) *API {
	return &API{
		Config: config,
	}
}

// Run runs the api
func (w *API) Run() {
	fmt.Println("API server started .....")

	r := gin.Default()

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	r.GET("/api/_health", controller.HealthCheck)

	if viper.GetBool("api.tls.status") {
		r.RunTLS(
			fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("api.port"))),
			viper.GetString("api.tls.pemPath"),
			viper.GetString("api.tls.keyPath"),
		)
	} else {
		r.Run(
			fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("api.port"))),
		)
	}
}
