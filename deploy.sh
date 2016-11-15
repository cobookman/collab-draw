#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/config.sh

# Deploy index.yaml
gcloud app deploy index.yaml

# Deploy each service
for d in $DIR/services/* ; do
  cd $d
  ./deploy.sh
  #aedeploy gcloud app deploy --project ${PROJECT_ID} --no-promote
done
