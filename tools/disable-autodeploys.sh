#!/bin/bash

kubectl -n deploy patch cronjob autodeploy -p '{"spec":{"suspend":true}}'
