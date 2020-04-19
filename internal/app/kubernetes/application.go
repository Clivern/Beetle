// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"fmt"
	"context"

	"github.com/clivern/beetle/internal/app/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetApplicationVersion gets current application version
func (c *Cluster) GetApplicationVersion(ctx context.Context, namespace, id, format string) (string, error) {
	err := c.Config()

	if err != nil {
		return id, err
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
		return id, err
	}

	for _, deployment := range data.Items {
		fmt.Printf("%+v\n", deployment.Spec.Template.Spec.Containers)



		result = append(result, model.Deployment{
			Name: deployment.ObjectMeta.Name,
			UID: string(deployment.ObjectMeta.UID),
		})
	}

	return result, nil

	return id, nil
}
