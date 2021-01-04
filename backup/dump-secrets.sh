#!/bin/sh

SECRET_PATH=/backup/secrets
SECRETS=$(cat <<END
auth/guard-pki
auth/oauth2-proxy
backup/postgresql
backup/tarsnap
deploy/deploy
detect-presence/detect-presence
go-links/go-links
homebase/homebase-bot-srv
monitoring/alertmanager
monitoring/grafana
monitoring/unifi-exporter
pi-hole/pi-hole
security/ca-key-pair
storage/postgresql
END
)

for secret in $SECRETS; do
  secret_dest="$SECRET_PATH/$secret"
  mkdir -p "$(dirname "$secret_dest")"
  namespace=$(echo $secret | awk -F'/' '{ print $1 }')
  name=$(echo $secret | awk -F'/' '{ print $2 }')
  kubectl get secret -o yaml -n "$namespace" "$name" > $secret_dest
  echo "Wrote secret $secret to $secret_dest"
done
