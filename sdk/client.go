// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package sdk

import (
	"github.com/clivern/beetle/core/module"
)

// Client struct
type Client struct {
	APIKey     string
	APIURL     string
	HTTPClient *module.HTTPClient
}

// SetHTTPClient sets http client
func (c *Client) SetHTTPClient(httpClient *module.HTTPClient) {
	c.HTTPClient = httpClient
}

// SetAPIURL sets api url
func (c *Client) SetAPIURL(APIURL string) {
	c.APIURL = APIURL
}

// SetAPIKey sets api key
func (c *Client) SetAPIKey(APIKey string) {
	c.APIKey = APIKey
}
