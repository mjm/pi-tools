template {
  source      = "/usr/local/etc/paperless.conf.tpl"
  destination = "/usr/local/etc/paperless.conf"
  command     = "/usr/local/bin/paperless-restart"
}
