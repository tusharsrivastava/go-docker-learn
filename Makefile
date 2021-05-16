build:
	docker build . -t go-docker-learn
	docker image prune -f
