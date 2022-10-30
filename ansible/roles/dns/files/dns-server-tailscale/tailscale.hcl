template {
  source               = "/usr/local/etc/namedb/master/homelab.db.tpl"
  destination          = "/usr/local/etc/namedb/master/homelab.db"
  command              = "service named reload"
  error_on_missing_key = true
  left_delimiter       = "<<"
  right_delimiter      = ">>"
}

template {
  source               = "/usr/local/etc/namedb/master/home.mattmoriarity.com.db.tpl"
  destination          = "/usr/local/etc/namedb/master/home.mattmoriarity.com.db"
  command              = "service named reload"
  error_on_missing_key = true
  left_delimiter       = "<<"
  right_delimiter      = ">>"
}
