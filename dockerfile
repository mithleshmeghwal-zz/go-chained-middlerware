FROM golang:latest

RUN mkdir -p /go/src/server
WORKDIR /go/src/server

ADD . /go/src/server

RUN go get -v

RUN  go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o server" -command="./server"
