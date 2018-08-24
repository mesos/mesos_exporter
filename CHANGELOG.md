# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [1.1.0] - 2018-08-23
### Added
- Added a new `-enableMasterState` flag to prevent collection of metrics that
  require polling the Mesos mater `/state` endpoint.
- Added new `-clientCertFile` and `-clientKeyFile` flags that allow the
  exporter to use TLS client authentication.
- Added exporter version information to the exported metrics, the HTTP user
  agent, and the `-version` flag.
- Updated the exported metrics to include new metrics from recent
  Mesos versions, including allocator metrics, message type metrics and
  improved resource, task and agent state metrics.
- Added strict authentication mode for use in DC/OS clusters.

### Changed
- Moved the `tasks_killing` metrics label from the `mesos_master_task_states_exit_total`
  metric to `mesos_master_task_states_current` since it is a gauge, not a counter.
- Correctly marked the `mesos_master_task_states_current` metric as a gauge,
  not a counter.
- Fixed embedded whitespace in master metric names.
- Improved the exporter error logging to include contextual information.
- Improved the JSON unmarshalling to correctly accept Mesos agent attributes.

### Removed
- Several built-in metrics of the go Prometheus client have been removed. This includes:
  `http_request_duration_microseconds`, `http_request_size_bytes`, `http_requests_total`,
  and `http_response_size_bytes`.

## [1.0.0] - 2017-09-02
### Changed
- Build releases with cgo enabled. This allows the exporter to use the system cert store
  which can simplify deployments significantly.

## [1.0.0-rc2] - 2017-09-01
### Added
- First release with changelog. All changes before have been untracked.
