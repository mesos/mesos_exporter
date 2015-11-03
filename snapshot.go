package main

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type snapshotMetrics map[string]float64

func snapshotCollect(ch chan<- prometheus.Metric, r io.Reader) {
	var metrics snapshotMetrics
	if err := json.NewDecoder(r).Decode(&metrics); err != nil {
		log.Print(err)
		return
	}

	for name, value := range metrics {
		mn := strings.Join(append([]string{"mesos"}, strings.Split(name, "/")...), "_")
		desc := prometheus.NewDesc(mn, "Exposed from /metrics/snapshot", nil, nil)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, value)
	}
}
