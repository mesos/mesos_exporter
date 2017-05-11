# Prometheus Mesos Exporter
Exporter for Mesos master and agent metrics

## Installing
```sh
$ go get github.com/mesosphere/mesos-exporter
```

## Using
The Mesos Exporter can either expose cluster wide metrics from a master or task
metrics from an agent.

```sh
Usage of mesos-exporter:
  -addr string
        Address to listen on (default ":9110")
  -exportedSlaveAttributes string
        Comma-separated list of slave attributes to include in the corresponding metric
  -exportedTaskLabels string
        Comma-separated list of task labels to include in the corresponding metric
  -ignoreCompletedFrameworkTasks
        Don't export task_state_time metric
  -master string
        Expose metrics from master running on this URL
  -slave string
        Expose metrics from slave running on this URL
  -timeout duration
        Master polling timeout (default 5s)
  -trustedCerts string
        Comma-separated list of certificates (.pem files) trusted for requests to
        Mesos endpoints
```
When using HTTP authentication, the following values are read from the environment:
- `MESOS_EXPORTER_USERNAME`
- `MESOS_EXPORTER_PASSWORD`

---

Usually you would run one exporter with `-master` pointing to the current
leader and one exporter for each slave with `-slave` pointing to it. In
a default Mesos / DC/OS setup, you should be able to run the mesos-exporter
like this:

- Master: `mesos-exporter -master http://leader.mesos:5050`
- Agent: `mesos-exporter -slave http://localhost:5051`
