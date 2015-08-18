FROM golang
MAINTAINER THUNDERGROOVE

ENV DOCKER 1

ADD *.tmpl /assets/

copy stats-bot /bin/
ENTRYPOINT /bin/stats-bot
