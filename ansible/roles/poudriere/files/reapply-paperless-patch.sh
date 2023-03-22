#!/bin/sh

cd /usr/local/poudriere/ports/current
git fetch origin
git reset --hard origin/master
git clean -df
arc --conduit-uri https://code.home.mattmoriarity.com export --revision D1 --git | git apply -C1 -
git add -A
git commit -m "Add py-paperless-ng port"
