#!/bin/sh

set -e

service paperless-ng-webserver restart
service paperless-ng-consumer restart
service paperless-ng-scheduler restart
