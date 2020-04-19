// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"

	"github.com/clivern/beetle/internal/app/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDeployments gets a list of deployments
func (c *Cluster) GetDeployments(ctx context.Context, namespace string, label string) ([]model.Deployment, error) {
	result := []model.Deployment{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	data, err := c.ClientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: label,
	})

	if err != nil {
		return result, err
	}

	for _, deployment := range data.Items {
		result = append(result, model.Deployment{
			Name: deployment.ObjectMeta.Name,
		})
	}

	return result, nil
}

// GetDeployment gets a deployment by name
func (c *Cluster) GetDeployment(ctx context.Context, namespace, name string) (model.Deployment, error) {
	result := model.Deployment{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	deployment, err := c.ClientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return result, err
	}

	result.Name = deployment.ObjectMeta.Name

	return result, nil
}
