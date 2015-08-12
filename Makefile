all:
	go build

docker-build:
	docker build -t thundergroove/stats-bot .
docker-push:
	docker push thundergroove/stats-bot
docker-run:
	docker run --name test
