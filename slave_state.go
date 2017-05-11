// Scrape the /slave(1)/state endpoint to get information on the tasks running
// on executors. Information scraped at this point:
//
// * Labels of running tasks ("mesos_slave_task_labels" series)
// * Attributes of mesos slaves ("mesos_slave_attributes")
package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	slaveFramework struct {
		ID        string     `json:"ID"`
		Executors []executor `json:"executors"`
	}

	// similar to /master/state's 'slave', but with small differences
	slaveState struct {
		PID        string            `json:"pid"`
		Attributes map[string]string `json:"attributes"`
		Frameworks []slaveFramework  `json:"frameworks"`
	}

	slaveStateCollector struct {
		*httpClient
		metrics map[prometheus.Collector]func(*slaveState, prometheus.Collector)
	}
)

func newSlaveStateCollector(httpClient *httpClient, userTaskLabelList []string, slaveAttributeLabelList []string) *slaveStateCollector {
	var metrics = map[prometheus.Collector]func(*slaveState, prometheus.Collector){}
	defaultLabels := []string{"slave"}

	if len(userTaskLabelList) > 0 {
		defaultTaskLabels := []string{"source", "framework_id", "executor_id"}
		normalisedUserTaskLabelList := normaliseLabelList(userTaskLabelList)
		taskLabelList := append(defaultLabels, defaultTaskLabels...)
		taskLabelList = append(taskLabelList, normalisedUserTaskLabelList...)

		metrics[prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Help:      "Task labels",
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "task_labels",
		}, taskLabelList)] = func(st *slaveState, c prometheus.Collector) {
			for _, f := range st.Frameworks {
				for _, e := range f.Executors {
					for _, t := range e.Tasks {
						// Default labels
						taskLabels := map[string]string{
							"slave":        st.PID,
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
							if stringinSlice(normalisedLabel, normalisedUserTaskLabelList) {
								taskLabels[normalisedLabel] = label.Value
							}
						}
						c.(*prometheus.GaugeVec).With(taskLabels).Set(1)
					}
				}
			}
		}
	}

	if len(slaveAttributeLabelList) > 0 {
		normalisedAttributeLabels := normaliseLabelList(slaveAttributeLabelList)
		AttributeLabelList := append(defaultLabels, normalisedAttributeLabels...)
		metrics[prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Help:      "Slave attributes",
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "attributes",
		}, AttributeLabelList)] = func(st *slaveState, c prometheus.Collector) {
			slaveAttributesExport := map[string]string{
				"slave": st.PID,
			}

			// (Empty) user labels
			for _, label := range normalisedAttributeLabels {
				slaveAttributesExport[label] = ""
			}
			for key, value := range st.Attributes {
				normalisedLabel := normaliseLabel(key)
				if stringinSlice(normalisedLabel, normalisedAttributeLabels) {
					slaveAttributesExport[normalisedLabel] = value
				}
			}
			c.(*prometheus.GaugeVec).With(slaveAttributesExport).Set(1)
		}
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
