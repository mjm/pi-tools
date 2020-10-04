#!/bin/bash

# Build a tarball of all of the pi-tools and extract it on the Raspberry Pi

PI_USER=pi
PI_HOST=10.0.0.2

ssh $PI_USER@$PI_HOST "cd /home/pi/pi-tools && git pull && ./tools/deploy-local.sh"
