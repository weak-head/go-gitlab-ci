# Go GitLab CI/CD <!-- omit from toc -->

[![go1.20](https://img.shields.io/badge/Go-1.20-00ADD8.svg)](https://go.dev/dl/)
[![gitlab](https://img.shields.io/badge/GitLab-CI/CD-fc6d27.svg)](https://docs.gitlab.com/ee/ci/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-2.0-84e92d.svg)](https://swagger.io/specification/v2/)
[![docker](https://img.shields.io/badge/Docker-24-0db7ed.svg)](https://docs.docker.com/engine/release-notes/24.0/)
[![helm3](https://img.shields.io/badge/Helm-3.12-0f1689.svg)](https://helm.sh/)

## Table of Contents <!-- omit from toc -->

- [Overview](#overview)
- [GitLab CI/CD setup](#gitlab-cicd-setup)
- [Installation](#installation)
  - [Configure GitLab Container registry](#configure-gitlab-container-registry)
  - [Configure GitLab Helm registry](#configure-gitlab-helm-registry)
  - [Deploy to Kubernetes](#deploy-to-kubernetes)
- [Building and testing](#building-and-testing)
- [Troubleshooting](#troubleshooting)
  - [Render Helm template](#render-helm-template)
  - [Deploy using local Helm template](#deploy-using-local-helm-template)
- [CLI usage](#cli-usage)
  - [Options](#options)

## Overview

This is a simple Go service with the complete end-to-end GitLab CI/CD pipeline that includes: 
- Compiling API documentation
- Building executable
- Running unit tests and race detection
- Generating code coverage report and publishing it to GitLab Pages
- Bundling docker image and publishing it to the private GitLab Container Registry
- Bundling helm chart and publishing it to the private GitLab Package Registry
- Deploying the application on tag to a configured Kubernetes cluster

This GitLab CI pipeline tags docker images based on the git tags.
For example the git tag `v1.2.0` will result in a docker container with tag `1.2.0` pushed to the private GitLab Container Registry and the application deployed to Kubernetes using `1.2.0` docker image tag.

Also this example demonstrates the typical Go project structure that includes:
- [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) setup and configuration
- HTTP REST API request and middleware handling using [Gin](https://github.com/gin-gonic/gin)
- RESTful API documentation using [swag](https://github.com/swaggo/swag) and [gin-swagger](https://github.com/swaggo/gin-swagger)
- gRPC status and health check using [grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe)
- Customized logger with support of custom fields using [logrus](https://github.com/sirupsen/logrus)
- Handling correlation id that is provided via [x-request-id](https://http.dev/x-request-id) header.
- Prometheus metrics using [client_golang](https://github.com/prometheus/client_golang)
- Simple unit tests and code coverage reports
- [Multi-stage](https://docs.docker.com/build/building/multi-stage/) Dockerfile
- Deployment-ready Helm chart with ingress, secrets, auto-scaling and service account for pulling docker images from private registry

## GitLab CI/CD setup

In order to have the deployment automation, the following environment variables should be set via `Settings -> CI/CD -> Variables` on GitLab side:
- `KUBERNETES_SERVER` - The endpoint to the Kubernetes API.
- `KUBERNETES_CERTIFICATE_AUTHORITY_DATA` - The CA configuration for the Kubernetes cluster.
- `KUBERNETES_USER_NAME` - Kubernetes user name.
- `KUBERNETES_USER_TOKEN` - Kubernetes user token.
- `GITLAB_REGISTRY` - GitLab container registry that Kubernetes should use for auth.
- `GITLAB_REGISTRY_USER_NAME` - Persistent GitLab user name to pull image from the GitLab container registry.
- `GITLAB_REGISTRY_USER_TOKEN` - Persistent GitLab user token to pull image from the GitLab container registry. 
- `APPLICATION_DOMAIN` - The FQDN the application is deployed to.
- `APPLICATION_KUBERNETES_NAMESPACE` - The Kubernetes namespace the application is deployed to.

Refer to the [gitlab-ci.yml](./.gitlab-ci.yml) for the configured GitLab CI/CD steps.

The deployment account could be crated via `Settings -> Repository -> Deploy tokens` (`GITLAB_REGISTRY_USER_NAME` and `GITLAB_REGISTRY_USER_TOKEN`). The account should have the `read_registry` scope in order to pull a docker image from the GitLab Container registry.

## Installation

This readme refers to the two private registries that I use internally for my home infrastructure.  
These two resources are not publicly available and require private GitLab authentication:  
- GitLab Package registry ( https://git.lothric.net/api/v4/projects/examples%2Fgo%2Fgogin/packages/helm/stable )
- GitLab Container registry ( https://registry.lothric.net )

Make sure to replace the references to these two private resources with your own links.

### Configure GitLab Container registry

This docker configuration is used by Kubernetes to pull a docker image from the private GitLab Container registry.

```bash
# GitLab: Settings -> Repository -> Deploy tokens
export GITLAB_CONTAINER_REGISTRY=https://registry.lothric.net
export GITLAB_CONTAINER_REGISTRY_USER_NAME=registryUserName
export GITLAB_CONTAINER_REGISTRY_USER_TOKEN=registryUserToken

# Create base64 encoded docker auth config
export DOCKER_CONFIG=$({
cat << EOF
{
    "auths": {
        "${GITLAB_CONTAINER_REGISTRY}":{
            "auth":"`echo -n "${GITLAB_CONTAINER_REGISTRY_USER_NAME}:${GITLAB_CONTAINER_REGISTRY_USER_TOKEN}" | base64 -w 0`"
        }
    }
}
EOF
} | base64 -w 0)
```

Use this value for the `image.registry.dockerConfig` in the Helm chart.

### Configure GitLab Helm registry

This is setup and connection to the private GitLab Helm chart registry.

```sh
# GitLab: Settings -> Repository -> Deploy tokens
export GITLAB_HELM_REGISTRY=https://git.lothric.net/api/v4/projects/examples%2Fgo%2Fgogin/packages/helm/stable
export GITLAB_HELM_USER_NAME=deploymentUserName
export GITLAB_HELM_USER_TOKEN=deploymentUserToken

helm repo add \
    --username ${GITLAB_HELM_USER_NAME} \
    --password ${GITLAB_HELM_USER_TOKEN} \
    gogin-repo \
    ${GITLAB_HELM_REGISTRY}

helm repo update
```

### Deploy to Kubernetes

Make sure ingress is correctly configured on your Kubernetes cluster and you have sub-domain routing rules for the specified application domain.

```sh
export DOCKER_CONFIG=ewCAogog...
export APPLICATION_DOMAIN=gogin.k8s.lothric.net

helm upgrade \
    --install gogin \
    gogin-repo/gogin \
    --version 0.2.0 \
    --namespace=services \
    --create-namespace \
    --set="image.registry.dockerConfig=${DOCKER_CONFIG}" \
    --set="ingress.host.goginDomain=${APPLICATION_DOMAIN}"

helm uninstall \
    --namespace=services \
    gogin
```

## Building and testing

Run `make help` for the list of available commands. The most useful commands are the following:
```sh
# Generate OpenAPI documentation
make swagger

# Create executable
make build

# Run unit tests
make test

# Generate code coverage report
make coverage
```

## Troubleshooting

This is how you can troubleshoot, test and verify the helm template locally.

### Render Helm template

```sh
helm template \
    gogin ./helm
```

### Deploy using local Helm template

```sh
export DOCKER_CONFIG=ewodGhwog...
export APPLICATION_DOMAIN=gogin.k8s.lothric.net

helm upgrade \
    --install gogin \
    ./helm \
    --namespace=services \
    --create-namespace \
    --set="image.registry.dockerConfig=${DOCKER_CONFIG}" \
    --set="ingress.host.goginDomain=${APPLICATION_DOMAIN}"

helm uninstall \
    --namespace=services \
    gogin
```

## CLI usage

This command could be used to start the application locally.
> `gogin [flags]`

### Options
```
      --config string                    Path to config file.
  -h, --help                             help for gogin
      --http.gin.mode string             Gin mode. (default "release")
      --http.port string                 HTTP API port. (default "8080")
      --log.formatter string             Log formatter. (default "json")
      --log.level string                 Log level. (default "info")
      --metrics.prometheus.addr string   HTTP address of prometheus metrics endpoint. (default ":8880")
      --metrics.prometheus.path string   HTTP URL endpoint of prometheus metrics endpoint. (default "/metrics")
      --node.name string                 Unique server ID.
      --status.rpc.addr string           Rpc address of status server. (default ":8400")
```