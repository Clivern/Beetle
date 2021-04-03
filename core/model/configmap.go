// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
)

// ConfigMap struct
type ConfigMap struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	CreationTimestamp string            `json:"creation_timestamp"`
	Data              map[string]string `json:"data"`
	Labels            map[string]string `json:"labels"`
}

// LoadFromJSON update object from json
func (d *ConfigMap) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &d)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (d *ConfigMap) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
