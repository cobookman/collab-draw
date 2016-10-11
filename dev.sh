#!/bin/bash

# The dev project Id environment project id
PROJECT_ID="strong-moose"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

go build
if [ $? -ne 0 ]; then
  echo "Failed to build go binary"
  exit
fi

killall -9 collabdraw 2> /dev/null
GCLOUD_DATASET_ID=${PROJECT_ID} IS_DEV="1" bash -c './collabdraw' &

python -mwebbrowser http://localhost:8080 > /dev/null

