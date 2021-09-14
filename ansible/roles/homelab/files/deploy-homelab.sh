#!/bin/sh

if ! [ -f /shared/homelab/homelab.tar.gz ]; then
  exit 0
fi

cd /opt/homelab
tar xzvf /shared/homelab/homelab.tar.gz
bin/homelab restart
rm /shared/homelab/homelab.tar.gz
