#!/bin/bash

set -e

pid=$(pgrep kube-apiserver)
rss_kb=$(cat /proc/${pid}/smaps_rollup | grep Rss\: | awk '{ print $2 }')

if (( rss_kb > 2000000 )); then
  echo "Current kube-apiserver memory is ${rss_kb} kB which is more than 2 gB. Restarting kube-apiserver." >&2
  systemctl restart snap.microk8s.daemon-apiserver
fi
