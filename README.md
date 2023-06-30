# gogin <!-- omit from toc -->

## Table of Contents <!-- omit from toc -->

- [Overview](#overview)
- [Installation](#installation)
  - [Create docker config for private registry access](#create-docker-config-for-private-registry-access)
  - [Add private GitLab Helm registry](#add-private-gitlab-helm-registry)
  - [Install gogin using Helm chart](#install-gogin-using-helm-chart)
- [Troubleshooting Helm chart](#troubleshooting-helm-chart)
  - [Get rendered Helm chart template](#get-rendered-helm-chart-template)
  - [Use local version of Helm chart to deploy gogin](#use-local-version-of-helm-chart-to-deploy-gogin)

## Overview

This is an end-to-end example of a Golang project that shows the complete GitLab CI pipeline that includes: 
- Compiling swagger API documentation
- Building executable
- Executing unit tests
- Running race detection
- Generating code coverage report
- Bundling docker container and publishing it to private registry
- Bundling helm chart and publishing it to private registry
- Publishing GitLab pages with code coverage reports

This GitLab Ci supports tagging docker images based on the git tags. For example the git tag `v1.2.0` will result in the docker container with tag `1.2.0` pushed to the private container registry.

Also this example demonstrates the typical Golang project structure that covers:
- HTTP Rest API handlers
- Middleware handling
- Cobra and Viper configuration
- Swagger documentation
- Unit tests and coverage reports
- Production-ready Helm chart

## Installation

This section explains how to install gogin using the private GitLab container registry and the private GitLab Helm chart registry:
- Helm Chart registry: https://git.lothric.net/api/v4/projects/examples%2Fgo%2Fgogin/packages/helm/stable
- Docker container registry: https://registry.lothric.net

### Create docker config for private registry access

```bash
export REGISTRY_USERNAME=registryUsername
export REGISTRY_PASSWORD=registryPassword 

# Create base64 encoded docker auth config
{
cat << EOF
{
    "auths": {
        "https://registry.lothric.net":{
            "auth":"`echo -n "${REGISTRY_USERNAME}:${REGISTRY_PASSWORD}" | base64`"
        }
    }
}
EOF
} | base64
```

Use the output as `image.registry.dockerConfig` value. This will allow pulling docker image from the GitLab container registry using the gogin service account.

### Add private GitLab Helm registry

```sh
# GitLab: Settings -> Repository -> Deploy tokens
export GITLAB_DEPLOYTOKEN_USERNAME=deploymentUsername
export GITLAB_DEPLOYTOKEN_SECRET=deploymentSecretKey

helm repo add \
    --username ${GITLAB_DEPLOYTOKEN_USERNAME} \
    --password ${GITLAB_DEPLOYTOKEN_SECRET} \
    gogin-repo \
    https://git.lothric.net/api/v4/projects/examples%2Fgo%2Fgogin/packages/helm/stable

helm repo update
```

### Install gogin using Helm chart

```sh
export DOCKER_CONFIG=ewogICAgImF1dGhzIjogewog...
export GOGIN_DOMAIN=gogin.k8s.lothric.net

helm upgrade \
    --install gogin \
    gogin-repo/gogin \
    --version 0.1.0-alpha \
    --namespace=services \
    --create-namespace \
    --set="image.registry.dockerConfig=${DOCKER_CONFIG}" \
    --set="ingress.host.goginDomain=${GOGIN_DOMAIN}"

helm uninstall \
    --namespace=services \
    gogin
```

## Troubleshooting Helm chart

This is how you can troubleshoot the helm chart and play with it locally.

### Get rendered Helm chart template

```sh
helm template gogin ./helm
```

### Use local version of Helm chart to deploy gogin

```sh
export DOCKER_CONFIG=ewogICAgImF1dGhzIjogewog...
export GOGIN_DOMAIN=gogin.k8s.lothric.net

helm upgrade \
    --install gogin \
    ./helm \
    --namespace=services \
    --create-namespace \
    --set="image.registry.dockerConfig=${DOCKER_CONFIG}" \
    --set="ingress.host.goginDomain=${GOGIN_DOMAIN}"

helm uninstall \
    --namespace=services \
    gogin
```