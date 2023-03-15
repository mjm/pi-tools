#!/bin/sh

version="1.5.0"
extensions="quickconnect header jdbc"

mkdir -p /usr/local/etc/guacamole-client/extensions

for extension in $extensions; do
  if [ ! -f /usr/local/etc/guacamole-client/extensions/guacamole-auth-$extension-$version.jar ]; then
    fetch https://apache.org/dyn/closer.lua/guacamole/1.5.0/binary/guacamole-auth-$extension-$version.tar.gz\?action=download -o - | tar -C /tmp -x
    cp /tmp/guacamole-auth-$extension-$version/guacamole-auth-$extension-$version.jar /usr/local/etc/guacamole-client/extensions/
  fi
done
