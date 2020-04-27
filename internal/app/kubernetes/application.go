// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Container struct
type Container struct {
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	Version    string     `json:"version"`
	Deployment Deployment `json:"deployment"`
}

// Deployment struct
type Deployment struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// Application struct
type Application struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Format     string      `json:"format"`
	Containers []Container `json:"containers"`
}

// GetApplication gets current application version
func (c *Cluster) GetApplication(ctx context.Context, namespace, id, name, format string) (Application, error) {
	result := Application{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	data, err := c.ClientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(
			"%s=%s,%s=%s",
			"app.clivern.com/managed-by",
			"beetle",
			"app.clivern.com/application-id",
			id,
		),
	})

	if err != nil {
		return result, err
	}

	result.ID = id
	result.Name = name
	result.Format = format
	result.Containers = []Container{}

	for _, deployment := range data.Items {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			result.Containers = append(result.Containers, Container{
				Name:  container.Name,
				Image: container.Image,
				Version: strings.Replace(
					container.Image,
					strings.Replace(format, "[.Release]", "", -1),
					"",
					-1,
				),
				Deployment: Deployment{
					Name: deployment.ObjectMeta.Name,
					UID:  string(deployment.ObjectMeta.UID),
				},
			})
		}
	}

	return result, nil
}
