#!/bin/bash

set -o pipefail
set -e

bazel run --platforms @io_bazel_rules_go//go/toolchain:linux_arm64 "$@"
