# Gogin


## Generate private registry dockerconfigjson for k8s

Create `dockerconfigjson` that is required for private container repositories auth on kubernetes:

```bash
{
cat << EOF
{
    "auths": {
        "https://registry.domain.net":{
            "auth":"`echo -n "REGISTRY_USERNAME:REGISTRY_PASSWORD" | base64`"
        }
    }
}
EOF
} | base64
```

The output is for `image.registry.dockerconfigjson` to authenticate for image pull.

## Add Helm registry

```sh
helm repo add \
    --username <username> \
    --password <personal_access_token> \
    gogin-repo \
    https://git.lothric.net/api/v4/projects/<project_id>/packages/helm/stable
```

## Installation

```sh
helm install \
	--namespace=services \
	--create-namespace \
    --set="image.registry.dockerconfigjson=${DOCKER_CONFIG_BASE64}" \
    --set="ingress.host.goginDomain=${GOGIN_DOMAIN}" \
	gogin helm
```