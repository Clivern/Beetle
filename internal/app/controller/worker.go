// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/internal/app/util"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Worker controller
func Worker(id int, messages <-chan string) {
	var ok bool
	var err error
	messageObj := model.Message{}

	logger, _ := module.NewLogger(
		viper.GetString("log.level"),
		viper.GetString("log.format"),
		[]string{viper.GetString("log.output")},
	)

	defer func() {
		_ = logger.Sync()
	}()

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
	}
}
