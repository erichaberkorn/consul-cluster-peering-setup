#!/usr/bin/env bash

set -euo pipefail

# killall server envoy server

rm -rf /tmp/consul-dc1
rm -rf /tmp/consul-dc2
rm -rf /tmp/consul-dc3

# Start Consul
consul agent -dev -config-file config_dc1.hcl 1>/dev/null &
consul agent -dev -config-file config_dc2.hcl 1>/dev/null &
consul agent -dev -config-file config_dc3.hcl 1>/dev/null &

sleep 2

CONSUL_HTTP_ADDR=localhost:7500 consul partition create -name ap1
