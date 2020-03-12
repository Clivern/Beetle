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
	"go.uber.org/zap"
)

// CreateRollback controller
func CreateRollback(c *gin.Context, messages chan<- string) {
	// Validate
	// ...
	// ...

	// Then create async job
	logger, _ := module.NewLogger()

	defer logger.Sync()

	db := module.Database{}
	err := db.AutoConnect()

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

	for db.JobExistByUUID(uuid) {
		uuid = util.GenerateUUID4()
	}

	job := db.CreateJob(&model.Job{
		UUID:   uuid,
		Status: model.JobPending,
		Type:   model.JobRollbackCreate,
		RunAt:  time.Now(),
	})

	messageObj := model.Message{
		UUID: c.Request.Header.Get("X-Correlation-ID"),
		Job:  job.ID,
	}

	message, _ := messageObj.ConvertToJSON()

	// Send the job to workers
	messages <- message

	c.JSON(http.StatusAccepted, gin.H{
		"id":        job.UUID,
		"type":      job.Type,
		"status":    job.Status,
		"createdAt": job.CreatedAt,
	})
}
