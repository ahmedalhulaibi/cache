# cache-api
Backend API for cache-api

Please see [the onboarding notebook](./notebooks/onboarding.ipynb) to get started

# Prerequisites

- [Go 1.24](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
 
# Running the API

To build and run the API locally, you can use the following command:

```bash
make run-local
```

By default, this will listen on two addresses, one for HTTP RESTful API `ADDR=:8080` and another for the gRPC API `GRPC_ADDR=:8090`.

You can override these defaults by setting the `ADDR` and `GRPC_ADDR` environment variables.

```bash
ADDR=:8081 GRPC_ADDR=:8091 make run-local
```

## Running the API in Docker

To build the docker image

```bash
make docker-build
```

To run the docker image

```bash
make docker-run
```

## Running in local k8s (k3d)

If you are running on a Linux machine, you can bootstrap a local k8s cluster.

```bash
make bootstrap
```

This will install various tools (k3d, kubectl, skaffold) and setup a local k8s cluster.

To deploy the application to the cluster

```bash
make run
```

This repo uses skaffold to deploy the application to the cluster. 

The cluster will expose an HTTP ingress on `localhost:8888` and the application will be available at `localhost:8888/cache-api/`.

You can override this default port by overriding the `K3D_HOST_PORT` environment variable.

# API Documentation

You can find the protobuf definitions in [`./proto/cacheapi/v1/api.proto`](./proto/cacheapi/v1/api.proto). This is the source of truth for the API.

You can find the openapi/swagger specification in [`./gen/api/swagger/cacheapi/v1/api.swagger.json`](./gen/api/swagger/cacheapi/v1/api.swagger.json). This is generated from the protobuf definitions.

## Set a key

To set a key

```http
POST http://localhost:8080/v1/set HTTP/1.1
Content-Type: application/json

{
  "bucket": "my-bucket",
  "key": "my-key",
  "value": "my-value",
  "options": {
    "ttlSeconds": 60,
    "evictionPolicy": "EVICTION_LRU"
  }
}
```

Or in curl

```bash
curl -X POST "http://localhost:8080/v1/set" -H "Content-Type: application/json" -d '{
  "bucket": "my-bucket",
  "key": "my-key",
  "value": "my-value",
  "options": {
    "ttlSeconds": 60,
    "evictionPolicy": "EVICTION_LRU"
  }
}'
```

## Get a key

To get a key

```http
GET http://localhost:8080/v1/get/my-bucket/my-key HTTP/1.1
```

Or in curl

```bash
curl -X GET "http://localhost:8080/v1/get/my-bucket/my-key" -H "accept: application/json"
```

## Get cache stats

To get cache stats

```http
GET http://localhost:8080/v1/stats HTTP/1.1
```

Or in curl

```bash
curl -X GET "http://localhost:8080/v1/stats" -H "accept: application/json"
```