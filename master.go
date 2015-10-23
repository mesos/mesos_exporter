package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	resources struct {
		CPUs  float64 `json:"cpus"`
		Disk  float64 `json:"disk"`
		Mem   float64 `json:"mem"`
		Ports ranges  `json:"ports"`
	}

	task struct {
		Name        string    `json:"name"`
		ID          string    `json:"id"`
		ExecutorID  string    `json:"executor_id"`
		FrameworkID string    `json:"framework_id"`
		SlaveID     string    `json:"slave_id"`
		State       string    `json:"state"`
		Labels      []label   `json:"labels"`
		Resources   resources `json:"resources"`
		Statuses    []status  `json:"statuses"`
	}

	label struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	status struct {
		State     string  `json:"state"`
		Timestamp float64 `json:"timestamp"`
	}

	slave struct {
		PID        string    `json:"pid"`
		Used       resources `json:"used_resources"`
		Unreserved resources `json:"unreserved_resources"`
		Total      resources `json:"resources"`
	}

	framework struct {
		Active    bool   `json:"active"`
		Tasks     []task `json:"tasks"`
		Completed []task `json:"completed_tasks"`
	}

	state struct {
		Slaves     []slave     `json:"slaves"`
		Frameworks []framework `json:"frameworks"`
	}

	masterCollector struct {
		*http.Client
		url     string
		metrics map[prometheus.Collector]func(*state, prometheus.Collector)
	}
)

func newMasterCollector(url string, timeout time.Duration) *masterCollector {
	labels := []string{"slave"}
	return &masterCollector{
		Client: &http.Client{Timeout: timeout},
		url:    url,
		metrics: map[prometheus.Collector]func(*state, prometheus.Collector){
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave CPUs (fractional)",
				Name:      "cpus",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Total.CPUs)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave CPUs (fractional)",
				Name:      "cpus_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Used.CPUs)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave CPUs (fractional)",
				Name:      "cpus_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Unreserved.CPUs)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave memory in MB",
				Name:      "mem",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Total.Mem)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave memory in MB",
				Name:      "mem_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Used.Mem)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave memory in MB",
				Name:      "mem_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Unreserved.Mem)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave disk in MB",
				Name:      "disk",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Total.Disk)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave disk in MB",
				Name:      "disk_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Used.Disk)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave disk in MB",
				Name:      "disk_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(s.Unreserved.Disk)
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Total slave ports",
				Name:      "ports",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					size := s.Total.Ports.size()
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(float64(size))
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Used slave ports",
				Name:      "ports_used",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					size := s.Used.Ports.size()
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(float64(size))
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Unreserved slave ports",
				Name:      "ports_unreserved",
				Namespace: "mesos",
				Subsystem: "slave",
			}, labels): func(st *state, c prometheus.Collector) {
				for _, s := range st.Slaves {
					size := s.Unreserved.Ports.size()
					c.(*prometheus.GaugeVec).WithLabelValues(s.PID).Set(float64(size))
				}
			},
			prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Help:      "Framework tasks",
				Name:      "tasks",
				Namespace: "mesos",
				Subsystem: "slave",
			}, []string{"slave", "executor", "name", "framework", "timestamp", "state"}): func(st *state, c prometheus.Collector) {
				for _, f := range st.Frameworks {
					if !f.Active {
						continue
					}
					tasks := f.tasks()
					for _, task := range tasks {
						for _, status := range task.Statuses {
							values := []string{
								task.SlaveID,
								task.ExecutorID,
								task.Name,
								task.FrameworkID,
								strconv.FormatFloat(status.Timestamp, 'f', -1, 64),
								status.State,
							}
							c.(*prometheus.GaugeVec).WithLabelValues(values...).Set(float64(len(tasks)))
						}
					}
				}
			},
		},
	}
}

func (c *masterCollector) Collect(ch chan<- prometheus.Metric) {
	res, err := c.Get(c.url + "/state.json")
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

	for c, set := range c.metrics {
		set(&s, c)
		c.Collect(ch)
	}

	res, err = c.Get(c.url + "/metrics/snapshot")
	if err != nil {
		log.Print(err)
		return
	}
	defer res.Body.Close()
	snapshotCollect(ch, res.Body)
}

func (c *masterCollector) Describe(ch chan<- *prometheus.Desc) {
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

func (f framework) tasks() []task {
	tasks := make([]task, len(f.Tasks)+len(f.Completed))
	tasks = append(tasks, f.Tasks...)
	tasks = append(tasks, f.Completed...)
	return tasks
}
