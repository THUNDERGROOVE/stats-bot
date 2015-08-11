FROM golang

add . /go/src/github.com/THUNDERGROOVE/stats-bot

WORKDIR /go/src/github.com/THUNDERGROOVE/stats-bot

RUN go get golang.org/x/text/encoding

RUN go get ./... # It's magic!

RUN go install github.com/THUNDERGROOVE/stats-bot

ENTRYPOINT /go/bin/stats-bot