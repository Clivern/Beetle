// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
)

// Container struct
type Container struct {
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	Version    string     `json:"version"`
	Deployment Deployment `json:"deployment"`
}

// Application struct
type Application struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Format     string      `json:"format"`
	Containers []Container `json:"containers"`
}

// Deployment struct
type Deployment struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// LoadFromJSON update object from json
func (c *Application) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (c *Application) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (d *Deployment) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &d)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (d *Deployment) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
