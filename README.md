# Distributed Container Launcher

## Overview
Development of a pull-based system for launching containers across multiple hosts
with centralized control via a master node. Configuration is defined using YAML manifests that describe the containers and 
assign them to specific hosts.
## How to run 
### on master node:
```
go run ./cmd/master [flags] # Use -h to see available flags
```
### on slave nodes:
```
go run ./cmd/slave [flags] # Use -h to see available flags
```
### on client(cli):
```
go run ./cmd/cli [flags] # Use -h to see available flags
```

## How to run tests:
```
make test
```

## Design Document
https://docs.google.com/document/d/1FeeSc4tqoPcfpIUSRdpBGDKMQZmfi9AWYnqN7t8_m7w/edit?tab=t.0
