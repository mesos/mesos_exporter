package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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
		if err := prometheus.Register(newMasterCollector(*masterURL, *timeout)); err != nil {
			log.Fatal(err)
		}

		if err := prometheus.Register(newMasterStateCollector(*masterURL, *timeout)); err != nil {
			log.Fatal(err)
		}
		log.Printf("Exposing master metrics on %s", *addr)
	case *slaveURL != "":
		if err := prometheus.Register(newSlaveCollector(*slaveURL, *timeout)); err != nil {
			log.Fatal(err)
		}
		if err := prometheus.Register(newSlaveMonitorCollector(*slaveURL, *timeout)); err != nil {
			log.Fatal(err)
		}
		log.Printf("Exposing slave metrics on %s", *addr)

	default:
		log.Fatal("Either -master or -slave is required")
	}

	http.Handle("/metrics", prometheus.Handler())
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
