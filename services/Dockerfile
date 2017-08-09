FROM golang:1.9

ARG service
ENV service $service

ADD ./services/$service $GOPATH/src/github.com/zalfonse/$service

WORKDIR $GOPATH/src/github.com/zalfonse/$service
RUN go get
RUN go build

ENTRYPOINT ./$service
