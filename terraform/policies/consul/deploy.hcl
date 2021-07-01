key "service/deploy/leader" {
  policy = "write"
}

session_prefix "" {
  policy = "write"
}

service_prefix "" {
  policy     = "write"
  intentions = "write"
}
