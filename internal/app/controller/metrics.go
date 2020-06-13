// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"
	"time"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	workersCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "beetle",
			Name:      "workers_count",
			Help:      "Number of Async Workers",
		})

	queueCapacity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "beetle",
			Name:      "workers_queue_capacity",
			Help:      "The maximum number of messages queue can process",
		})

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
	prometheus.MustRegister(workersCount)
	prometheus.MustRegister(queueCapacity)
	prometheus.MustRegister(pendingJobs)
	prometheus.MustRegister(failedJobs)
	prometheus.MustRegister(successJobs)
}

// Metrics controller
func Metrics() http.Handler {
	var err error
	var pendingJobsCount int
	var failedJobsCount int
	var successfulJobsCount int

	workersCount.Set(float64(viper.GetInt("app.broker.native.workers")))
	queueCapacity.Set(float64(viper.GetInt("app.broker.native.capacity")))

	db := module.Database{}

	// spin a goroutine to update db metrics
	go func() {
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
			time.Sleep(2 * time.Second)
		}
	}()

	return promhttp.Handler()
}
