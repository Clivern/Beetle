// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"

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

	messageObj := model.Message{}

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

		err = job.Run()

		if err != nil {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Error(fmt.Sprintf(`Worker [%d] failure while executing async job [id=%d] [uuid=%s]: %s`, id, messageObj.Job, job.UUID, err.Error()))
		} else {
			log.WithFields(log.Fields{
				"CorrelationId": messageObj.UUID,
			}).Info(fmt.Sprintf(`Worker [%d] processed async job [id=%d] [uuid=%s]`, id, messageObj.Job, job.UUID))
		}

		db.UpdateJobByID(&job)
	}
}
