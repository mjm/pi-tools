template {
  source      = "/usr/local/etc/phabricator.conf.tpl"
  destination = "/usr/local/lib/php/phabricator/conf/local/local.json"
  command     = "service phd restart"
}
