// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"
)

// Job struct
type Job struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Result    string    `json:"result"`
	RunAt     time.Time `json:"run_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoadFromJSON update object from json
func (j *Job) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &j)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (j *Job) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&j)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
