// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"time"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
)

func init() {
	prometheus.MustRegister(pendingJobs)
	prometheus.MustRegister(failedJobs)
	prometheus.MustRegister(successJobs)
}

// Daemon function
func Daemon() {
	var err error
	var pendingJobsCount int
	var failedJobsCount int
	var successfulJobsCount int

	db := module.Database{}

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

		log.WithFields(log.Fields{
			"correlation_id":        "",
			"pending_jobs_count":    pendingJobsCount,
			"failed_jobs_count":     failedJobsCount,
			"successful_jobs_count": successfulJobsCount,
		}).Debug(`Update metrics`)

		pendingJobs.Set(float64(pendingJobsCount))
		failedJobs.Set(float64(failedJobsCount))
		successJobs.Set(float64(successfulJobsCount))

		// Run Pending Jobs (HTTP Notification)
		time.Sleep(2 * time.Second)
	}
}
