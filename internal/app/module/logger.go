// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// NewLogger returns a logger instance
func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	rawJSON := []byte(fmt.Sprintf(`{
      		"level": "%s",
      		"encoding": "%s",
      		"outputPaths": []
    	}`, viper.GetString("log.level"), viper.GetString("log.format")))

	err := json.Unmarshal(rawJSON, &cfg)

	if err != nil {
		panic(err)
	}

	cfg.Encoding = viper.GetString("log.format")
	cfg.OutputPaths = []string{viper.GetString("log.output")}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}

	return logger, nil
}
