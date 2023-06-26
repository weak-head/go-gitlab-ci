# Gogin

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

```yaml
registry.dockerconfigjson: "BASE64_VALUE"
```