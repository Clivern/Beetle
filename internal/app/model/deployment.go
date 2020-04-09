// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
)

// Deployment struct
type Deployment struct {
	Name string `json:"name"`
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
