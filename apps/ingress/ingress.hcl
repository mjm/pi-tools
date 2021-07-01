# Allow issuing *.homelab certificates for serving HTTPS pages
path "pki-homelab/issue/homelab" {
  capabilities = ["update"]
}

# Allow issuing client certs for accessing the Nomad API over mTLS
path "pki-int/issue/nomad-cluster" {
  capabilities = ["update"]
}
