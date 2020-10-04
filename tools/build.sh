#!/bin/bash

# Build a Bazel target on the Raspberry Pi

TARGET="$1"

# Workaround something broken in Bazel/rules_cc/rules_go's handling of these flags
export BAZEL_LINKOPTS="-lstdc++:-lm:-latomic"
bazel build -c opt "$TARGET"
