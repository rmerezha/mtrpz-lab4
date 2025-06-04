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
go run ./cmd/cli [subcomands] [flags]
```

## How to run tests:
```
make test
```

## Design Document
https://docs.google.com/document/d/1FeeSc4tqoPcfpIUSRdpBGDKMQZmfi9AWYnqN7t8_m7w/edit?tab=t.0

## Master Node

The master node is responsible for centralized orchestration and coordination. It exposes a set of HTTP API endpoints used by slave nodes and the CLI client. These endpoints include:

- POST /api/v1/state – Update the state of a container (e.g., stop, restart). (for slave node)
- GET /api/v1/container – Retrieve a list of containers running on a specific host. (for slave node)
- POST /api/v1/container/action – Apply a container action (stop, kill, restart, remove).
- POST /api/v1/manifest/up – Register a new manifest (YAML file with container configuration).
- POST /api/v1/manifest/down – Mark a manifest for removal.
- POST /api/v1/manifest/ps – List containers defined by a specific manifest.
- POST /api/v1/token – Generate a new authentication token.

All endpoints except /api/v1/token require a valid Bearer token provided via the Authorization header.

The master node maintains internal state using a Planner, ensuring all updates to manifests and container states are consistent and thread-safe.

## Slave Node

The slave node includes two pull-based listeners for communication with the master node:
* PollingListener — periodically pulls the list of containers assigned to the current host:

  - Makes a GET /api/v1/container?host=... request.

  - Applies the desired state (e.g., run, stop, remove) using the local Runner implementation.

* StateWatcherListener — periodically checks the actual state of running containers and reports any changes:

  - Sends updates via POST /api/v1/state.

  - Reports to the master only if the container state has changed.

Both listeners run at a configured interval in parallel and use a token for authentication.

## Client (CLI)

The CLI client provides a command-line interface for interacting with the master node’s API. It supports three main command groups:

* manifest — manage manifests describing container deployments:

  - up — upload a YAML manifest to the master.

  - down — remove a manifest by name.

  - ps — list containers from a manifest.
  
    *Flags: -f for manifest file, --url for master API base URL, --token for authentication token.*

* container — control individual containers on hosts:

  - Subcommands: stop, kill, restart, rm.
  - Flags: -h for host, -c for container name, --url and --token for authentication.

* token generate — generate an access token by providing a password.

Each command constructs and sends HTTP requests with proper authorization headers to the master node, handles responses, and outputs the result or errors.
