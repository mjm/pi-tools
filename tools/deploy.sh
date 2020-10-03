#!/bin/bash

# Build a tarball of all of the pi-tools and extract it on the Raspberry Pi

PI_USER=pi
PI_HOST=10.0.0.2

TARGET=":pkg"
BUILT_PATH="bazel-bin/pkg.tar"
DEST_PATH="/deploy/pi-tools/pkg.tar"

./tools/build.sh "$TARGET"
file "$BUILT_PATH"

scp "$BUILT_PATH" $PI_USER@$PI_HOST:$DEST_PATH
ssh $PI_USER@$PI_HOST "cd /deploy/pi-tools && tar xvf pkg.tar && rm -rf pkg.tar"
