// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/clivern/beetle/internal/app/controller"
	"github.com/clivern/beetle/internal/app/middleware"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start beetle server",
	Run: func(cmd *cobra.Command, args []string) {
		configUnparsed, err := ioutil.ReadFile(config)

		if err != nil {
			panic(fmt.Sprintf(
				"Error while reading config file [%s]: %s",
				config,
				err.Error(),
			))
		}

		configParsed, err := envsubst.EvalEnv(string(configUnparsed))

		if err != nil {
			panic(fmt.Sprintf(
				"Error while parsing config file [%s]: %s",
				config,
				err.Error(),
			))
		}

		viper.SetConfigType("yaml")
		err = viper.ReadConfig(bytes.NewBuffer([]byte(configParsed)))

		if err != nil {
			panic(fmt.Sprintf(
				"Error while loading configs [%s]: %s",
				config,
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
			log.SetOutput(os.Stdout)
		} else {
			f, _ := os.Create(viper.GetString("log.output"))
			gin.DefaultWriter = io.MultiWriter(f)
			log.SetOutput(f)
		}

		lvl := strings.ToLower(viper.GetString("log.level"))
		level, err := log.ParseLevel(lvl)

		if err != nil {
			level = log.InfoLevel
		}

		log.SetLevel(level)

		if viper.GetString("app.mode") == "prod" {
			gin.SetMode(gin.ReleaseMode)
			gin.DefaultWriter = ioutil.Discard
			gin.DisableConsoleColor()
		}

		if viper.GetString("log.format") == "json" {
			log.SetFormatter(&log.JSONFormatter{})
		} else {
			log.SetFormatter(&log.TextFormatter{})
		}

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
		r.Use(middleware.Auth())
		r.Use(middleware.Logger())
		r.Use(middleware.Metric())

		r.GET("/favicon.ico", func(c *gin.Context) {
			c.String(http.StatusNoContent, "")
		})

		r.GET("/", controller.HealthCheck)
		r.GET("/_health", controller.HealthCheck)
		r.GET("/_ready", controller.ReadyCheck)
		r.GET(viper.GetString("app.metrics.prometheus.endpoint"), gin.WrapH(controller.Metrics()))
		r.GET("/api/v1/cluster", controller.Clusters)
		r.GET("/api/v1/cluster/:cn", controller.Cluster)
		r.GET("/api/v1/cluster/:cn/namespace", controller.Namespaces)
		r.GET("/api/v1/cluster/:cn/namespace/:ns", controller.Namespace)
		r.GET("/api/v1/cluster/:cn/namespace/:ns/app", controller.Applications)
		r.GET("/api/v1/cluster/:cn/namespace/:ns/app/:id", controller.Application)
		r.POST("/api/v1/cluster/:cn/namespace/:ns/app/:id/deployment", func(c *gin.Context) {
			controller.CreateDeployment(c, messages)
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

	},
}

func init() {
	serveCmd.Flags().StringVarP(&config, "config", "c", "config.prod.yml", "Absolute path to config file (required)")
	serveCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(serveCmd)
}
