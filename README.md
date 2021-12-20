# audit-exporter

[![Build Status](https://img.shields.io/drone/build/thegeeklab/audit-exporter?logo=drone&server=https%3A%2F%2Fdrone.thegeeklab.de)](https://drone.thegeeklab.de/thegeeklab/audit-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/audit-exporter)](https://goreportcard.com/report/github.com/thegeeklab/audit-exporter)
[![Codecov](https://img.shields.io/codecov/c/github/thegeeklab/audit-exporter)](https://codecov.io/gh/thegeeklab/audit-exporter)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/audit-exporter)](https://github.com/thegeeklab/audit-exporter/graphs/contributors)
[![License: MIT](https://img.shields.io/github/license/thegeeklab/audit-exporter)](https://github.com/thegeeklab/audit-exporter/blob/main/LICENSE)

## Installation

Prebuild multiarch binaries are availabe for Linux only:

```Shell
curl -L https://github.com/thegeeklab/audit-exporter/releases/download/v0.1.0/audit-exporter-0.1.0-linux-amd64 > /usr/local/bin/audit-exporter
chmod +x /usr/local/bin/audit-exporter
audit-exporter --help
```

## Build

Build the binary from source with the following command:

```Shell
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/audit-exporter ./cmd/audit-exporter/
```

## Usage

TBD

## Examples

TBD

## Contributors

Special thanks goes to all [contributors](https://github.com/thegeeklab/audit-exporter/graphs/contributors). If you would like to contribute,
please see the [instructions](https://github.com/thegeeklab/audit-exporter/blob/main/CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/thegeeklab/audit-exporter/blob/main/LICENSE) file for details.
