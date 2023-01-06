# CF CLI Syslog Drain Plugin

CF CLI plugin to help manage syslog drains

The plugin currently has three different commands that all list the syslog drain details in CSV format.
The output CSV has the following fields:

```csv
Org, Space, Bound App Name, Drain Name, Drain URL, Drain GUID, Drain Service Last Operation
```

## Installation

```shell
make && cf install-plugin -f ./bin/syslog-drain-plugin
```

## Usage

List all syslog log drains for an entire foundation in CSV format

```shell
$ cf list-syslog-drains
```

List all syslog log drains for the currently targeted org in CSV format

```shell
$ cf list-org-syslog-drains
```

List all syslog log drains for the currently targeted space in CSV format

```shell
$ cf list-space-syslog-drains
```
