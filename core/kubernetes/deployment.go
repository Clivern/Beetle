// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/clivern/beetle/core/model"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

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
func (c *Cluster) PatchDeployment(ctx context.Context, namespace, name, data string) (bool, error) {
	err := c.Config()

	if err != nil {
		return false, err
	}

	_, err = c.ClientSet.AppsV1().Deployments(namespace).Patch(
		ctx,
		name,
		types.JSONPatchType,
		[]byte(data),
		metav1.PatchOptions{},
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

// FetchDeploymentStatus get deployment status
func (c *Cluster) FetchDeploymentStatus(ctx context.Context, namespace, name string, limit int) (bool, error) {
	err := c.Config()

	if err != nil {
		return false, err
	}

	// Wait till k8s pick the deployment
	time.Sleep(10 * time.Second)

	deployment, err := c.ClientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	status := true

	for i := 0; i < limit; i++ {
		status = true

		deployment, err = c.ClientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})

		if err != nil {
			return false, err
		}

		if int(deployment.Generation) != int(deployment.Status.ObservedGeneration) {
			status = false
		}

		if int(deployment.Status.UnavailableReplicas) > 0 {
			status = false
		}

		if int(int32(*deployment.Spec.Replicas)) != int(deployment.Status.AvailableReplicas) {
			status = false
		}

		if !status {
			log.WithFields(log.Fields{
				"deployment.Generation":                 int(deployment.Generation),
				"deployment.Status.ObservedGeneration":  int(deployment.Status.ObservedGeneration),
				"deployment.Spec.Replicas":              int(int32(*deployment.Spec.Replicas)),
				"deployment.Status.AvailableReplicas":   int(deployment.Status.AvailableReplicas),
				"deployment.Status.UnavailableReplicas": int(deployment.Status.UnavailableReplicas),
			}).Debug(`Deployment Success`)
			time.Sleep(2 * time.Second)
		} else {
			log.WithFields(log.Fields{
				"deployment.Generation":                 int(deployment.Generation),
				"deployment.Status.ObservedGeneration":  int(deployment.Status.ObservedGeneration),
				"deployment.Spec.Replicas":              int(int32(*deployment.Spec.Replicas)),
				"deployment.Status.AvailableReplicas":   int(deployment.Status.AvailableReplicas),
				"deployment.Status.UnavailableReplicas": int(deployment.Status.UnavailableReplicas),
			}).Debug(`Deployment Success`)

			return true, nil
		}
	}

	log.WithFields(log.Fields{
		"deployment.Generation":                 int(deployment.Generation),
		"deployment.Status.ObservedGeneration":  int(deployment.Status.ObservedGeneration),
		"deployment.Spec.Replicas":              int(int32(*deployment.Spec.Replicas)),
		"deployment.Status.AvailableReplicas":   int(deployment.Status.AvailableReplicas),
		"deployment.Status.UnavailableReplicas": int(deployment.Status.UnavailableReplicas),
	}).Debug(`Deployment failure`)

	return false, fmt.Errorf(fmt.Sprintf(
		"Deployment %s failed: namespace %s, Generation %d, ObservedGeneration %d,"+
			" UnavailableReplicas %d, Replicas %d, AvailableReplicas %d",
		name,
		namespace,
		int(deployment.Generation),
		int(deployment.Status.ObservedGeneration),
		int(deployment.Status.UnavailableReplicas),
		int(int32(*deployment.Spec.Replicas)),
		int(deployment.Status.AvailableReplicas),
	))
}
