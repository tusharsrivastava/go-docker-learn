build:
	docker build . -t go-docker-learn
	docker image prune -f

build-compose:
	docker-compose build
