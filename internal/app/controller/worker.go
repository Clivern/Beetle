// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/internal/app/util"

	"go.uber.org/zap"
)

// Worker controller
func Worker(id int, messages <-chan string) {
	var ok bool
	var err error
	var db module.Database
	var job model.Job

	messageObj := model.Message{}

	logger, _ := module.NewLogger()

	defer logger.Sync()

	logger.Info(fmt.Sprintf(
		`Worker [%d] started`,
		id,
	), zap.String("CorrelationId", util.GenerateUUID4()))

	for message := range messages {
		ok, err = messageObj.LoadFromJSON([]byte(message))

		if !ok || err != nil {
			logger.Warn(fmt.Sprintf(
				`Worker [%d] received invalid message: %s`,
				id,
				message,
			), zap.String("CorrelationId", messageObj.UUID))
			continue
		}

		logger.Info(fmt.Sprintf(
			`Worker [%d] received: %d`,
			id,
			messageObj.Job,
		), zap.String("CorrelationId", messageObj.UUID))

		db = module.Database{}

		err = db.AutoConnect()

		if err != nil {
			logger.Error(fmt.Sprintf(
				`Worker [%d] unable to connect to database: %s`,
				id,
				err.Error(),
			), zap.String("CorrelationId", messageObj.UUID))
			continue
		}

		defer db.Close()

		job = db.GetJobByID(messageObj.Job)

		err = job.Run()

		if err != nil {
			logger.Error(fmt.Sprintf(
				`Worker [%d] failure while executing async job [%d] [%s]: %s`,
				id,
				messageObj.Job,
				job.UUID,
				err.Error(),
			), zap.String("CorrelationId", messageObj.UUID))
		}

		logger.Info(fmt.Sprintf(
			`Worker [%d] processed async job [%d] [%s]`,
			id,
			messageObj.Job,
			job.UUID,
		), zap.String("CorrelationId", messageObj.UUID))

		db.UpdateJobByID(&job)
	}
}
