# FROM golang:onbuild
FROM golang:latest

# Add source
ADD . /go/src/github.com/cobookman/collabdraw

# Install deps & Build service
RUN go get github.com/cobookman/collabdraw

# Run the collabdraw command by default when the container starts.
ENTRYPOINT /go/bin/collabdraw

# Document that the service listens on port 8080.
EXPOSE 8080 65080
