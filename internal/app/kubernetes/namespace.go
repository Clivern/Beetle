// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"strings"

	"github.com/clivern/beetle/internal/app/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNamespaces gets a list of cluster namespaces
func (c *Cluster) GetNamespaces(ctx context.Context) ([]model.Namespace, error) {
	result := []model.Namespace{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	data, err := c.ClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})

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

	err := c.Config()

	if err != nil {
		return result, err
	}

	namespace, err := c.ClientSet.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return result, err
	}

	result.Name = namespace.ObjectMeta.Name
	result.UID = string(namespace.ObjectMeta.UID)
	result.Status = strings.ToLower(string(namespace.Status.Phase))

	return result, nil
}
