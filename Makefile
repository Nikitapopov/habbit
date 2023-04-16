build:
	go build -o ./.bin/main ./cmd/main/main.go

run: build
	./.bin/main

run-mongodb:
	docker-compose --env-file .env up