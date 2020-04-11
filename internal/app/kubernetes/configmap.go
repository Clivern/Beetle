// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetConfigMap gets a configmap data
func (c *Cluster) GetConfigMap(ctx context.Context, namespace, name string) (model.ConfigMap, error) {
	result := model.ConfigMap{}

	fs := module.FileSystem{}

	if !fs.FileExists(c.Kubeconfig) {
		return result, fmt.Errorf(
			"cluster [%s] config file [%s] not exist",
			c.Name,
			c.Kubeconfig,
		)
	}

	config, err := clientcmd.BuildConfigFromFlags("", c.Kubeconfig)

	if err != nil {
		return result, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return result, err
	}

	configmap, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})

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
