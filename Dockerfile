# FROM golang:onbuild
FROM golang:latest

ENV GCLOUD_DATASET_ID strong-moose
ENV GOOGLE_APPLICATION_CREDENTIALS /go/src/github.com/cobookman/collabdraw/service-account.json

# Add source
ADD . /go/src/github.com/cobookman/collabdraw

# Install deps & Build service
RUN go get github.com/cobookman/collabdraw

# Run the collabdraw command by default when the container starts.
ENTRYPOINT /go/bin/collabdraw

# Document that the service listens on port 8080.
EXPOSE 8080 65080
