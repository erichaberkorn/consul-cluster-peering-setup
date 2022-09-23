ports {
  dns = 7600
  http = 7500
  serf_lan = 7301
  serf_wan = 7401
  server = 7300
  grpc = 7502
}
server = true
bootstrap = true
data_dir = "/tmp/consul-dc3"
bind_addr = "127.0.0.1"
datacenter = "dc3"
node_name = "secondary"
ui_config {
  enabled = true
}
connect {
  enabled = true
}

# verify_incoming = true
# verify_outgoing = true
# verify_server_hostname = true
# ca_file = "consul-agent-ca.pem"
# cert_file = "dc2-server-consul-0.pem"
# key_file = "dc2-server-consul-0-key.pem"
# auto_encrypt {
#   allow_tls = true
# }
peering {
  enabled = true
}
