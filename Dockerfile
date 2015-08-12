FROM golang
MAINTAINER THUNDERGROOVE

#add . /go/src/github.com/THUNDERGROOVE/stats-bot

#WORKDIR /go/src/github.com/THUNDERGROOVE/stats-bot

#RUN go get golang.org/x/text/encoding

#RUN go get ./... # It's magic!

#RUN go install github.com/THUNDERGROOVE/stats-bot


ADD lookup_template.tmpl /assets/

copy stats-bot /bin/
ENTRYPOINT /bin/stats-bot