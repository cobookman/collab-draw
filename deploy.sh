#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/config.sh

# Deploy each service
for d in $DIR/services/* ; do
  cd $d
  aedeploy gcloud app deploy --project ${PROJECT_ID} --no-promote
done
