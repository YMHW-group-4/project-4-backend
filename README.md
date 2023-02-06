![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/YMHW-group-4/project-4-backend)
![GitHub issues](https://img.shields.io/github/issues/YMHW-group-4/project-4-backend)
[![go](https://github.com/YMHW-group-4/project-4-backend/actions/workflows/build.yml/badge.svg)](https://github.com/YMHW-group-4/project-4-backend/actions/workflows/build.yml)

# Project-4
This is the repository for the cryptonode of the project of year 4. The project is about creating a cryptocurrency. 
This project is created by Group 4.

## Requirements
- Docker
- Go 1.19

## Configuration

Configuring the node can be done by using enviroment variables:

* `"DEBUG", "false"` Sets log level to debug.
* `"PORT", "30333"` Sets the port for the node.
* `"API_PORT", "8080"` Sets the API port.
* `"DNS_SEED", "localhost:3000"` Sets the address of the DNS seed.
* `"INTERVAL", "20m"` Sets the interval of the scheduler.

To set multiple enviroments variables on a local machine (when not using a supervisor, or docker)
a file that specifies all the enviroment variables can be made. For example a file `node.env` can be created, 
and within this file multiple enviroment variables can be set.

```env
DEBUG=true
INTERVAL=20m15s
```

The variables can be set by executing the following command:
```shell
$ export $(cat node.env | xargs)
```

When the node is launched, it will use the enviroment variables that have been set.

## Development

### Cloning

```shell
git clone git@github.com:YMHW-group-4/project-4-backend.git
cd project-4-backend/
```

### Install dependencies

```
$ make deps
```

### Building

The project can be build for multiple platforms:

```
$ make build_amd64    # Linux AMD64 binary
$ make build_arm64v8  # Linux arm (armv8) binary
$ make build_windows  # Windows binary
```

To build for the current platform use:

```
$ make build
```

To build for all platforms use:

```
$ make build_all
```

### Unit test

Unit tests can be run by using the command:

```
$ make test
```

### Linter

Various linters can be run to check the quality of the code.
```
$ make lint
```

The code can be formatted by using:
```
$ make format
```

### Docker

A docker container can be made by either running the Dockerfile, or using the following commands:
```
$ make docker_amd64     # Build amd64 docker container
$ make docker_arm64v8   # Build arm64v8 docker container
$ make docker_all       # Builds all containers
```



