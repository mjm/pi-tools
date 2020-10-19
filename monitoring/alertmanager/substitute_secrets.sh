#!/bin/sh

sed \
  -e "s/__PUSHOVER_USER_KEY__/$PUSHOVER_USER_KEY/" \
  -e "s/__PUSHOVER_TOKEN__/$PUSHOVER_TOKEN/" \
  /etc/alertmanager/alertmanager.yml \
  > /cfg/alertmanager/alertmanager.yml
