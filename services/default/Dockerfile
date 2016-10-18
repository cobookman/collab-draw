# FROM golang:onbuild
FROM golang:latest

# Add source
ADD . /go/src/github.com/cobookman/collabdraw/services/default

# Install deps & Build service
RUN go get github.com/cobookman/collabdraw/services/default

# Run the collabdraw command by default when the container starts.
ENTRYPOINT /go/bin/default

# Document that the service listens on port 8080.
EXPOSE 8080
