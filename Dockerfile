FROM golang:1.6

RUN mkdir /go/src/gbot
COPY . /go/src/gbot/

WORKDIR /go/src/gbot/

RUN go get
RUN go build
