# Consul Cluster Peering Setup

## Steps

Cluster peering hasn't been added to the Consul Terraform provider yet, so we need to set `~/.terraformrc` to:

```hcl
provider_installation {
  dev_overrides {
    "hashicorp/consul" = "<PATH_TO_TF_PROVIDER>"
  }

  direct {}
}
```

Ensure that `consul` is an enterprise binary and `CONSUL_LICENSE` is exported.

Run `./run.sh`.

Run `cd tf`

Write the following to `input.json`:

```json
[
  {
    "address": "localhost:7500",
    "peer_name": "cluster-01"
  },
  {
    "address": "localhost:7500",
    "peer_name": "cluster-02",
    "partition": "ap1"
  },
  {
    "address": "localhost:8500",
    "peer_name": "cluster-03"
  },
  {
    "address": "localhost:9500",
    "peer_name": "cluster-04"
  }
]
```

Run `go build . && ./consul-peering-setup input.json` to generate the Terraform files.


Run `terraform init`.

Run `terraform apply`.

Run `curl localhost:<PORT>/v1/peerings` to see the peerings on each server.
