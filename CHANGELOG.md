# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [UNRELEASED]
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
