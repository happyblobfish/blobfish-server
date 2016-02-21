FROM golang:1.5.3

RUN apt-get update -qq && apt-get install -y build-essential

# Create app directory.
ENV APP_HOME=${GOPATH}/src/github.com/kruszczynki/blobfish-server
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME
ADD . $APP_HOME

# Build dependencies and the ws binary
RUN go get ./...
RUN go build
