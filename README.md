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
  -addr=":9110": Address to listen on
  -master="": Expose metrics from master running on this URL
  -slave="": Expose metrics from agent running on this URL
  -timeout=5s: Master polling timeout
```

Usually you would run one exporter with `-master` pointing to the current
leader and one exporter for each slave with `-slave` pointing to it. In
a default Mesos / DC/OS setup, you should be able to run the mesos-exporter
like this:

- Master: `mesos-exporter -master http://leader.mesos:5050`
- Agent: `mesos-exporter -slave http://localhost:5051`
