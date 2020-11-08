#!/bin/bash

set -o pipefail
set -e

org=mjm
repo=pi-tools
branch=main
workflow=k8s.yaml

file_to_apply="$1"

artifacts_url=$(curl -s "https://api.github.com/repos/$org/$repo/actions/workflows/$workflow/runs?branch=$branch&status=success&per_page=1" | jq -r '.workflow_runs[0].artifacts_url')
archive_url=$(curl -s "$artifacts_url" | jq -r '.artifacts[] | select(.name == "all_k8s") | .archive_download_url')

echo "Downloading archive from ${archive_url}" >&2

creds=$(cat "$HOME/.github_auth.txt")

cd "$(mktemp -d)" || exit 1
curl -sL "$archive_url" -u "$creds" -o archive.zip
unzip archive.zip

/snap/bin/microk8s kubectl apply -f "$file_to_apply.yaml"

rm *.yaml archive.zip
