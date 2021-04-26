consul {
  address = "http://10.0.2.10:8500"
}

vault {
  address     = "http://vault.service.consul:8200"
  renew_token = true
}
