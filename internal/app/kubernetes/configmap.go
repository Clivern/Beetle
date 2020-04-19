// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"

	"github.com/clivern/beetle/internal/app/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfigMap gets a configmap data
func (c *Cluster) GetConfigMap(ctx context.Context, namespace, name string) (model.ConfigMap, error) {
	result := model.ConfigMap{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	configmap, err := c.ClientSet.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return result, err
	}

	result.Name = configmap.ObjectMeta.Name
	result.Namespace = configmap.ObjectMeta.Namespace
	result.UID = string(configmap.ObjectMeta.UID)
	result.CreationTimestamp = configmap.ObjectMeta.CreationTimestamp.String()
	result.Data = configmap.Data
	result.Labels = configmap.ObjectMeta.Labels

	return result, nil
}
