// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/internal/app/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CreateDeployment controller
func CreateDeployment(c *gin.Context, messages chan<- string) {
	rawBody, _ := c.GetRawData()

	deploymentRequest := model.DeploymentRequest{}

	_, err := deploymentRequest.LoadFromJSON(rawBody)

	deploymentRequest.Cluster = c.Param("cn")
	deploymentRequest.Namespace = c.Param("ns")
	deploymentRequest.Application = c.Param("id")

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Info(`Invalid request`)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request!",
		})
		return
	}

	err = deploymentRequest.Validate([]string{
		model.RecreateStrategy,
		model.RampedStrategy,
		model.CanaryStrategy,
		model.BlueGreenStrategy,
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Info(`Invalid request`)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := deploymentRequest.ConvertToJSON()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Info(`Invalid request`)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Then create async job
	db := module.Database{}
	err = db.AutoConnect()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Error(`Failure while connecting database`)

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	uuid := util.GenerateUUID4()

	for db.JobExistByUUID(uuid) {
		uuid = util.GenerateUUID4()
	}

	job := db.CreateJob(&model.Job{
		UUID:    uuid,
		Payload: result,
		Status:  model.JobPending,
		Type:    model.JobDeploymentCreate,
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
