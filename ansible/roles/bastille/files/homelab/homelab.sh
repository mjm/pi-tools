#!/bin/sh

# PROVIDE: homelab
# REQUIRE: DAEMON
# KEYWORD: shutdown
#
# Add the following lines to /etc/rc.conf.local or /etc/rc.conf
# to enable this service:
#
# homelab_enable (bool):	Set it to YES to enable Homelab.
#     Default is "NO".
# homelab_user (user):	Set user to run Homelab.
#     Default is "root".
# homelab_group (group):	Set group to run Homelab.
#     Default is "wheel".

. /etc/rc.subr

name=homelab
rcvar=homelab_enable

load_rc_config $name

: ${homelab_enable:="NO"}
: ${homelab_user:="root"}
: ${homelab_group:="wheel"}

procname="/usr/local/homelab/current/bin/homelab"

start_cmd=homelab_start
stop_cmd=homelab_stop
deploy_cmd=homelab_deploy
extra_commands=deploy

homelab_start()
{
  . /usr/local/homelab/.env.sh
  ${procname} daemon
}

homelab_stop()
{
  . /usr/local/homelab/.env.sh
  ${procname} stop >/dev/null 2>&1 || echo "Homelab app not already running"

  while ${procname} pid >/dev/null 2>&1; do
    echo "Waiting for Homelab app to shutdown"
    sleep 1
  done
}

homelab_deploy()
{
  if ! [ -f /var/db/homelab/homelab.tar.gz ]; then
    echo "No homelab tarball found to deploy."
    exit 0
  fi

  cd /usr/local/homelab/current
  tar xzvf /var/db/homelab/homelab.tar.gz
  run_rc_command restart
  rm /var/db/homelab/homelab.tar.gz
}

run_rc_command "$1"
