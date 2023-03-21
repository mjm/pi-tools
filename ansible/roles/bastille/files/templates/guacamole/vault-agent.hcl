template {
  contents = <<EOF
http-auth-header: X-Auth-Request-User

{{ with secret "database/creds/guacamole" }}
postgresql-hostname: postgresql.service.consul
postgresql-database: guacamole
postgresql-username: {{ .Data.username }}
postgresql-password: {{ .Data.password }}
{{ end }}
EOF
  destination = "/usr/local/etc/guacamole-client/guacamole.properties"
  command = "service tomcat9 restart"
}
