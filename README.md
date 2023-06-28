# Gogin


## Access private docker registry from k8s

Create `dockerconfigjson` that is required for private container repositories auth on kubernetes:

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

The output is for `image.registry.dockerconfigjson` to authenticate for image pull.

## Installation

```sh

# Settings -> Repository -> Deploy tokens
export GITLAB_DEPLOYTOKEN_USERNAME = "gitlab-deployment-1"
export GITLAB_DEPLOYTOKEN_SECRET = "secret"

# Generate for private registry
export DOCKER_CONFIG = "dockerconfigjson"

# Deployment domain for ingress
export GOGIN_DOMAIN = "gogin.k8s.lothric.net"

helm repo add \
    --username ${GITLAB_DEPLOYTOKEN_USERNAME} \
    --password ${GITLAB_DEPLOYTOKEN_SECRET} \
    gogin-repo \
    https://git.lothric.net/api/v4/projects/examples%2Fgo%2Fgogin/packages/helm/stable

helm install \
	--namespace=services \
	--create-namespace \
    --set="image.registry.dockerconfigjson=${DOCKER_CONFIG}" \
    --set="ingress.host.goginDomain=${GOGIN_DOMAIN}" \
	gogin helm
```

## Publish devel Helm chart (CI pipeline)

```sh
# Add GitLab helm repo
helm repo add  \
    --username gitlab-ci-token \
    --password ${CI_JOB_TOKEN} \
    ${CI_PROJECT_NAME} \
    ${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/helm/devel

# For cm-push
helm plugin install https://github.com/chartmuseum/helm-push.git

# Package helm chart
helm package ./helm

# Publish helm chart
helm cm-push gogin-0.1.0.tgz ${CI_PROJECT_NAME}
```