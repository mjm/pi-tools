#!/bin/bash

service="$1"
name="$2"

case "$service" in
detect-presence)
  migration_dir=detect-presence/database/migrate
  ;;
go-links)
  migration_dir=go-links/database/migrate
  ;;
*)
  echo "new-migration.sh: unknown service name $service" >&2
  exit 1
esac

migrate create -ext sql -dir "$migration_dir" -seq "$name"
