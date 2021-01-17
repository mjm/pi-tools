job "named" {
  datacenters = ["dc1"]

  type     = "system"
  priority = 90

  group "named" {
    network {
      mode = "host"
      port "dns" {
        static = 53
        to     = 53
      }
    }

    service {
      name = "named"
      port = "dns"
    }

    task "named" {
      driver = "docker"
      config {
        image        = "eafxx/bind@sha256:9c15e971a7a358a4ba248e02154b7d5a6b37803bdf65371325364f3cbae9dd43"
        ports        = ["dns"]
        network_mode = "host"

        logging {
          type = "journald"
          config {
            tag = "named"
          }
        }
      }

      env {
        WEBMIN_ENABLED = "false"
        DATA_DIR       = "${NOMAD_TASK_DIR}"
      }

      resources {
        cpu    = 50
        memory = 200
      }

      template {
        destination = "local/bind/etc/named.conf"
        data        = <<EOF
include "/etc/bind/named.conf.options";
include "/etc/bind/named.conf.local";
EOF
      }

      template {
        destination = "local/bind/etc/named.conf.options"
        data        = <<EOF
acl tailscale {
  100.64.0.0/10; // range of tailscale IPs
};
acl local {
  10.0.0.0/8;    // everything inside the cluster, and in the local network;
};

options {
  directory "/var/cache/bind";
  recursion yes;
  forwarders {
    {{ range service "dns.pihole" }}
    {{ .Address }} port {{ .Port }};
    {{ else }}
    8.8.8.8;
    8.8.4.4;
    {{ end }}
  };

  auth-nxdomain no;

  dnssec-enable no;
  dnssec-validation no;

  listen-on-v6 { any; };

  zero-no-soa-ttl-cache yes;
  max-ncache-ttl 15;
};

// let the bind-exporter get stats
statistics-channels {
  inet 127.0.0.1 port 8053 allow { 127.0.0.1; };
};
EOF
      }

      template {
        destination = "local/bind/etc/named.conf.local"
        data        = <<EOF
view tailscale {
  match-clients { tailscale; };
  allow-recursion { any; };

  zone "homelab" {
    type master;
    file "/etc/bind/db.homelab.tailscale";
    zone-statistics yes;
  };

  include "/etc/bind/named.conf.default-zones";
};
view local {
  match-clients { local; };
  allow-recursion { any; };

  zone "homelab" {
    type master;
    file "/etc/bind/db.homelab.local";
    zone-statistics yes;
  };

  include "/etc/bind/named.conf.default-zones";
};
EOF
      }

      template {
        destination = "local/bind/etc/named.conf.default-zones"
        data        = <<EOF
zone "." {
  type hint;
  file "/usr/share/dns/root.hints";
};

zone "localhost" {
  type master;
  file "/etc/bind/db.local";
};

zone "127.in-addr.arpa" {
  type master;
  file "/etc/bind/db.127";
};

zone "0.in-addr.arpa" {
  type master;
  file "/etc/bind/db.0";
};

zone "255.in-addr.arpa" {
  type master;
  file "/etc/bind/db.255";
};

zone "consul" IN {
  type forward;
  forward only;
  forwarders { 127.0.0.1 port 8600; };
};
EOF
      }

      template {
        destination = "local/bind/etc/db.local"
        data        = <<EOF
;
; BIND data file for local loopback interface
;
$TTL	604800
@	IN	SOA	localhost. root.localhost. (
                  2		; Serial
             604800		; Refresh
              86400		; Retry
            2419200		; Expire
             604800 )	; Negative Cache TTL
;
@	IN	NS	localhost.
@	IN	A	127.0.0.1
@	IN	AAAA	::1
EOF
      }

      template {
        destination = "local/bind/etc/db.127"
        data        = <<EOF
;
; BIND reverse data file for local loopback interface
;
$TTL	604800
@	IN	SOA	localhost. root.localhost. (
                  1		; Serial
             604800		; Refresh
              86400		; Retry
            2419200		; Expire
             604800 )	; Negative Cache TTL
;
@	IN	NS	localhost.
1.0.0	IN	PTR	localhost.
EOF
      }

      template {
        destination = "local/bind/etc/db.0"
        data        = <<EOF
;
; BIND reverse data file for broadcast zone
;
$TTL	604800
@	IN	SOA	localhost. root.localhost. (
                  1		; Serial
             604800		; Refresh
              86400		; Retry
            2419200		; Expire
             604800 )	; Negative Cache TTL
;
@	IN	NS	localhost.
EOF
      }

      template {
        destination = "local/bind/etc/db.255"
        data        = <<EOF
;
; BIND reverse data file for broadcast zone
;
$TTL	604800
@	IN	SOA	localhost. root.localhost. (
                  1		; Serial
             604800		; Refresh
              86400		; Retry
            2419200		; Expire
             604800 )	; Negative Cache TTL
;
@	IN	NS	localhost.
EOF
      }

      template {
        destination = "local/bind/etc/db.homelab.local"
        data        = <<EOF
; BIND data file for homelab zone, when served from a local LAN IP
$TTL  1m
@   IN  SOA localhost. matt.mattmoriarity.com. (
                  1
                 1m     ; Refresh
                 1h		; Retry
                 1w		; Expire
                 1h )	; Negative Cache TTL
@   IN  NS  localhost.
;{{ range nodes }}
{{ .Node }} IN  A {{ .Address }}{{ end }}
*.homelab.  IN  CNAME {{ env "node.unique.name" }}.homelab.
unifi	IN  A 10.0.0.1
nas	IN  A 10.0.0.10
EOF
      }

      template {
        destination = "local/bind/etc/db.homelab.tailscale"
        data        = <<EOF
; BIND data file for homelab zone, when served from a Tailscale IP
$TTL  1m
@   IN  SOA localhost. matt.mattmoriarity.com. (
                  1
                 1m     ; Refresh
                 1h		; Retry
                 1w		; Expire
                 1h )	; Negative Cache TTL
@   IN  NS  localhost.
;{{ range nodes }}
{{ .Node }} IN  A {{ .Meta.tailscale_ip }}{{ end }}
*.homelab.  IN  CNAME {{ env "node.unique.name" }}.homelab.
unifi	IN  A 10.0.0.1
nas	IN  A 10.0.0.10
EOF
      }
    }
  }
}
