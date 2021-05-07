template {
  source      = "/usr/local/www/paperless-ng/paperless.conf.tpl"
  destination = "/usr/local/www/paperless-ng/paperless.conf"
  command     = "/usr/local/bin/paperless-restart"
}
