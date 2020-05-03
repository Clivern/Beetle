// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/json"
)

type clusterResponse struct {
	Name   string `json:"name"`
	Health bool   `json:"health"`
}

type clustersResponse struct {
	Clusters []clusterResponse `json:"clusters"`
}

// LoadFromJSON update object from json
func (d *clusterResponse) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &d)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (d *clusterResponse) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (d *clustersResponse) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &d)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (d *clustersResponse) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
