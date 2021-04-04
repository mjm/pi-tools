#!/bin/bash

set -e

function join { local IFS="$1"; shift; echo "$*"; }

cd "$(git rev-parse --show-toplevel)"

out_path="$(mktemp -d)"
function finish {
  rm -rf "$out_path"
}
trap finish EXIT

jobs_pattern=".*/($(join '|' "$@"))\\.nomad"

find -E "$PWD/jobs" -regex "$jobs_pattern" -exec \
  bazel run //deploy/cmd/job-resolver -- -root "$PWD/jobs" -out "$out_path" '{}' +

bazel run //deploy/cmd/deploy -- -dry-run $out_path/*.json
