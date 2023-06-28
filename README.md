# Gogin <!-- omit from toc -->

## Table of Contents <!-- omit from toc -->

- [Overview](#overview)
- [Installation](#installation)
  - [Add GitLab Helm registry](#add-gitlab-helm-registry)
  - [Authorize private docker registry access](#authorize-private-docker-registry-access)
  - [Install/update gogin](#installupdate-gogin)
  - [Uninstall gogin](#uninstall-gogin)

## Overview

TBD

## Installation

TBD 

### Add GitLab Helm registry

```sh
# Settings -> Repository -> Deploy tokens
export GITLAB_DEPLOYTOKEN_USERNAME = "gitlab-deployment-1"
export GITLAB_DEPLOYTOKEN_SECRET = "secretKey"

helm repo add \
    --username ${GITLAB_DEPLOYTOKEN_USERNAME} \
    --password ${GITLAB_DEPLOYTOKEN_SECRET} \
    gogin-repo \
    https://git.lothric.net/api/v4/projects/examples%2Fgo%2Fgogin/packages/helm/stable

helm repo update
```

### Authorize private docker registry access

TBD

```bash
{
cat << EOF
{
    "auths": {
        "https://registry.lothric.net":{
            "auth":"`echo -n "REGISTRY_USERNAME:REGISTRY_PASSWORD" | base64`"
        }
    }
}
EOF
} | base64
```

Use the `base64` output as `image.registry.dockerconfigjson` to allow pulling docker image from GitLab container registry using the service account.

### Install/update gogin

```sh
export DOCKER_CONFIG = "dockerconfigjson"
export GOGIN_DOMAIN = "gogin.k8s.lothric.net"

# or load ENV vars from .env file
set -o allexport; source .env; set +o allexport

helm upgrade \
    --install gogin \
    gogin-repo/gogin \
    --version 0.1.0-alpha \
    --namespace=services \
    --create-namespace \
    --set="image.registry.dockerconfigjson=${DOCKER_CONFIG}" \
    --set="ingress.host.goginDomain=${GOGIN_DOMAIN}"
```

### Uninstall gogin

```sh
helm uninstall \
    --namespace=services \
    gogin
```