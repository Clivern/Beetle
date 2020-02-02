// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//+build wireinject

package main

import (
	"github.com/clivern/beetle/internal/app/cmd"
	"github.com/clivern/beetle/internal/app/config"
	"github.com/google/wire"
)

func InitializeNewAPI() *cmd.Agent {
	wire.Build(cmd.NewAPI, config.NewConfig)
	return &cmd.Agent{}
}
