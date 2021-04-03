// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

// App struct
type App struct {
	ID          string
	Name        string
	ImageFormat string
}

// Configs struct
type Configs struct {
	Exists       bool
	Version      string
	Applications []App
}
