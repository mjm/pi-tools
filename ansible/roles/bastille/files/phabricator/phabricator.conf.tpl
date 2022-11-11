{
  "phd.taskmasters": 10,
  "cluster.mailers": [
    {
      "key": "smtp-fastmail",
      "type": "smtp",
      "options": {
        "host": "smtp.fastmail.com",
        "port": 587,
        "user": "matt@mattmoriarity.com",
        "password": {{ with secret "kv/phabricator" }}{{ .Data.data.fastmail_password | toJSON }}{{ end }},
        "protocol": "tls",
        "message-id": true
      }
    }
  ],
  "amazon-s3.endpoint": "minio.home.mattmoriarity.com",
  "storage.s3.bucket": "phabricator-files",
  "amazon-s3.region": "us-east-1",
  "amazon-s3.secret-key": {{ with secret "kv/phabricator" }}{{ .Data.data.minio_secret_key | toJSON }}{{ end }},
  "amazon-s3.access-key": "phabricator",
  "phabricator.base-uri": "https://code.home.mattmoriarity.com",
  "phabricator.timezone": "UTC",
  {{ with secret "database/creds/phabricator" -}}
  "mysql.pass": {{ .Data.password | toJSON }},
  "mysql.user": {{ .Data.username | toJSON }},
  {{ end -}}
  "mysql.host": "mysql.service.consul"
}
