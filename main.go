package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var errorCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "mesos",
	Subsystem: "collector",
	Name:      "errors_total",
	Help:      "Total number of internal mesos-collector errors.",
})

func init() {
	prometheus.MustRegister(errorCounter)
}

func mkHttpClient(url string, timeout time.Duration, auth authInfo) *httpClient {
	return &httpClient{
		http.Client{Timeout: timeout},
		url,
		auth,
	}
}

func main() {
	fs := flag.NewFlagSet("mesos-exporter", flag.ExitOnError)
	addr := fs.String("addr", ":9110", "Address to listen on")
	masterURL := fs.String("master", "", "Expose metrics from master running on this URL")
	slaveURL := fs.String("slave", "", "Expose metrics from slave running on this URL")
	timeout := fs.Duration("timeout", 5*time.Second, "Master polling timeout")
	exportedTaskLabels := fs.String("exportedTaskLabels", "", "Comma-separated list of task labels to include in the task_labels metric")
	ignoreCompletedFrameworkTasks := fs.Bool("ignoreCompletedFrameworkTasks", false, "Don't export task_state_time metric");

	fs.Parse(os.Args[1:])
	if *masterURL != "" && *slaveURL != "" {
		log.Fatal("Only -master or -slave can be given at a time")
	}

	auth := authInfo{
		os.Getenv("MESOS_EXPORTER_USERNAME"),
		os.Getenv("MESOS_EXPORTER_PASSWORD"),
	}

	switch {
	case *masterURL != "":
		for _, f := range []func(*httpClient) prometheus.Collector{
			newMasterCollector,
			func(c *httpClient) prometheus.Collector {
				return newMasterStateCollector(c, *ignoreCompletedFrameworkTasks)
			},
		} {
			c := f(mkHttpClient(*masterURL, *timeout, auth));
			if err := prometheus.Register(c); err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("Exposing master metrics on %s", *addr)

	case *slaveURL != "":
		slaveCollectors := []func(*httpClient) prometheus.Collector{
			func(c *httpClient) prometheus.Collector {
				return newSlaveCollector(c)
			},
			func(c *httpClient) prometheus.Collector {
				return newSlaveMonitorCollector(c)
			},
		}
		if *exportedTaskLabels != "" {
			slaveLabels := strings.Split(*exportedTaskLabels, ",");
			slaveCollectors = append(slaveCollectors, func (c *httpClient) prometheus.Collector{
				return newSlaveStateCollector(c, slaveLabels)
			})
		}

		for _, f := range slaveCollectors {
			c := f(mkHttpClient(*slaveURL, *timeout, auth));
			if err := prometheus.Register(c); err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("Exposing slave metrics on %s", *addr)

	default:
		log.Fatal("Either -master or -slave is required")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>Mesos Exporter</title></head>
            <body>
            <h1>Mesos Exporter</h1>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})
	http.Handle("/metrics", prometheus.Handler())
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
