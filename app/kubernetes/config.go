// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"

	"github.com/clivern/beetle/app/model"
)

// GetConfig gets a beetle configs for a specific namespace
func (c *Cluster) GetConfig(ctx context.Context, namespace string) (model.Configs, error) {
	result := model.Configs{}

	// Get ConfigMap
	item, err := c.GetConfigMap(ctx, namespace, c.ConfigMapName)

	if err != nil {
		return result, err
	}

	// If configs not exist, fallback to defaults
	if _, ok := item.Data["config"]; !ok {
		result.Exists = false
		return result, nil
	}

	// Convert to struct
	ok, err := result.LoadFromYAML([]byte(item.Data["config"]))

	if !ok || err != nil {
		return result, err
	}

	result.Exists = true

	return result, nil
}
