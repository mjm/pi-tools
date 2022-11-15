#!/bin/sh

SCRIPTNAME=`basename "$0"`

# check for running script
STATUS=`ps ax | grep "$SCRIPTNAME" | grep -v grep | wc -l`

# compare to 2 because the ` create a sub process
if [ "$STATUS" -gt 2 ]; then
  echo "already running ... exit"
  exit 0
fi

# The build
POUDRIERE="/usr/local/bin/poudriere"
PORTLIST="/usr/local/etc/pkglist"
JAILS="amd64-13-1"
REPOS="current"
URL="https://pkg.home.mattmoriarity.com"

poudriere_build() {
    for JAIL in $JAILS; do
      for REPO in $REPOS; do
        echo "Started $JAIL / $REPO ("`/bin/date | /usr/bin/tr -d '\n'`")"
        "$POUDRIERE" bulk -j "$JAIL" -p "$REPO"  -f "$PORTLIST"
        echo "    Cleaning $REPO ("`/bin/date | /usr/bin/tr -d '\n'`")"
        "$POUDRIERE" pkgclean -j "$JAIL" -p "$REPO" -f "$PORTLIST" -y
        echo "    Finished $REPO ("`/bin/date | /usr/bin/tr -d '\n'`")"
      done
    done
}

repos_update() {
  echo "[$SCRIPTNAME] Updating ports tree..."

  for REPO in $REPOS; do
    echo "[$SCRIPTNAME] Updating ports tree... $REPO"
    "$POUDRIERE" ports -p "$REPO" -u

    if [ $? -ne 0 ]; then
      echo "    Error updating ports tree."
      exit 1
    fi

  echo "    Ports tree has been updated."
  done
}

echo "This is a log of poudriere. Details: $URL"
echo ""

repos_update
poudriere_build

# echo "[$SCRIPTNAME] Cleaning distfiles..."
# "$POUDRIERE" distclean -p "$REPOS" -f "$PORTLIST" -y > /dev/null

echo "[$SCRIPTNAME] Finished. ("`/bin/date | /usr/bin/tr -d '\n'`")"
exit 0
