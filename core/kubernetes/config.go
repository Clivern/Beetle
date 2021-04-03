// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"

	"github.com/clivern/beetle/core/model"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfig gets a beetle configs for a specific namespace
func (c *Cluster) GetConfig(ctx context.Context, namespace string) (model.Configs, error) {
	result := model.Configs{}

	err := c.Config()

	if err != nil {
		return result, err
	}

	data, err := c.ClientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(
			"%s=%s",
			"beetle.clivern.com/status",
			"enabled",
		),
	})

	if err != nil {
		return result, err
	}

	for _, deployment := range data.Items {
		applicationName := ""
		imageFormat := ""
		applicationID := ""
		status := "disabled"

		for key, value := range deployment.ObjectMeta.Annotations {
			if key == "beetle.clivern.com/application-name" {
				applicationName = value
			}
			if key == "beetle.clivern.com/image-format" {
				imageFormat = value
			}
		}
		for key, value := range deployment.ObjectMeta.Labels {
			if key == "beetle.clivern.com/status" {
				status = value
			}
			if key == "beetle.clivern.com/application-id" {
				applicationID = value
			}
		}

		if status == "enabled" && applicationID != "" && imageFormat != "" {
			result.Applications = append(result.Applications, model.App{
				ID:          applicationID,
				Name:        applicationName,
				ImageFormat: imageFormat,
			})
		} else {
			log.WithFields(log.Fields{
				"application_id": applicationID,
			}).Debug(`Application status disabled`)
		}
	}

	result.Exists = true

	return result, nil
}
