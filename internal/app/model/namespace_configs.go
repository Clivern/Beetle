// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"gopkg.in/yaml.v3"
)

// NamespaceConfigs struct
type NamespaceConfigs struct {
	Version string `yaml:"version"`
}

// LoadFromYAML update object from yaml
func (n *NamespaceConfigs) LoadFromYAML(data []byte) (bool, error) {
	err := yaml.Unmarshal(data, &n)

	if err != nil {
		return false, err
	}

	return true, nil
}

// ConvertToYAML convert object to yaml
func (n *NamespaceConfigs) ConvertToYAML() (string, error) {
	data, err := yaml.Marshal(&n)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
