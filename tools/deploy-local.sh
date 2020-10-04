#!/bin/bash

set -xe

TARGET=":pkg"
BUILT_PATH="bazel-bin/pkg.tar"
DEST_PATH="/deploy/pi-tools"

./tools/build.sh "$TARGET"

tar -xvf "$BUILT_PATH" -C "$DEST_PATH"

sudo systemctl restart detect-presence
