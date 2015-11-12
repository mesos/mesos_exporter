package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func newSlaveCollector(url string, timeout time.Duration) *metricCollector {
	metrics := map[prometheus.Collector]func(metricMap, prometheus.Collector) error{}

	return &metricCollector{
		Client:  &http.Client{Timeout: timeout},
		url:     url,
		metrics: metrics,
	}
}
