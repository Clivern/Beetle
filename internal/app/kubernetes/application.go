// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"context"
)

// GetApplicationVersion gets current application version
func (c *Cluster) GetApplicationVersion(ctx context.Context, namespace, appID, imageFormat string) (string, error) {
	return appID, nil
}
