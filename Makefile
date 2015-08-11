all:
	go build

docker-build:
	docker build -t stats-bot .
docker-push:

docker-run:
	docker run --name test
