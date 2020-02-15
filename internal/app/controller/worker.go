// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Worker controller
func Worker(id int, messages <-chan string) {
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
	), zap.String("CorrelationId", ""))

	for message := range messages {
		logger.Info(fmt.Sprintf(
			`Worker [%d] Received: %s`,
			id,
			message,
		), zap.String("CorrelationId", ""))
	}
}
