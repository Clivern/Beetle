// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/clivern/beetle/internal/app/config"
	"github.com/clivern/beetle/internal/app/controller"
	"github.com/clivern/beetle/internal/app/middleware"

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

// Run runs the app
func (w *API) Run() {

	if viper.GetString("log.output") == "stdout" {
		gin.DefaultWriter = os.Stdout
	} else {
		f, _ := os.Create(viper.GetString("log.output"))
		gin.DefaultWriter = io.MultiWriter(f)
	}

	if viper.GetString("app.mode") == "prod" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
		gin.DisableConsoleColor()
	}

	r := gin.Default()

	r.Use(middleware.Correlation())
	r.Use(middleware.Logger())

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	r.GET("/_health", controller.HealthCheck)

	if viper.GetBool("app.tls.status") {
		r.RunTLS(
			fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
			viper.GetString("app.tls.pemPath"),
			viper.GetString("app.tls.keyPath"),
		)
	} else {
		r.Run(
			fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
		)
	}
}
