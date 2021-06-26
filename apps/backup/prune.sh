#!/bin/sh

set -e

tarsnap --keyfile {{ env "NOMAD_SECRETS_DIR" }}/tarsnap.key --cachedir /var/lib/tarsnap/cache --no-default-config --list-archives \
  | prunef --keep-daily 7 --keep-weekly 4 --keep-monthly 6 daily-backup-%Y-%m-%d_%H-%M-%S \
  | xargs -r -n1 tarsnap --keyfile {{ env "NOMAD_SECRETS_DIR" }}/tarsnap.key --cachedir /var/lib/tarsnap/cache --no-default-config -v -d -f
