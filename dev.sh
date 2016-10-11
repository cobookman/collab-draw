#!/bin/bash

# The dev project Id environment project id
PROJECT_ID="strong-moose"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

# Build go binary, quit if this fails
# This step isn't necessary, but its useful for
# debugging to fail fast here vs deal with parsing
# docker output when the go build tool would have more
# useful output
go build
if [ $? -ne 0 ]; then
  echo "Failed to build go binary"
  exit
fi

# Stop previous dev container
docker stop collabdraw-dev 2>&1 > /dev/null

# Delete old dev container (if it exists)
OLD_DOCKER=$(docker ps -a --filter 'name=collabdraw-dev' | awk '{if(NR>1) print $1 }')
if [ -n "$OLD_DOCKER" ]; then
  echo "Cleaning up old docker dev container"
  docker rm $OLD_DOCKER
fi

# Build docker image
docker build -t collabdraw .
if [ $? -ne 0 ]; then
  echo "Failed to build container"
  exit
fi

# Run newly built docker image
echo "Running docker image"
docker run -i -t -d -p 8080:8080 -p 65080:65080 --name collabdraw-dev collabdraw

# python -mwebbrowser http://localhost:8080 > /dev/null

