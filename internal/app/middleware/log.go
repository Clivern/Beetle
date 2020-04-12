// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "beetle",
			Name:      "total_http_requests",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		}, []string{"code", "method", "handler", "host", "url"})

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "beetle",
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"code", "method", "url"},
	)

	responseSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "beetle",
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(responseSize)
}

// Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		var bodyBytes []byte
		start := time.Now()

		// Workaround for issue https://github.com/gin-gonic/gin/issues/1651
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Info(fmt.Sprintf(`Incoming Request %s:%s Body: %s`, c.Request.Method, c.Request.URL, string(bodyBytes)))

		c.Next()

		// after request
		elapsed := float64(time.Since(start)) / float64(time.Second)
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Collect Metrics
		httpRequests.WithLabelValues(
			strconv.Itoa(c.Writer.Status()),
			c.Request.Method,
			c.HandlerName(),
			c.Request.Host,
			c.Request.URL.Path,
		).Inc()

		requestDuration.WithLabelValues(
			strconv.Itoa(c.Writer.Status()),
			c.Request.Method,
			c.Request.URL.Path,
		).Observe(elapsed)

		responseSize.Observe(float64(c.Writer.Size()))

		log.WithFields(log.Fields{
			"CorrelationId": c.Request.Header.Get("X-Correlation-ID"),
		}).Info(fmt.Sprintf(`Outgoing Response Code %d, Size %d`, status, size))
	}
}
