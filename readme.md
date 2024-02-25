# Healthpose

This project provides an HTTP server that performs health checks to the specified services and exposes HTTP endpoints to get status of configured health checks. Main use-cases of this project are:

- Exposing HTTP health check endpoints for services that don't have a native HTTP health check API (i.e. databases or containers).
- Exposing HTTP health check endpoints from the private network to the public (i.e. you are using external monitoring tool/service that doesn't have access to the private network).

# Features

- Supports all health checks from [hellofresh/health-go](https://github.com/hellofresh/health-go) library except gRPC (will be enabled later), and additionally:
    - DNS hostname resolve check (A, CNAME, PTR, TXT records)
    - ICMP ping

# Getting started

[//]: # (## Releases)

[//]: # ()
[//]: # (1. Download a binary for your OS from the [GitLab releases]&#40;https://gitlab.com/buzzer13/healthpose/-/releases&#41; page.)

[//]: # (2. Prepare a [configuration file]&#40;#configuration&#41; and put it in the supported directory.)

## Container

1. Pull and run `registry.gitlab.com/buzzer13/healthpose:<tag>` image (tag can be either `latest`, or a specific version like `v1.0.1`).
    1. `docker run --volume="healthpose-config:/config" --name=healthpose -it "registry.gitlab.com/buzzer13/healthpose:latest"`
    2. To use ping health check with runtimes other than Docker engine - you may need to apply `net.ipv4.ping_group_range=0 2147483647` sysctl variable to the container.
2. Update an example [configuration file](#configuration) that was created in the volume.

# Configuration

Config file can be placed at the following paths:

- `/etc/healthpose/healthpose.yaml`
- `/config/healthpose.yaml`
- `healthpose.yaml`
- At the path, specified by the `CONFIG_FILE` environment variable.

Here is the general configuration file structure:

```yaml
# Internal HTTP server configuration dictionary.
http:
  # Address and port HTTP server should listen at
  listen: :8080

# Health check configuration dictionary.
# Key is the service name, that will be used in the heath check URL.
# If you define `test` service below - then `GET /test` endpoint will be created
# in the internal HTTP server.
services:
  test:
    # Required. Name of the service.
    name: example
    # Required. Version of the service.
    version: v1.0
    # Required. List of service health checks, that should be performed.
    checks:
      - # Required. Health check name. It will be shown in the endpoint response if check fails.
        name: cassandra
        # Required. Health check interval.
        timeout: 60
        # Optional. If check is marked as optional - it won't fail health check of the whole service.
        optional: true
        # Required. One of the available health check configuration dictionaries.
        # Examples can be found here: https://gitlab.com/buzzer13/healthpose/-/blob/master/misc/config/healthpose.yaml
        dns:
          address: example.com
          server: 8.8.8.8:53
          type: a
          request_timeout: 5
          fallback_delay: 0.3
```

You can check out [example config here](https://gitlab.com/buzzer13/healthpose/-/blob/master/misc/config/healthpose.yaml).
