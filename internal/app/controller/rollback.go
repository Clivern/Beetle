// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/internal/app/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// CreateRollback controller
func CreateRollback(c *gin.Context, messages chan<- string) {
	// Validate
	// ...
	// ...

	// Then create async job
	logger, _ := module.NewLogger(
		viper.GetString("log.level"),
		viper.GetString("log.format"),
		[]string{viper.GetString("log.output")},
	)

	defer func() {
		_ = logger.Sync()
	}()

	db := module.Database{}
	err := db.Connect(model.DSN{
		Driver:   viper.GetString("app.database.driver"),
		Username: viper.GetString("app.database.username"),
		Password: viper.GetString("app.database.password"),
		Hostname: viper.GetString("app.database.host"),
		Port:     viper.GetInt("app.database.port"),
		Name:     viper.GetString("app.database.name"),
	})

	if err != nil {
		logger.Error(fmt.Sprintf(
			`Error while connecting to database: %s`,
			err.Error(),
		), zap.String("CorrelationId", c.Request.Header.Get("X-Correlation-ID")))

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	uuid := util.GenerateUUID4()

	for db.GetJobByUUID(uuid).ID != 0 {
		uuid = util.GenerateUUID4()
	}

	job := db.CreateJob(&model.Job{
		UUID:   uuid,
		Status: model.JobPending,
		Type:   "rollback.create",
		RunAt:  time.Now(),
	})

	messageObj := model.Message{
		UUID: c.Request.Header.Get("X-Correlation-ID"),
		Job:  job.ID,
	}

	message, _ := messageObj.ConvertToJSON()

	messages <- message

	c.JSON(http.StatusAccepted, gin.H{
		"id":        job.UUID,
		"type":      job.Type,
		"status":    job.Status,
		"createdAt": job.CreatedAt,
	})
}

// GetRollback controller
func GetRollback(c *gin.Context) {
	c.Status(http.StatusOK)
}
