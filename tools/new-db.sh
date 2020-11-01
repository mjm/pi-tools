#!/bin/bash

namespace="$1"
secret_name="$2"
db_name="$3"

# first create the dev db locally
createdb "${db_name}_dev"

# now create a random password for the prod db user
db_password=$(dd if=/dev/urandom of=/dev/stdout count=1 bs=32 | base64)

# and add it as a secret in the cluster
if ! kubectl -n "$namespace" create secret generic "$secret_name" "--from-literal=db-password=$db_password"; then
  # assume that if the command fails, it's because the secret already exists, so patch it to have this new value
  kubectl -n "$namespace" patch secret "$secret_name" -p="{\"stringData\":{\"db-password\":\"$db_password\"}}"
fi

# then set up the db and user on the database in the cluster
kubectl -n storage exec postgresql-0 -i -- psql -U postgres <<SQL
create database $db_name;
create user $db_name with encrypted password '$db_password';
grant all privileges on database $db_name to $db_name;
SQL
