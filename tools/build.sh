#!/bin/bash

# Build a Bazel target targeting the Raspberry Pi

TARGET="$1"

bazel build -c opt --platforms=@io_bazel_rules_go//go/toolchain:linux_arm "$TARGET"
