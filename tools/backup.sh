#!/bin/bash

kubectl -n backup create job "tarsnap-backup-$(date +%s)" --from=cronjob/tarsnap-backup-daily
