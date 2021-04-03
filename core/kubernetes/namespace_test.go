// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/clivern/beetle/core/module"
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestNamespace test cases
func TestNamespace(t *testing.T) {
	testingConfig := "config.testing.yml"

	// LoadConfigFile
	t.Run("LoadConfigFile", func(t *testing.T) {
		fs := module.FileSystem{}

		dir, _ := os.Getwd()
		configFile := fmt.Sprintf("%s/%s", dir, testingConfig)

		for {
			if fs.FileExists(configFile) {
				break
			}
			dir = filepath.Dir(dir)
			configFile = fmt.Sprintf("%s/%s", dir, testingConfig)
		}

		t.Logf("Load Config File %s", configFile)

		configUnparsed, _ := ioutil.ReadFile(configFile)
		configParsed, _ := envsubst.EvalEnv(string(configUnparsed))
		viper.SetConfigType("yaml")
		viper.ReadConfig(bytes.NewBuffer([]byte(configParsed)))
	})

	// TestGetNamespaces
	t.Run("TestGetNamespaces", func(t *testing.T) {
		cluster, err := GetCluster("production")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, cluster.Name, "production")
		pkg.Expect(t, cluster.Kubeconfig, "/app/configs/production-cluster-kubeconfig.yaml")

		cluster.Override(
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					UID:  "9d0cdf8a-dedc-11e9-bf91-42010a800167",
				},
			},
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "beetle",
					UID:  "9d0cdf8a-dedc-11e9-bf91-42010a800168",
				},
			},
		)

		namespaces, err := cluster.GetNamespaces(context.TODO())

		pkg.Expect(t, nil, err)
		pkg.Expect(t, namespaces[1].Name, "default")
		pkg.Expect(t, namespaces[1].UID, "9d0cdf8a-dedc-11e9-bf91-42010a800167")
		pkg.Expect(t, namespaces[1].Status, "")

		pkg.Expect(t, namespaces[0].Name, "beetle")
		pkg.Expect(t, namespaces[0].UID, "9d0cdf8a-dedc-11e9-bf91-42010a800168")
		pkg.Expect(t, namespaces[0].Status, "")
	})

	// TestGetNamespaceBeetle
	t.Run("TestGetNamespaceBeetle", func(t *testing.T) {
		cluster, err := GetCluster("production")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, cluster.Name, "production")
		pkg.Expect(t, cluster.Kubeconfig, "/app/configs/production-cluster-kubeconfig.yaml")

		cluster.Override(
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					UID:  "9d0cdf8a-dedc-11e9-bf91-42010a800167",
				},
			},
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "beetle",
					UID:  "9d0cdf8a-dedc-11e9-bf91-42010a800168",
				},
			},
		)

		namespace, err := cluster.GetNamespace(context.TODO(), "beetle")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, namespace.Name, "beetle")
		pkg.Expect(t, namespace.UID, "9d0cdf8a-dedc-11e9-bf91-42010a800168")
		pkg.Expect(t, namespace.Status, "")
	})

	// TestGetNamespaceDefault
	t.Run("TestGetNamespaceDefault", func(t *testing.T) {
		cluster, err := GetCluster("production")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, cluster.Name, "production")
		pkg.Expect(t, cluster.Kubeconfig, "/app/configs/production-cluster-kubeconfig.yaml")

		cluster.Override(
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					UID:  "9d0cdf8a-dedc-11e9-bf91-42010a800167",
				},
			},
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "beetle",
					UID:  "9d0cdf8a-dedc-11e9-bf91-42010a800168",
				},
			},
		)

		namespace, err := cluster.GetNamespace(context.TODO(), "default")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, namespace.Name, "default")
		pkg.Expect(t, namespace.UID, "9d0cdf8a-dedc-11e9-bf91-42010a800167")
		pkg.Expect(t, namespace.Status, "")
	})
}
