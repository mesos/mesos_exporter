// Scrape the /slave(1)/state endpoint to get information on the tasks running
// on executors. Information scraped at this point:
//
// * Labels of running tasks ("mesos_slave_task_labels" series)
package main

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	slaveFramework struct {
		ID        string     `json:"ID"`
		Executors []executor `json:"executors"`
	}

	slaveState struct {
		Frameworks []slaveFramework `json:"frameworks"`
	}

	slaveStateCollector struct {
		*httpClient
		metrics map[prometheus.Collector]func(*slaveState, prometheus.Collector)
	}
)

// Task labels must be alphanumeric, no leading digits.
var invalidLabelNameCharRE = regexp.MustCompile("[^a-zA-Z_0-9]")

// Replace invalid task label digits by underscores
func normaliseLabel(label string) string {
	if len(label) > 0 && '0' <= label[0] && label[0] <= '9' {
		return "_" + invalidLabelNameCharRE.ReplaceAllString(label[1:], "_")
	}
	return invalidLabelNameCharRE.ReplaceAllString(label, "_")
}

// Return true if `needle` is in `haystack`
func inArray(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if needle == elem {
			return true
		}
	}
	return false
}

func newSlaveStateCollector(httpClient *httpClient, userTaskLabelList []string) *slaveStateCollector {
	defaultLabels := []string{"source", "framework_id", "executor_id"}

	// Sanitise user-supplied list of task labels that should be included in the series
	normalisedUserTaskLabelList := []string{}
	for _, label := range userTaskLabelList {
		normalisedUserTaskLabelList = append(normalisedUserTaskLabelList, normaliseLabel(label))
	}

	taskLabelList := append(defaultLabels, normalisedUserTaskLabelList...)

	metrics := map[prometheus.Collector]func(*slaveState, prometheus.Collector){
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Help:      "Task labels",
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "task_labels",
		}, taskLabelList): func(st *slaveState, c prometheus.Collector) {
			for _, f := range st.Frameworks {
				for _, e := range f.Executors {
					for _, t := range e.Tasks {
						// Default labels
						taskLabels := map[string]string{
							"source":       e.Source,
							"framework_id": f.ID,
							"executor_id":  e.ID,
						}
						// User labels
						for _, label := range normalisedUserTaskLabelList {
							taskLabels[label] = ""
						}
						for _, label := range t.Labels {
							normalisedLabel := normaliseLabel(label.Key)
							// Ignore labels not explicitly whitelisted by user
							if inArray(normalisedLabel, normalisedUserTaskLabelList) {
								taskLabels[normalisedLabel] = label.Value
							}
						}
						c.(*prometheus.GaugeVec).With(taskLabels).Set(1)
					}
				}
			}
		},
	}

	return &slaveStateCollector{httpClient, metrics}
}

func (c *slaveStateCollector) Collect(ch chan<- prometheus.Metric) {
	var s slaveState
	c.fetchAndDecode("/slave(1)/state", &s)
	for c, set := range c.metrics {
		set(&s, c)
		c.Collect(ch)
	}
}

func (c *slaveStateCollector) Describe(ch chan<- *prometheus.Desc) {
	for metric := range c.metrics {
		metric.Describe(ch)
	}
}
