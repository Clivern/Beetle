// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

// NewLogger returns a logger instance
func NewLogger(level, encoding string, outputPaths []string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	rawJSON := []byte(fmt.Sprintf(`{
      		"level": "%s",
      		"encoding": "%s",
      		"outputPaths": []
    	}`, level, encoding))

	err := json.Unmarshal(rawJSON, &cfg)

	if err != nil {
		panic(err)
	}

	cfg.Encoding = encoding
	cfg.OutputPaths = outputPaths

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}

	return logger, nil
}
