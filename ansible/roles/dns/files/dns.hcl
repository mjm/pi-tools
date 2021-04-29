consul {
  address = "10.0.2.10:8500"
}

template {
  source               = "/usr/local/etc/namedb/named.forwarders.conf.tpl"
  destination          = "/usr/local/etc/namedb/named.forwarders.conf"
  command              = "service named reload"
  error_on_missing_key = true
}

template {
  source               = "/usr/local/etc/namedb/master/homelab.local.db.tpl"
  destination          = "/usr/local/etc/namedb/master/homelab.local.db"
  command              = "service named reload"
  error_on_missing_key = true
  left_delimiter       = "<<"
  right_delimiter      = ">>"
}

template {
  source               = "/usr/local/etc/namedb/master/home.mattmoriarity.com.local.db.tpl"
  destination          = "/usr/local/etc/namedb/master/home.mattmoriarity.com.local.db"
  command              = "service named reload"
  error_on_missing_key = true
  left_delimiter       = "<<"
  right_delimiter      = ">>"
}

template {
  source               = "/usr/local/etc/namedb/master/homelab.tailscale.db.tpl"
  destination          = "/usr/local/etc/namedb/master/homelab.tailscale.db"
  command              = "service named reload"
  error_on_missing_key = true
  left_delimiter       = "<<"
  right_delimiter      = ">>"
}

template {
  source               = "/usr/local/etc/namedb/master/home.mattmoriarity.com.tailscale.db.tpl"
  destination          = "/usr/local/etc/namedb/master/home.mattmoriarity.com.tailscale.db"
  command              = "service named reload"
  error_on_missing_key = true
  left_delimiter       = "<<"
  right_delimiter      = ">>"
}
