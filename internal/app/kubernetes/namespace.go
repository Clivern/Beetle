// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetNamespaces gets a list of cluster namespaces
func (c *Cluster) GetNamespaces(ctx context.Context) ([]model.Namespace, error) {
	result := []model.Namespace{}

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

	data, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})

	if err != nil {
		return result, err
	}

	for _, namespace := range data.Items {
		result = append(result, model.Namespace{
			Name:   namespace.ObjectMeta.Name,
			UID:    string(namespace.ObjectMeta.UID),
			Status: strings.ToLower(string(namespace.Status.Phase)),
		})
	}

	return result, nil
}

// GetNamespace gets a namespace by name
func (c *Cluster) GetNamespace(ctx context.Context, name string) (model.Namespace, error) {
	result := model.Namespace{}

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

	namespace, err := clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return result, err
	}

	result.Name = namespace.ObjectMeta.Name
	result.UID = string(namespace.ObjectMeta.UID)
	result.Status = strings.ToLower(string(namespace.Status.Phase))

	return result, nil
}
