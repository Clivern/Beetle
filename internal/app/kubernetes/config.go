// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
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
// clus, _ := kubernetes.GetClusters()
// fmt.Println(clus[0].Ping())
func (c *Cluster) Ping() (bool, error) {
	config, err := clientcmd.BuildConfigFromFlags("", c.Kubeconfig)

	if err != nil {
		return false, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return false, err
	}

	data, err := clientset.RESTClient().Get().AbsPath("/api/v1").DoRaw()

	if err != nil {
		return false, err
	}

	return (string(data) != ""), nil
}
