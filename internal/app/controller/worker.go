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
)

// Worker controller
func Worker(id int, messages <-chan string) {
	var ok bool
	var err error
	var db module.Database
	var job model.Job
	var cluster *kubernetes.Cluster

	messageObj := model.Message{}
	deploymentRequest := model.DeploymentRequest{}

	log.WithFields(log.Fields{
		"CorrelationId": util.GenerateUUID4(),
	}).Info(fmt.Sprintf(`Worker [%d] started`, id))

	for message := range messages {
		ok, err = messageObj.LoadFromJSON([]byte(message))

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Warn(fmt.Sprintf(`Worker [%d] received invalid message: %s`, id, message))
			continue
		}

		log.WithFields(log.Fields{
			"CorrelationId": messageObj.UUID,
		}).Info(fmt.Sprintf(`Worker [%d] received: %d`, id, messageObj.Job))

		db = module.Database{}

		err = db.AutoConnect()

		if err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Error(fmt.Sprintf(`Worker [%d] unable to connect to database: %s`, id, err.Error()))
			continue
		}

		defer db.Close()

		job = db.GetJobByID(messageObj.Job)

		job.Status = model.JobPending

		ok, err = deploymentRequest.LoadFromJSON([]byte(job.Payload))

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Error(fmt.Sprintf(`Worker [%d] failure while executing async job [id=%d] [uuid=%s]: %s`, id, messageObj.Job, job.UUID, err.Error()))
			continue
		} else {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Info(fmt.Sprintf(`Worker [%d] processed async job [id=%d] [uuid=%s]`, id, messageObj.Job, job.UUID))
		}

		cluster, err = kubernetes.GetCluster(deploymentRequest.Cluster)

		if err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Info(fmt.Sprintf(`Worker [%d] Cluster not found %s: %s`, id, deploymentRequest.Cluster, err.Error()))
			continue
		}

		ok, err = cluster.Ping(context.Background())

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Error(fmt.Sprintf(`Worker [%d] Unable to connect to cluster %s error: %s`, id, deploymentRequest.Cluster, err.Error()))
		}

		ok, err = cluster.Deploy(deploymentRequest)

		if !ok || err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Error(fmt.Sprintf(`Worker [%d] Unable to connect to cluster %s error: %s`, id, deploymentRequest.Cluster, err.Error()))
			continue
		}

		// Wait for the deployment to check the final
		// Status https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#deployment-status

		now := time.Now()

		job.RunAt = &now

		db.UpdateJobByID(&job)
	}
}
