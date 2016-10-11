#!/bin/bash
PROJECT_ID="strong-moose"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

echo "Getting all go deps"
go get -u ./...

echo "Allowing websocket tcp traffic on port 65080 for all compute nodes tagged 'websocket'"
gcloud compute firewall-rules create default-allow-websockets \
  --allow tcp:65080 \
  --target-tags websocket \
  --description "Allow websocket traffic on port 65080" \
  --project=${PROJECT_ID}

