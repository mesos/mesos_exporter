# Prometheus Mesos Exporter

[![Build Status](https://travis-ci.org/mesos/mesos_exporter.svg?branch=master)](https://travis-ci.org/mesos/mesos_exporter)

Exporter for Mesos master and agent metrics.

## Using
The Mesos Exporter can either expose cluster wide metrics from a master or task
metrics from an agent.

```sh
Usage of mesos_exporter:
  -addr string
        Address to listen on (default ":9105")
  -clientCert string
        Path to Mesos client TLS certificate (.pem file)
  -clientKey string
        Path to Mesos client TLS key file (.pem file)
  -enableMasterState
        Enable collection from the master's /state endpoint (default true)
  -exportedSlaveAttributes string
        Comma-separated list of slave attributes to include in the corresponding metric
  -exportedTaskLabels string
        Comma-separated list of task labels to include in the corresponding metric
  -logLevel string
        Log level (default "error")
  -loginURL string
        URL for strict mode authentication (default "https://leader.mesos/acs/api/v1/auth/login")
  -master string
        Expose metrics from master running on this URL
  -password string
        Password for authentication
  -privateKey string
        File path to certificate for strict mode authentication
  -skipSSLVerify
        Skip SSL certificate verification
  -slave string
        Expose metrics from slave running on this URL
  -strictMode
        Use strict mode authentication
  -timeout duration
        Master polling timeout (default 10s)
  -trustedCerts string
        Comma-separated list of certificates (.pem files) trusted for requests to Mesos endpoints
  -username string
        Username for authentication
  -version
        Show version
```

When using HTTP or strict mode authentication, the following values are read from the environment, if they are not specified at run time:
- `MESOS_EXPORTER_USERNAME`
- `MESOS_EXPORTER_PASSWORD`
- `MESOS_EXPORTER_PRIVATE_KEY`

When collecting metrics from the master, the `-enableMasterState`
flag will enable the Mesos Exporter to fetch the master's
[state](http://mesos.apache.org/documentation/latest/endpoints/master/state/)
endpoint in order to publish metrics about the resources available
on registered agents. In large clusters, polling this endpoint can
degrade master performance. In this case, `-enableMasterState` can
be disabled on the master exporter and equivalent metrics can be
collected by running the Mesos Exporter on each agent.

When `-enableMasterState` is true, the master exporter will publish
the following additional metrics labeled agent ID:

| mesos_slave_cpus |
| mesos_slave_cpus_unreserved |
| mesos_slave_cpus_used |
| mesos_slave_disk_bytes |
| mesos_slave_disk_unreserved_bytes |
| mesos_slave_disk_used_bytes |
| mesos_slave_mem_bytes |
| mesos_slave_mem_unreserved_bytes|
| mesos_slave_mem_used_bytes |
| mesos_slave_ports |
| mesos_slave_ports_unreserved |
| mesos_slave_ports_used |

## Prometheus Configuration

Usually you would run one exporter with `-master` for each master and one
exporter for each slave with `-slave`. Monitoring each master individually
ensures that the cluster can be monitored even if the underlying Mesos cluster
is in a degraded state.

- Master: `mesos_exporter -master http://localhost:5050`
- Agent: `mesos_exporter -slave http://localhost:5051`

The necessary Prometheus configuration could look like this:

```
- job_name: mesos-master
  scrape_interval: 15s
  scrape_timeout: 10s
  static_configs:
  - targets:
    - master1.mesos.example.org:9105
    - master2.mesos.example.org:9105
    - master3.mesos.example.org:9105

- job_name: mesos-slave
  scrape_interval: 15s
  scrape_timeout: 10s
  static_configs:
  - targets:
    - node1.mesos.example.org:9105
    - node2.mesos.example.org:9105
    - node3.mesos.example.org:9105
```


A minimal set of alerts to ensure your cluster is operational could then be defined
as follows:

```
ALERT MesosDown
  IF (up{job=~"mesos.*"} == 0) or (irate(mesos_collector_errors_total[5m]) > 0)
  FOR 5m
  LABELS { severity="warning" }
  ANNOTATIONS {
    description="Either the exporter or the associated Mesos component is down.",
    summary="The Mesos instance {{$labels.instance}} cannot be scraped."
  }

ALERT MesosMasterLeader
  IF sum(mesos_master_elected{job="mesos-master"}) != 1
  FOR 5m
  LABELS { severity="page" }
  ANNOTATIONS {
    description="Agents and frameworks require a unique leading Mesos master.",
    summary="Expected one leading Mesos master but there are {{ $value }}."
  }

ALERT MesosMasterTooManyRestarts
  IF resets(mesos_master_uptime_seconds{job="mesos-master"}[1h]) > 10
  FOR 5m
  LABELS { severity="page" }
  ANNOTATIONS {
    description="The number of seconds the process has been running is resetting regularly.",
    summary="The Mesos master {{$labels.instance}} has restarted {{ $value }} times in the last hour."
  }

ALERT MesosSlaveActive
  IF sum(mesos_master_slaves_state{state="active"}) < 0.9 * count(up{job="mesos-slave"})
  FOR 5m
  LABELS { severity="page" }
  ANNOTATIONS {
    description="Mesos agents must be registered with the master in order to receive tasks.",
    summary="More than 10% of all Mesos agents dropped out. Only {{ $value }} active agents remaining."
  }

ALERT MesosSlaveTooManyRestarts
  IF resets(mesos_slave_uptime_seconds{job="mesos-slave"}[1h]) > 10
  FOR 5m
  LABELS { severity="page" }
  ANNOTATIONS {
    description="The number of seconds the process has been running is resetting regularly.",
    summary="The Mesos agent {{$labels.instance}} has restarted {{ $value }} times in the last hour."
  }
```
