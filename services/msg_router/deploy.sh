#!/bin/bash
BUCKET="strong-moose.appspot.com"
TOPIC="incoming-drawing"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/../../config.sh

# Create topic if not exists
gcloud alpha pubsub topics create $TOPIC --project=$PROJECT_ID
gcloud alpha functions deploy forwardDrawing --stage-bucket=$BUCKET --trigger-topic=$TOPIC --project=$PROJECT_ID
gcloud alpha functions deploy saveDrawing --stage-bucket=$BUCKET --trigger-topic=$TOPIC --project=$PROJECT_ID
