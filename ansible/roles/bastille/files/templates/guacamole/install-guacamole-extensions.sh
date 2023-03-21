#!/bin/sh

set -e

version="1.5.0"
extensions="quickconnect header"

needs_install() {
  extension=$1
  if [ -f /usr/local/etc/guacamole-client/extensions/guacamole-auth-$extension-$version.jar ]; then
    return 1
  else
    return 0
  fi
}

fetch_extension() {
  extension=$1
  fetch https://apache.org/dyn/closer.lua/guacamole/1.5.0/binary/guacamole-auth-$extension-$version.tar.gz\?action=download -o - | tar -C /tmp -x
}

copy_extension() {
  extension=$1
  cp /tmp/guacamole-auth-$extension-$version/guacamole-auth-$extension-$version.jar /usr/local/etc/guacamole-client/extensions/
}

mkdir -p /usr/local/etc/guacamole-client/extensions /usr/local/etc/guacamole-client/lib
ln -sf /usr/local/share/java/classes/postgresql.jar /usr/local/etc/guacamole-client/lib/postgresql.jar

for extension in $extensions; do
  if needs_install $extension; then
    fetch_extension $extension
    copy_extension $extension
  fi
done

if needs_install jdbc-postgresql; then
  fetch_extension jdbc
  cp /tmp/guacamole-auth-jdbc-$version/postgresql/guacamole-auth-jdbc-postgresql-$version.jar /usr/local/etc/guacamole-client/extensions/
fi
