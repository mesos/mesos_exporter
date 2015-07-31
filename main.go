package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	fs := flag.NewFlagSet("mesos-exporter", flag.ExitOnError)
	addr := fs.String("addr", ":9110", "Address to listen on")
	port := fs.Int("port", 5050, "Master port")
	name := fs.String("name", "leader.mesos.", "Leading Mesos master DNS name")
	timeout := fs.Duration("timeout", 5*time.Second, "Master polling timeout")

	fs.Parse(os.Args[1:])

	c := newCollector(*name, *port, *timeout)
	if err := prometheus.Register(c); err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", prometheus.Handler())
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

type collector struct {
	*http.Client
	name    string
	port    int
	metrics map[*prometheus.GaugeVec]func(*slave, prometheus.Gauge)
}

func newCollector(name string, port int, timeout time.Duration) *collector {
	labels := []string{"slave"}
	return &collector{
		Client: &http.Client{Timeout: timeout},
		name:   name,
		port:   port,
		metrics: map[*prometheus.GaugeVec]func(*slave, prometheus.Gauge){
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave CPUs (fractional)",
				Name:      "cpus",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Total.CPUs) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave CPUs (fractional)",
				Name:      "cpus_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Used.CPUs) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave CPUs (fractional)",
				Name:      "cpus_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Unreserved.CPUs) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave memory in MB",
				Name:      "mem",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Total.Mem) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave memory in MB",
				Name:      "mem_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Used.Mem) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave memory in MB",
				Name:      "mem_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Unreserved.Mem) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave disk in MB",
				Name:      "disk",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Total.Disk) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave disk in MB",
				Name:      "disk_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Used.Disk) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave disk in MB",
				Name:      "disk_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(s.Unreserved.Disk) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave ports",
				Name:      "ports",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(float64(s.Total.Ports.size())) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave ports",
				Name:      "ports_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(float64(s.Used.Ports.size())) },
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave ports",
				Name:      "ports_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(s *slave, g prometheus.Gauge) { g.Set(float64(s.Unreserved.Ports.size())) },
		},
	}
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	masters, err := net.LookupHost(c.name)
	if err != nil || len(masters) == 0 {
		log.Printf("failed to DNS lookup %s: %s", c.name, err)
		return
	}

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(masters[0], strconv.Itoa(c.port)),
			Path:   "/state.json",
		},
	}

	res, err := c.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer res.Body.Close()

	var s state
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		log.Print(err)
		return
	}

	for _, slave := range s.Slaves {
		for metric, set := range c.metrics {
			m := metric.WithLabelValues(slave.PID)
			set(&slave, m)
			ch <- m
		}
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for metric := range c.metrics {
		metric.Describe(ch)
	}
}

type ranges [][2]uint64

func (rs *ranges) UnmarshalJSON(data []byte) (err error) {
	if data = bytes.Trim(data, `[]"`); len(data) == 0 {
		return nil
	}

	var rng [2]uint64
	for _, r := range bytes.Split(data, []byte(",")) {
		ps := bytes.SplitN(r, []byte("-"), 2)
		if len(ps) != 2 {
			return fmt.Errorf("bad range: %s", r)
		}

		rng[0], err = strconv.ParseUint(string(bytes.TrimSpace(ps[0])), 10, 64)
		if err != nil {
			return err
		}

		rng[1], err = strconv.ParseUint(string(bytes.TrimSpace(ps[1])), 10, 64)
		if err != nil {
			return err
		}

		*rs = append(*rs, rng)
	}

	return nil
}

func (rs ranges) size() uint64 {
	var sz uint64
	for i := range rs {
		sz += 1 + (rs[i][1] - rs[i][0])
	}
	return sz
}

type resources struct {
	CPUs  float64 `json:"cpus"`
	Disk  float64 `json:"disk"`
	Mem   float64 `json:"mem"`
	Ports ranges  `json:"ports"`
}

type slave struct {
	PID        string    `json:"pid"`
	Used       resources `json:"used_resources"`
	Unreserved resources `json:"unreserved_resources"`
	Total      resources `json:"resources"`
}

type state struct {
	Slaves []slave `json:"slaves"`
}
