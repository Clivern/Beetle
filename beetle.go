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

	"github.com/clivern/beetle/internal/app/cmd"
	"github.com/clivern/beetle/internal/app/controller"
	"github.com/clivern/beetle/internal/app/middleware"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	flag.StringVar(&configFile, "config", "config.prod.yml", "Absolute path to config file")
	flag.StringVar(&get, "get", "", "Use with version or release (./beetle --get release) (./beetle --get version) (./beetle --get health)")
	flag.Parse()

	if get == "release" || get == "version" {
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

	if get == "health" {
		cmd.HealthCheck()
		return
	}

	if viper.GetString("log.output") == "stdout" {
		gin.DefaultWriter = os.Stdout
		log.SetOutput(os.Stdout)
	} else {
		f, _ := os.Create(viper.GetString("log.output"))
		gin.DefaultWriter = io.MultiWriter(f)
	}

	if viper.GetString("log.level") == "info" {
		log.SetLevel(log.InfoLevel)
	}

	if viper.GetString("app.mode") == "prod" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
		gin.DisableConsoleColor()
	}

	log.SetFormatter(&log.JSONFormatter{})

	// Init DB Connection
	db := module.Database{}
	err = db.AutoConnect()

	if err != nil {
		panic(err.Error())
	}

	// Migrate Database
	success := db.Migrate()

	if !success {
		panic("Error! Unable to migrate database tables.")
	}

	defer db.Close()

	messages := make(chan string, viper.GetInt("app.broker.native.capacity"))

	r := gin.Default()

	r.Use(middleware.Correlation())
	r.Use(middleware.Logger())
	r.Use(middleware.Metric())

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	r.GET("/", controller.HealthCheck)
	r.GET("/_health", controller.HealthCheck)
	r.GET("/_metrics", gin.WrapH(controller.Metrics()))
	r.GET("/api/v1/cluster", controller.Clusters)
	r.GET("/api/v1/cluster/:cn", controller.Cluster)
	r.GET("/api/v1/cluster/:cn/namespace", controller.Namespaces)
	r.GET("/api/v1/cluster/:cn/namespace/:ns", controller.Namespace)
	r.GET("/api/v1/cluster/:cn/namespace/:ns/app", controller.Applications)
	r.GET("/api/v1/cluster/:cn/namespace/:ns/app/:id", controller.Application)
	r.POST("/api/v1/cluster/:cn/namespace/:ns/app/:id/deployment", func(c *gin.Context) {
		controller.CreateDeployment(c, messages)
	})
	r.POST("/api/v1/cluster/:cn/namespace/:ns/app/:id/rollback", func(c *gin.Context) {
		controller.CreateRollback(c, messages)
	})
	r.GET("/api/v1/job", controller.Jobs)
	r.GET("/api/v1/job/:uuid", controller.GetJob)
	r.DELETE("/api/v1/job/:uuid", controller.DeleteJob)

	for i := 0; i < viper.GetInt("app.broker.native.workers"); i++ {
		go controller.Worker(i+1, messages)
	}

	var runerr error

	if viper.GetBool("app.tls.status") {
		runerr = r.RunTLS(
			fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
			viper.GetString("app.tls.pemPath"),
			viper.GetString("app.tls.keyPath"),
		)
	} else {
		runerr = r.Run(
			fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
		)
	}

	if runerr != nil {
		panic(runerr.Error())
	}
}
