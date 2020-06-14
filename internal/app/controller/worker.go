// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/clivern/beetle/internal/app/kubernetes"
	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/internal/app/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Worker controller
func Worker(workerID int, messages <-chan string) {
	var ok bool
	var err error
	var job model.Job
	var cluster *kubernetes.Cluster
	var uuid string

	messageObj := model.Message{}
	deploymentRequest := model.DeploymentRequest{}

	log.WithFields(log.Fields{
		"correlation_id": util.GenerateUUID4(),
		"worker_id":      workerID,
	}).Info(`Worker started`)

	db := module.Database{}

	for message := range messages {
		ok, err = messageObj.LoadFromJSON([]byte(message))

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"correlation_id": messageObj.UUID,
				"worker_id":      workerID,
				"message":        message,
			}).Warn(`Worker received invalid message`)
			continue
		}

		log.WithFields(log.Fields{
			"correlation_id": messageObj.UUID,
			"worker_id":      workerID,
			"job_id":         messageObj.Job,
		}).Info(`Worker received a new job`)

		err = db.AutoConnect()

		if err != nil {
			log.WithFields(log.Fields{
				"correlation_id": messageObj.UUID,
				"worker_id":      workerID,
				"error":          err.Error(),
			}).Error(`Worker unable to connect to database`)
			continue
		}

		job = db.GetJobByID(messageObj.Job)

		ok, err = deploymentRequest.LoadFromJSON([]byte(job.Payload))

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"correlation_id": messageObj.UUID,
				"worker_id":      workerID,
				"job_id":         messageObj.Job,
				"job_uuid":       job.UUID,
				"error":          err.Error(),
			}).Error(`Invalid job payload`)

			// Job Failed
			now := time.Now()
			job.Status = model.JobFailed
			job.RunAt = &now
			job.Result = fmt.Sprintf("Invalid job payload, UUID %s", messageObj.UUID)
			db.UpdateJobByID(&job)
			db.ReleaseChildJobs(job.ID)
			continue
		}

		log.WithFields(log.Fields{
			"correlation_id":      messageObj.UUID,
			"worker_id":           workerID,
			"job_id":              messageObj.Job,
			"job_uuid":            job.UUID,
			"request_cluster":     deploymentRequest.Cluster,
			"request_namespace":   deploymentRequest.Namespace,
			"request_application": deploymentRequest.Application,
			"request_version":     deploymentRequest.Version,
			"request_strategy":    deploymentRequest.Strategy,
		}).Info(`Worker accepted deployment request`)

		// Notify if there is a webhook
		if viper.GetString("app.webhook.url") != "" {
			uuid = util.GenerateUUID4()

			for db.JobExistByUUID(uuid) {
				uuid = util.GenerateUUID4()
			}

			db.CreateJob(&model.Job{
				UUID:    uuid,
				Payload: job.Payload,
				Status:  model.JobOnHold,
				Parent:  messageObj.Job,
				Type:    model.JobDeploymentNotify,
			})

			log.WithFields(log.Fields{
				"correlation_id":      messageObj.UUID,
				"worker_id":           workerID,
				"job_id":              messageObj.Job,
				"job_uuid":            job.UUID,
				"request_cluster":     deploymentRequest.Cluster,
				"request_namespace":   deploymentRequest.Namespace,
				"request_application": deploymentRequest.Application,
				"request_version":     deploymentRequest.Version,
				"request_strategy":    deploymentRequest.Strategy,
				"webhook_url":         viper.GetString("app.webhook.url"),
			}).Info(`HTTP webhook enabled`)
		} else {
			log.WithFields(log.Fields{
				"correlation_id":      messageObj.UUID,
				"worker_id":           workerID,
				"job_id":              messageObj.Job,
				"job_uuid":            job.UUID,
				"request_cluster":     deploymentRequest.Cluster,
				"request_namespace":   deploymentRequest.Namespace,
				"request_application": deploymentRequest.Application,
				"request_version":     deploymentRequest.Version,
				"request_strategy":    deploymentRequest.Strategy,
			}).Info(`HTTP webhook disabled`)
		}

		cluster, err = kubernetes.GetCluster(deploymentRequest.Cluster)

		if err != nil {
			log.WithFields(log.Fields{
				"correlation_id":      messageObj.UUID,
				"worker_id":           workerID,
				"error":               err.Error(),
				"request_cluster":     deploymentRequest.Cluster,
				"request_namespace":   deploymentRequest.Namespace,
				"request_application": deploymentRequest.Application,
				"request_version":     deploymentRequest.Version,
				"request_strategy":    deploymentRequest.Strategy,
			}).Error(`Worker can not find the cluster`)

			// Job Failed
			now := time.Now()
			job.Status = model.JobFailed
			job.RunAt = &now
			job.Result = fmt.Sprintf("Worker can not find the cluster, UUID %s", messageObj.UUID)
			db.UpdateJobByID(&job)
			db.ReleaseChildJobs(job.ID)
			continue
		}

		ok, err = cluster.Ping(context.TODO())

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"correlation_id":      messageObj.UUID,
				"worker_id":           workerID,
				"error":               err.Error(),
				"request_cluster":     deploymentRequest.Cluster,
				"request_namespace":   deploymentRequest.Namespace,
				"request_application": deploymentRequest.Application,
				"request_version":     deploymentRequest.Version,
				"request_strategy":    deploymentRequest.Strategy,
			}).Error(`Worker unable to ping cluster`)

			// Job Failed
			now := time.Now()
			job.Status = model.JobFailed
			job.RunAt = &now
			job.Result = fmt.Sprintf("Worker unable to ping cluster, UUID %s", messageObj.UUID)
			db.UpdateJobByID(&job)
			db.ReleaseChildJobs(job.ID)
			continue
		}

		ok, err = cluster.Deploy(deploymentRequest)

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"correlation_id":      messageObj.UUID,
				"worker_id":           workerID,
				"error":               err.Error(),
				"request_cluster":     deploymentRequest.Cluster,
				"request_namespace":   deploymentRequest.Namespace,
				"request_application": deploymentRequest.Application,
				"request_version":     deploymentRequest.Version,
				"request_strategy":    deploymentRequest.Strategy,
			}).Error(`Worker unable deploy`)

			// Job Failed
			now := time.Now()
			job.Status = model.JobFailed
			job.RunAt = &now
			job.Result = fmt.Sprintf("Failure during deployment, UUID %s", messageObj.UUID)
			db.UpdateJobByID(&job)
			db.ReleaseChildJobs(job.ID)
			continue
		}

		log.WithFields(log.Fields{
			"correlation_id":      messageObj.UUID,
			"worker_id":           workerID,
			"request_cluster":     deploymentRequest.Cluster,
			"request_namespace":   deploymentRequest.Namespace,
			"request_application": deploymentRequest.Application,
			"request_version":     deploymentRequest.Version,
			"request_strategy":    deploymentRequest.Strategy,
		}).Info(`Deployment finished successfully`)

		// Job Succeeded
		now := time.Now()
		job.Status = model.JobSuccess
		job.RunAt = &now
		job.Result = "Deployment finished successfully"
		db.UpdateJobByID(&job)
		db.ReleaseChildJobs(job.ID)
	}
}
