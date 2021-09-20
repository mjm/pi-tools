#!/bin/sh

# Load secrets into the environment
. /opt/.env.sh

cd /opt/homelab
bin/homelab stop >/dev/null 2>&1 || echo "Homelab app not already running"

while bin/homelab pid >/dev/null 2>&1; do
        echo "Waiting for Homelab app to shutdown"
        sleep 1
done

bin/homelab daemon
