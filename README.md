# Gogin


## Generate private registry dockerconfigjson for k8s

```sh
touch registry_auth
```

```sh
echo -n "{REGISTRY_USERNAME}:{REGISTRY_PASSWORD}" | base64
```

```json
{
    "auths": {
        "https://registry.domain.com":{
            "auth":"Z2l0XXXXXXXXXbGFiK29MYkJXXXXXXXXXXhdUJ6Z3Vi"
        }
    }
}
```

```sh
cat registry_auth | base64
```

## Installation

```sh

helm install \
	--namespace=services \
	--create-namespace \
    --set="image.repository=${IMAGE_REPOSITORY}" \
    --set="image.registry.dockerconfigjson=${DOCKER_CONFIG_BASE64}" \
    --set="ingress.host.goginDomain=${GOGIN_DOMAIN}" \
	gogin helm

```