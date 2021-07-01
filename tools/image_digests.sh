#!/bin/bash

dst_file="$1"
shift

: > "$dst_file"

while (( "$#" )); do
  repo_name="$1"
  digest_file="$2"
  digest=$(cat "$digest_file")

  echo "${repo_name}@${digest}" >> "$dst_file"

  shift
  shift
done
