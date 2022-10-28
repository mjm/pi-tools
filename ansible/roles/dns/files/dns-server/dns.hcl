template {
  source               = "/usr/local/etc/namedb/named.forwarders.conf.tpl"
  destination          = "/usr/local/etc/namedb/named.forwarders.conf"
  command              = "service named reload"
  error_on_missing_key = true
}
