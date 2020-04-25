// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "# namespace",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("namespace", args) },
}

func init() {
	shellCmd.AddCommand(namespaceCmd)
}
