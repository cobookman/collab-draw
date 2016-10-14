#!/bin/bash
PROJECT_ID="strong-moose"
#ADDR="strong-moose.appspot.com:65080"
ADDR="localhost:65080"
ADDR="104.198.178.125:65080"
GCLOUD_DATASET_ID=${PROJECT_ID} IS_DEV="1" bash -c "go test --addr=${ADDR}"

