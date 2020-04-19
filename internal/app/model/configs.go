// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"gopkg.in/yaml.v3"
)

// Application struct
type Application struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	ImageFormat  string `yaml:"imageFormat"`
	CurrentImage string `yaml:"currentImage"`
	Deployment   []Deployment `yaml:"deployment"`
}

// Configs struct
type Configs struct {
	Exists       bool          `yaml:"exists"`
	Version      string        `yaml:"version"`
	Applications []Application `yaml:"applications"`
}

// LoadFromYAML update object from yaml
func (n *Configs) LoadFromYAML(data []byte) (bool, error) {
	err := yaml.Unmarshal(data, &n)

	if err != nil {
		return false, err
	}

	return true, nil
}

// ConvertToYAML convert object to yaml
func (n *Configs) ConvertToYAML() (string, error) {
	data, err := yaml.Marshal(&n)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
