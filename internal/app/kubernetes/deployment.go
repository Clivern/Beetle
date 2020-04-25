// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"

	"github.com/clivern/beetle/internal/app/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

type patchUInt32Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

// GetDeployments gets a list of deployments
func (c *Cluster) GetDeployments(ctx context.Context, namespace, label string) ([]model.Deployment, error) {
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
			UID:  string(deployment.ObjectMeta.UID),
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
	result.UID = string(deployment.ObjectMeta.UID)

	return result, nil
}

// PatchDeployment updates the deployment
func (c *Cluster) PatchDeployment(ctx context.Context, namespace, name, data string) (model.Deployment, error) {
	result := model.Deployment{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	deployment, err := c.ClientSet.AppsV1().Deployments(namespace).Patch(
		ctx,
		name,
		types.JSONPatchType,
		[]byte(data),
		metav1.PatchOptions{},
	)

	if err != nil {
		return result, err
	}

	result.Name = deployment.ObjectMeta.Name
	result.UID = string(deployment.ObjectMeta.UID)

	return result, nil
}
