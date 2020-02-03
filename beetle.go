// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/clivern/beetle/internal/app/controller"
	"github.com/clivern/beetle/internal/app/middleware"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var configFile string
	var get string

	flag.StringVar(&configFile, "config", "config.prod.yml", "config")
	flag.StringVar(&get, "get", "", "get")
	flag.Parse()

	if get == "release" {
		fmt.Println(
			fmt.Sprintf(
				`Beetle Version %v Commit %v, Built @%v`,
				version,
				commit,
				date,
			),
		)
		return
	}

	if get == "health" {
		return
	}

	configUnparsed, err := ioutil.ReadFile(configFile)

	if err != nil {
		panic(fmt.Sprintf(
			"Error while reading config file [%s]: %s",
			configFile,
			err.Error(),
		))
	}

	configParsed, err := envsubst.EvalEnv(string(configUnparsed))

	if err != nil {
		panic(fmt.Sprintf(
			"Error while parsing config file [%s]: %s",
			configFile,
			err.Error(),
		))
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer([]byte(configParsed)))

	if err != nil {
		panic(fmt.Sprintf(
			"Error while loading configs [%s]: %s",
			configFile,
			err.Error(),
		))
	}

	if viper.GetString("log.output") != "stdout" {
		fs := module.FileSystem{}
		dir, _ := filepath.Split(viper.GetString("log.output"))

		if !fs.DirExists(dir) {
			if _, err := fs.EnsureDir(dir, 777); err != nil {
				panic(fmt.Sprintf(
					"Directory [%s] creation failed with error: %s",
					dir,
					err.Error(),
				))
			}
		}

		if !fs.FileExists(viper.GetString("log.output")) {
			f, err := os.Create(viper.GetString("log.output"))
			if err != nil {
				panic(fmt.Sprintf(
					"Error while creating log file [%s]: %s",
					viper.GetString("log.output"),
					err.Error(),
				))
			}
			defer f.Close()
		}
	}

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
	r.GET("/_metrics", controller.Metrics)

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
