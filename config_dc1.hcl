ports {
  dns = 8600
  http = 8500
  serf_lan = 8301
  serf_wan = 8401
  server = 8300
  grpc = 8502
}
bind_addr = "127.0.0.1"
bootstrap = true
server = true
data_dir = "/tmp/consul-dc1"
datacenter = "dc1"
node_name = "primary"
ui_config {
  enabled = true
}
connect {
  enabled = true
}

# verify_incoming = true
# verify_outgoing = true
# verify_server_hostname = true
# verify_incoming_https = false
# verify_incoming_rpc = true
# ca_file = "consul-agent-ca.pem"
# cert_file = "dc1-server-consul-0.pem"
# key_file = "dc1-server-consul-0-key.pem"
# auto_encrypt {
#   allow_tls = true
# }
peering {
  enabled = true
}
