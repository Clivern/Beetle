// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
	"fmt"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/spf13/viper"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Clusters struct
type Clusters struct {
	Clusters []*Cluster `mapstructure:",clusters"`
}

// Cluster struct
type Cluster struct {
	Name          string `mapstructure:",name"`
	Kubeconfig    string `mapstructure:",kubeconfig"`
	ConfigMapName string `mapstructure:",configMapName"`
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

// GetCluster get a list of clusters
func GetCluster(name string) (*Cluster, error) {
	var clusters Clusters

	err := viper.UnmarshalKey("app", &clusters)

	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters.Clusters {
		if name == cluster.Name {
			return cluster, nil
		}
	}

	return &Cluster{}, fmt.Errorf("Unable to find cluster %s", name)
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
