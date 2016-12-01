# FROM golang:onbuild
FROM golang:latest

# Get shared libraries from github. In realife this would be in its own
# github repo. And this line could be elimited as the next
# RUN line would fetch this dep.
RUN go get github.com/cobookman/collabdraw/shared/models

# Add our code for the drawings micro-service
ADD . /go/src/github.com/cobookman/collabdraw/services/drawings

# Install any deps for the drawings micro-service.
RUN go get github.com/cobookman/collabdraw/services/drawings

# Run the collabdraw command by default when the container starts.
ENTRYPOINT /go/bin/drawings

# Document that the service listens on port 8080.
EXPOSE 8080 65080
