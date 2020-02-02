// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		var bodyBytes []byte

		// Workaround for issue https://github.com/gin-gonic/gin/issues/1651
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		fmt.Println(
			fmt.Sprintf(
				`{"message": "Incoming Request %s:%s Body: %s", "correlationId": "%s"}`,
				c.Request.Method,
				c.Request.URL,
				string(bodyBytes),
				c.Request.Header.Get("X-Correlation-ID"),
			),
		)

		c.Next()

		// after request
		status := c.Writer.Status()
		size := c.Writer.Size()

		fmt.Println(
			fmt.Sprintf(
				`{"message": "Outgoing Response Code %d, Size %d", "correlationId": "%s"}`,
				status,
				size,
				c.Request.Header.Get("X-Correlation-ID"),
			),
		)
	}
}
