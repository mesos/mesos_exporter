# Prometheus Mesos Exporter
Exporter for Mesos master and agent metrics

## Installing
```sh
$ go get github.com/mesosphere/mesos_exporter
```

## Using
The Mesos Exporter can either expose cluster wide metrics from a master or task
metrics from an agent.

```sh
Usage of mesos_exporter:
  -addr string
       	Address to listen on (default ":9110")
  -ignoreCompletedFrameworkTasks
       	Don't export task_state_time metric
  -master string
       	Expose metrics from master running on this URL
  -slave string
       	Expose metrics from slave running on this URL
  -timeout duration
       	Master polling timeout (default 5s)
```

Usually you would run one exporter with `-master` pointing to the current
leader and one exporter for each slave with `-slave` pointing to it. In
a default Mesos / DC/OS setup, you should be able to run the mesos-exporter
like this:

- Master: `mesos_exporter -master http://leader.mesos:5050`
- Agent: `mesos_exporter -slave http://localhost:5051`
