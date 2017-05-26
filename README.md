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
  -username string
        Username to use for HTTP or strict mode authentication
  -password string
        Password to use for HTTP or strict mode authentication
  -loginUrl
        URL for strict mode authentication (default https://leader.mesos/acs/api/v1/auth/login).
  -trustedCerts string
        Comma-separated list of certificates (.pem files) trusted for requests to Mesos endpoints
  -strictMode
        Enable strict mode API authentication
  -privateKey
        Private key used for strict mode authentication. This must be provided
        when using strict mode. However, it can be read from the environment if the secret store is used. 
  -skipSSLVerify
        Disable SSL certificate verification
```
When using HTTP or strict mode authentication, the following values are read from the environment, if they are not specified at run time:
- `MESOS_EXPORTER_USERNAME`
- `MESOS_EXPORTER_PASSWORD`
- `MESOS_EXPORTER_PRIVATE_KEY`

---

Usually you would run one exporter with `-master` pointing to the current
leader and one exporter for each slave with `-slave` pointing to it. In
a default Mesos / DC/OS setup, you should be able to run the mesos-exporter
like this:

- Master: `mesos-exporter -master http://leader.mesos:5050`
- Agent: `mesos-exporter -slave http://localhost:5051`
