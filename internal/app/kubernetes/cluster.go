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

	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Clusters struct
type Clusters struct {
	Clusters []*Cluster `mapstructure:",clusters"`
}

// Cluster struct
type Cluster struct {
	Name       string `mapstructure:",name"`
	Kubeconfig string `mapstructure:",kubeconfig"`
}

// GetClusters get a list of clusters
func GetClusters() ([]*Cluster, error) {
	var clusters Clusters

	err := viper.UnmarshalKey("app", &clusters)

	if err != nil {
		return nil, err
	}

	return clusters.Clusters, nil
}

// Ping check the cluster
func (c *Cluster) Ping(ctx context.Context) (bool, error) {
	fs := module.FileSystem{}

	if !fs.FileExists(c.Kubeconfig) {
		return false, fmt.Errorf(
			"cluster [%s] config file [%s] not exist",
			c.Name,
			c.Kubeconfig,
		)
	}

	config, err := clientcmd.BuildConfigFromFlags("", c.Kubeconfig)

	if err != nil {
		return false, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return false, err
	}

	data, err := clientset.RESTClient().Get().AbsPath("/api/v1").DoRaw(ctx)

	if err != nil {
		return false, err
	}

	return (string(data) != ""), nil
}

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

// GetDeployments gets a list of deployments
func (c *Cluster) GetDeployments(ctx context.Context, namespace string, labels []string) ([]model.Deployment, error) {
	result := []model.Deployment{}

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

	data, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return result, err
	}

	for _ = range data.Items {
		//fmt.Printf("%v", deployment)
		result = append(result, model.Deployment{})
	}

	return result, nil
}

// GetDeployment gets a deployment by name
func (c *Cluster) GetDeployment(ctx context.Context, namespace, name string) (model.Deployment, error) {
	result := model.Deployment{}

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

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return result, err
	}

	return result, nil
}
