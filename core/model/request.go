// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	// RecreateStrategy var
	RecreateStrategy = "recreate"
	// RampedStrategy var
	RampedStrategy = "ramped"
	// CanaryStrategy var
	CanaryStrategy = "canary"
	// BlueGreenStrategy var
	BlueGreenStrategy = "blue_green"
)

// DeploymentRequest struct
type DeploymentRequest struct {
	Cluster     string `json:"cluster"`
	Namespace   string `json:"namespace"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Strategy    string `json:"strategy"`
	Status      string `json:"status"`

	// Ramped Strategy
	MaxSurge       string `json:"maxSurge"`
	MaxUnavailable string `json:"maxUnavailable"`
}

// LoadFromJSON update object from json
func (d *DeploymentRequest) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &d)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (d *DeploymentRequest) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Validate validates the request
func (d *DeploymentRequest) Validate(strategies []string) error {
	if d.Version == "" {
		return fmt.Errorf(
			"Error! version is required",
		)
	}

	if !In(d.Strategy, strategies) {
		return fmt.Errorf(
			"Error! strategy %s is invalid",
			d.Strategy,
		)
	}

	return nil
}

// In check if value is on array
func In(val interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}
