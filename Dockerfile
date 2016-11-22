FROM golang:1.7

COPY ./ /go/src/github.com/NeowayLabs/statsdig

WORKDIR /go/src/github.com/NeowayLabs/statsdig

RUN go get ./...
RUN go build ./cmd/sender
