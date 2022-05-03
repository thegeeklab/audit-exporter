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

make build
```

## Usage

TBD

```Text
 HELP trivy_vulnerabilities Vulnerabilities detected by trivy
# TYPE trivy_vulnerabilities gauge
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="1.0.2g-1ubuntu4.19",pkgName="libssl1.0.0",severity="LOW",vulnerabilityId="CVE-2021-3601"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="1.0.2g-1ubuntu4.19",pkgName="openssl",severity="LOW",vulnerabilityId="CVE-2021-3601"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="libc-bin",severity="LOW",vulnerabilityId="CVE-2021-33574"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="libc-bin",severity="MEDIUM",vulnerabilityId="CVE-2021-35942"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="libc-bin",severity="MEDIUM",vulnerabilityId="CVE-2021-38604"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="libc6",severity="LOW",vulnerabilityId="CVE-2021-33574"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="libc6",severity="MEDIUM",vulnerabilityId="CVE-2021-35942"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="libc6",severity="MEDIUM",vulnerabilityId="CVE-2021-38604"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="multiarch-support",severity="LOW",vulnerabilityId="CVE-2021-33574"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="multiarch-support",severity="MEDIUM",vulnerabilityId="CVE-2021-35942"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="",installedVersion="2.23-0ubuntu11.2",pkgName="multiarch-support",severity="MEDIUM",vulnerabilityId="CVE-2021-38604"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="2.23-0ubuntu11.3",installedVersion="2.23-0ubuntu11.2",pkgName="libc-bin",severity="LOW",vulnerabilityId="CVE-2009-5155"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="2.23-0ubuntu11.3",installedVersion="2.23-0ubuntu11.2",pkgName="libc-bin",severity="LOW",vulnerabilityId="CVE-2020-6096"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="2.23-0ubuntu11.3",installedVersion="2.23-0ubuntu11.2",pkgName="libc6",severity="LOW",vulnerabilityId="CVE-2009-5155"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="2.23-0ubuntu11.3",installedVersion="2.23-0ubuntu11.2",pkgName="libc6",severity="LOW",vulnerabilityId="CVE-2020-6096"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="2.23-0ubuntu11.3",installedVersion="2.23-0ubuntu11.2",pkgName="multiarch-support",severity="LOW",vulnerabilityId="CVE-2009-5155"} 1
trivy_vulnerabilities{artifactName="mongo:3.6",artifactType="container_image",fixedVersion="2.23-0ubuntu11.3",installedVersion="2.23-0ubuntu11.2",pkgName="multiarch-support",severity="LOW",vulnerabilityId="CVE-2020-6096"} 1
# HELP trivy_vulnerabilities_sum Vulnerabilities detected by trivy
# TYPE trivy_vulnerabilities_sum gauge
trivy_vulnerabilities_sum{artifactName="mongo:3.6",artifactType="container_image",severity="LOW"} 11
trivy_vulnerabilities_sum{artifactName="mongo:3.6",artifactType="container_image",severity="MEDIUM"} 6
```

## Examples

TBD

## Contributors

Special thanks goes to all [contributors](https://github.com/thegeeklab/audit-exporter/graphs/contributors). If you would like to contribute, please see the [instructions](https://github.com/thegeeklab/audit-exporter/blob/main/CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/thegeeklab/audit-exporter/blob/main/LICENSE) file for details.
