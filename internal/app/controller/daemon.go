// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	pendingJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "beetle",
			Name:      "workers_queue_pending_jobs",
			Help:      "The pending jobs in the queue",
		})

	failedJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "beetle",
			Name:      "workers_queue_failed_jobs",
			Help:      "The failed jobs in the queue",
		})

	successJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "beetle",
			Name:      "workers_queue_success_jobs",
			Help:      "The successful jobs in the queue",
		})

	onHoldJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "beetle",
			Name:      "workers_queue_on_hold_jobs",
			Help:      "The on hold jobs in the queue",
		})
)

func init() {
	prometheus.MustRegister(pendingJobs)
	prometheus.MustRegister(failedJobs)
	prometheus.MustRegister(successJobs)
	prometheus.MustRegister(onHoldJobs)
}

// Daemon function
func Daemon() {
	var err error
	var pendingJobsCount int
	var failedJobsCount int
	var successfulJobsCount int
	var onHoldJobsCount int
	var job model.Job
	var parentJob model.Job
	var deploymentRequest model.DeploymentRequest
	var payload string

	httpClient := module.NewHTTPClient()
	db := module.Database{}

	retry, err := strconv.Atoi(viper.GetString("app.webhook.retry"))

	if err != nil {
		panic(err.Error())
	}

	for {
		err = db.AutoConnect()

		if err != nil {
			log.WithFields(log.Fields{
				"correlation_id": "",
				"error":          err.Error(),
			}).Error(`Failure while connecting database`)
			time.Sleep(2 * time.Second)
			continue
		}
		// Update Metrics
		pendingJobsCount = db.CountJobs(model.JobPending)
		failedJobsCount = db.CountJobs(model.JobFailed)
		successfulJobsCount = db.CountJobs(model.JobSuccess)
		onHoldJobsCount = db.CountJobs(model.JobOnHold)

		log.WithFields(log.Fields{
			"correlation_id":        "",
			"pending_jobs_count":    pendingJobsCount,
			"failed_jobs_count":     failedJobsCount,
			"successful_jobs_count": successfulJobsCount,
			"on_hold_jobs_count":    onHoldJobsCount,
		}).Debug(`Update metrics`)

		pendingJobs.Set(float64(pendingJobsCount))
		failedJobs.Set(float64(failedJobsCount))
		successJobs.Set(float64(successfulJobsCount))
		onHoldJobs.Set(float64(onHoldJobsCount))

		// Run Pending Jobs (HTTP Notification)
		job = db.GetPendingJobByType(model.JobDeploymentNotify)

		if job.ID > 0 {
			if job.Retry > retry {
				now := time.Now()
				job.Status = model.JobFailed
				job.RunAt = &now
				job.Result = fmt.Sprintf("Failed to deliver the notification")
				db.UpdateJobByID(&job)
			} else {
				deploymentRequest.LoadFromJSON([]byte(job.Payload))

				if job.Parent > 0 {
					parentJob = db.GetJobByID(job.Parent)

					if parentJob.ID > 0 {
						deploymentRequest.Status = parentJob.Status
					}
				}

				payload, _ = deploymentRequest.ConvertToJSON()

				response, err := httpClient.Post(
					context.TODO(),
					viper.GetString("app.webhook.url"),
					payload,
					map[string]string{},
					map[string]string{
						"Content-Type":      "application/json",
						"X-AUTH-TOKEN":      viper.GetString("app.webhook.token"),
						"X-NOTIFICATION-ID": job.UUID,
						"X-ACTION-NAME":     job.Type,
						"X-DEPLOYMENT-ID":   parentJob.UUID,
					},
				)

				if httpClient.GetStatusCode(response) != http.StatusOK || err != nil {
					job.Status = model.JobFailed
					job.Result = fmt.Sprintf("Failed to deliver the notification")
				} else {
					job.Status = model.JobSuccess
					job.Result = fmt.Sprintf("Notification delivered successfully")
				}

				if job.Status == model.JobFailed && job.Retry <= retry {
					job.Status = model.JobPending
				}

				now := time.Now()
				job.Retry++
				job.RunAt = &now
				db.UpdateJobByID(&job)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
