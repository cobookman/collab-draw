#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/../../config.sh

# Create topic if not exists
gcloud alpha pubsub topics create $UPSTREAM_DRAWING_TOPIC --project=$PROJECT_ID
gcloud alpha functions deploy handleDrawing --stage-bucket=$BUCKET --trigger-topic=$UPSTREAM_DRAWING_TOPIC --project=$PROJECT_ID
