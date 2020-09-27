#!/bin/bash

# Build and copy a binary to the Raspberry Pi

PI_USER=pi
PI_HOST=10.0.0.2

TARGET="$1"

COMMAND_NAME=$(basename "$TARGET")
BUILT_PATH="bazel-bin$TARGET/${COMMAND_NAME}_/$COMMAND_NAME"
DEST_PATH="/deploy/pi-tools/$COMMAND_NAME"

./tools/build.sh "$TARGET"
file "$BUILT_PATH"

chmod u+w "$BUILT_PATH"
scp "$BUILT_PATH" $PI_USER@$PI_HOST:$DEST_PATH

if [[ "$COMMAND_NAME" = "detect-presence" ]]; then
  ssh $PI_USER@$PI_HOST sudo setcap cap_net_raw+ep "$DEST_PATH"
fi
