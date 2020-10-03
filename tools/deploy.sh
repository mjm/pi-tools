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

if [[ "$COMMAND_NAME" = "detect-presence-srv" ]]; then
  scp detect-presence/cmd/detect-presence-srv/detect-presence.service $PI_USER@$PI_HOST:/deploy/pi-tools/detect-presence.service
fi
