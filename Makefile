commit = `git rev-parse --short HEAD`
version = `git describe --abbrev=0`
all:
	go build -ldflags "-X main.Commit $(commit) -X main.Version $(version)"

docker:
	make
	make docker-build
	make docker-push

docker-build:
	docker build -t thundergroove/stats-bot .
docker-push:
	docker push thundergroove/stats-bot
docker-run:
	docker run --name test
