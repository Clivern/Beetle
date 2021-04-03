// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/clivern/beetle/core/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetApplication gets current application version
func (c *Cluster) GetApplication(ctx context.Context, namespace, id, name, format string) (model.Application, error) {
	result := model.Application{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	data, err := c.ClientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(
			"%s=%s,%s=%s",
			"beetle.clivern.com/status",
			"enabled",
			"beetle.clivern.com/application-id",
			id,
		),
	})

	if err != nil {
		return result, err
	}

	result.ID = id
	result.Name = name
	result.Format = format
	result.Containers = []model.Container{}

	for _, deployment := range data.Items {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			result.Containers = append(result.Containers, model.Container{
				Name:  container.Name,
				Image: container.Image,
				Version: strings.Replace(
					container.Image,
					strings.Replace(format, "[.Release]", "", -1),
					"",
					-1,
				),
				Deployment: model.Deployment{
					Name: deployment.ObjectMeta.Name,
					UID:  string(deployment.ObjectMeta.UID),
				},
			})
		}
	}

	return result, nil
}
