template {
  source          = "/usr/local/etc/alertmanager/alertmanager.yml.tpl"
  destination     = "/usr/local/etc/alertmanager/alertmanager.yml"
  left_delimiter  = "<<"
  right_delimiter = ">>"
  command         = "service alertmanager reload"
}
