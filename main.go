package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

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

func main() {
	fs := flag.NewFlagSet("mesos-exporter", flag.ExitOnError)
	addr := fs.String("addr", ":9110", "Address to listen on")
	masterURL := fs.String("master", "", "Expose metrics from master running on this URL")
	slaveURL := fs.String("slave", "", "Expose metrics from slave running on t his URL")
	timeout := fs.Duration("timeout", 5*time.Second, "Master polling timeout")

	fs.Parse(os.Args[1:])
	if *masterURL != "" && *slaveURL != "" {
		log.Fatal("Only -master or -slave can be given at a time")
	}

	switch {
	case *masterURL != "":
		for _, c := range []prometheus.Collector{
			newMasterCollector(*masterURL, *timeout),
			newMasterStateCollector(*masterURL, *timeout),
		} {
			if err := prometheus.Register(c); err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("Exposing master metrics on %s", *addr)

	case *slaveURL != "":
		for _, c := range []prometheus.Collector{
			newSlaveCollector(*slaveURL, *timeout),
			newSlaveMonitorCollector(*slaveURL, *timeout),
		} {
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
