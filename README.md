# GoGin GitLab CI <!-- omit from toc -->

## Table of Contents <!-- omit from toc -->

- [Overview](#overview)
- [Installation](#installation)
  - [Create docker config for the private registry access](#create-docker-config-for-the-private-registry-access)
  - [Add private GitLab Helm registry](#add-private-gitlab-helm-registry)
  - [Install gogin using Helm chart](#install-gogin-using-helm-chart)
- [Troubleshooting Helm chart](#troubleshooting-helm-chart)
  - [Get rendered Helm chart template](#get-rendered-helm-chart-template)
  - [Use local version of Helm chart to deploy gogin](#use-local-version-of-helm-chart-to-deploy-gogin)

## Overview

This is an end-to-end example of a simple Golang project that shows the complete GitLab CI pipeline that includes: 
- Compiling swagger API documentation
- Building executable
- Executing unit tests
- Running race detection
- Generating code coverage report
- Bundling docker image and publishing it to the private GitLab Container Registry
- Bundling helm chart and publishing it to the private GitLab Package Registry
- Publishing code coverage reports to GitLab Pages

This GitLab CI pipeline supports tagging docker images based on the git tags.
For example the git tag `v1.2.0` will result in the docker container with tag `1.2.0` pushed to the private GitLab Container Registry.

Also this example demonstrates the typical Golang project structure that covers:
- HTTP Rest API request handling
- HTTP Rest API middleware handling
- [Cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) setup
- Swagger documentation
- Examples of unit tests and coverage reports
- Multi-stage Dockerfile
- Production-ready Helm chart

## Installation

This section explains how to install gogin using the private GitLab Container Registry and the private GitLab Helm chart registry.
This Readme refers to the two private resources that I use internally for my home infrastructure. These two resources are not publicly available and depend on private GitLab auth.
- Helm Chart registry: https://git.lothric.net/api/v4/projects/examples%2Fgo%2Fgogin/packages/helm/stable
- Docker container registry: https://registry.lothric.net

Make sure to replace the references to these two private resources with your own links to GitLab Package Registry and GitLab Container Registry.

### Create docker config for the private registry access

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

Make sure ingress is correctly configured for your kubernetes and you have sub-domain routing.

```sh
export DOCKER_CONFIG=ewogICAgImF1dGhzIjogewog...
export GOGIN_DOMAIN=gogin.k8s.lothric.net

helm upgrade \
    --install gogin \
    gogin-repo/gogin \
    --version 0.2.0 \
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