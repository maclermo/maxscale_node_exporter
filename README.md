# MaxScale Prom Exporter

## Maxscale Prom Exporter written in Go-lang

### Kubernetes-friendly

This Prom Exporter exports the following:

1. Servers stats
1. Services stats

You have to create your own ``creds.json`` with the following structure :

```json
{
    "username": "admin",
    "password": "mariadb",
    "host": "http://127.0.0.1",
    "port": 8989
}
```

You can create the file using a Kubernetes secret_ref.

If using as a sidecar on Kubernetes, the hostname can be ``127.0.0.1`` and you have to specify the configuration path such as :

```bash
./main /etc/node_exporter/maxscale.json
```

Otherwise, you will get the following error :

```bash
❯ go run main.go
2021/08/11 12:46:05 Usage: ./main path_of_config_file
exit status 1
```
