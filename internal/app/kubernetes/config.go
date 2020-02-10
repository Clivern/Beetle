// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"github.com/spf13/viper"
)

// Clusters struct
type Clusters struct {
	Clusters []*Cluster `mapstructure:",clusters"`
}

// Cluster struct
type Cluster struct {
	Name    string `mapstructure:",name"`
	Config  string `mapstructure:",config"`
	Version string `mapstructure:",version"`
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

// Info fetch the cluster info
func (c *Cluster) Info() (bool, error) {
	// clus, _ := kubernetes.GetClusters()
	// clus[0].Info()
	// fmt.Println(clus[0].Version)
	c.Version = "1.17"

	return true, nil
}
