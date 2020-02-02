// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package config

// Log struct
type Log struct {
	Level  string `mapstructure:"level"`
	Output string `mapstructure:"output"`
	Format string `mapstructure:"format"`
}

// Config struct
type Config struct {
	Log *Log `mapstructure:"log"`
}

// NewConfig create a new instance
func NewConfig() *Config {
	return &Config{}
}
