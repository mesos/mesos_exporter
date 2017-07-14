package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	executor struct {
		ID          string      `json:"executor_id"`
		Name        string      `json:"executor_name"`
		FrameworkID string      `json:"framework_id"`
		Source      string      `json:"source"`
		Statistics  *statistics `json:"statistics"`
		Tasks       []task      `json:"tasks"`
	}

	statistics struct {
		Processes float64 `json:"processes"`
		Threads   float64 `json:"threads"`

		CpusLimit             float64 `json:"cpus_limit"`
		CpusSystemTimeSecs    float64 `json:"cpus_system_time_secs"`
		CpusUserTimeSecs      float64 `json:"cpus_user_time_secs"`
		CpusThrottledTimeSecs float64 `json:"cpus_throttled_time_secs"`
		CpusNrPeriods         float64 `json:"cpus_nr_periods"`
		CpusNrThrottled       float64 `json:"cpus_nr_throttled"`

		MemLimitBytes float64 `json:"mem_limit_bytes"`
		MemRssBytes   float64 `json:"mem_rss_bytes"`
		MemTotalBytes float64 `json:"mem_total_bytes"`
		MemCacheBytes float64 `json:"mem_cache_bytes"`
		MemSwapBytes  float64 `json:"mem_swap_bytes"`

		DiskLimitBytes float64 `json:"disk_limit_bytes"`
		DiskUsedBytes  float64 `json:"disk_used_bytes"`

		NetRxBytes   float64 `json:"net_rx_bytes"`
		NetRxDropped float64 `json:"net_rx_dropped"`
		NetRxErrors  float64 `json:"net_rx_errors"`
		NetRxPackets float64 `json:"net_rx_packets"`
		NetTxBytes   float64 `json:"net_tx_bytes"`
		NetTxDropped float64 `json:"net_tx_dropped"`
		NetTxErrors  float64 `json:"net_tx_errors"`
		NetTxPackets float64 `json:"net_tx_packets"`
	}

	slaveCollector struct {
		*httpClient
		metrics map[*prometheus.Desc]metric
	}

	metric struct {
		valueType prometheus.ValueType
		get       func(*statistics) float64
	}
)

func newSlaveMonitorCollector(httpClient *httpClient) prometheus.Collector {
	labels := []string{"id", "framework_id", "source"}

	return &slaveCollector{
		httpClient: httpClient,
		metrics: map[*prometheus.Desc]metric{
			// Processes
			prometheus.NewDesc(
				"processes",
				"Current number of processes",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.Processes }},
			prometheus.NewDesc(
				"threads",
				"Current number of threads",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.Threads }},

			// CPU
			prometheus.NewDesc(
				"cpus_limit",
				"Current limit of CPUs for task",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.CpusLimit }},
			prometheus.NewDesc(
				"cpu_system_seconds_total",
				"Total system CPU seconds",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.CpusSystemTimeSecs }},
			prometheus.NewDesc(
				"cpu_user_seconds_total",
				"Total user CPU seconds",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.CpusUserTimeSecs }},
			prometheus.NewDesc(
				"cpu_throttled_seconds_total",
				"Total time CPU was throttled due to CFS bandwidth control",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.CpusThrottledTimeSecs }},
			prometheus.NewDesc(
				"cpu_nr_periods_total",
				"Total number of elapsed CFS enforcement intervals",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.CpusNrPeriods }},
			prometheus.NewDesc(
				"cpu_nr_throttled_total",
				"Total number of throttled CFS enforcement intervals.",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.CpusNrThrottled }},

			// Memory
			prometheus.NewDesc(
				"mem_limit_bytes",
				"Current memory limit in bytes",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.MemLimitBytes }},
			prometheus.NewDesc(
				"mem_rss_bytes",
				"Current rss memory usage",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.MemRssBytes }},
			prometheus.NewDesc(
				"mem_total_bytes",
				"Current total memory usage",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.MemTotalBytes }},
			prometheus.NewDesc(
				"mem_cache_bytes",
				"Current page cache memory usage",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.MemCacheBytes }},
			prometheus.NewDesc(
				"mem_swap_bytes",
				"Current swap usage",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.MemSwapBytes }},

			// Disk
			prometheus.NewDesc(
				"disk_limit_bytes",
				"Current disk limit in bytes",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.DiskLimitBytes }},
			prometheus.NewDesc(
				"disk_used_bytes",
				"Current disk usage",
				labels, nil,
			): metric{prometheus.GaugeValue, func(s *statistics) float64 { return s.DiskUsedBytes }},

			// Network
			// - RX
			prometheus.NewDesc(
				"network_receive_bytes_total",
				"Total bytes received",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetRxBytes }},
			prometheus.NewDesc(
				"network_receive_dropped_total",
				"Total packets dropped while receiving",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetRxDropped }},
			prometheus.NewDesc(
				"network_receive_errors_total",
				"Total errors while receiving",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetRxErrors }},
			prometheus.NewDesc(
				"network_receive_packets_total",
				"Total packets received",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetRxPackets }},
			// - TX
			prometheus.NewDesc(
				"network_transmit_bytes_total",
				"Total bytes transmitted",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetTxBytes }},
			prometheus.NewDesc(
				"network_transmit_dropped_total",
				"Total packets dropped while transmitting",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetTxDropped }},
			prometheus.NewDesc(
				"network_transmit_errors_total",
				"Total errors while transmitting",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetTxErrors }},
			prometheus.NewDesc(
				"network_transmit_packets_total",
				"Total packets transmitted",
				labels, nil,
			): metric{prometheus.CounterValue, func(s *statistics) float64 { return s.NetTxPackets }},
		},
	}
}

func (c *slaveCollector) Collect(ch chan<- prometheus.Metric) {
	stats := []executor{}
	c.fetchAndDecode("/monitor/statistics", &stats)

	for _, exec := range stats {
		for desc, m := range c.metrics {
			ch <- prometheus.MustNewConstMetric(desc, m.valueType, m.get(exec.Statistics), exec.ID, exec.FrameworkID, exec.Source)
		}
	}
}

func (c *slaveCollector) Describe(ch chan<- *prometheus.Desc) {
	for metric := range c.metrics {
		ch <- metric
	}
}
