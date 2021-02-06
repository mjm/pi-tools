namespace "default" {
  capabilities = ["list-jobs", "read-job", "submit-job"]
}

host_volume "*" {
  policy = "write"
}
