.PHONY: build run docker docker-up docker-down test lint vuln

build:
	go build -o bin/dvt-node ./cmd/dvt-node

run: build
	./bin/dvt-node --validator-api 127.0.0.1:4600 --monitoring 127.0.0.1:4620

docker:
	docker build -t aequa-local:latest .

docker-up: docker
	docker compose up -d

docker-down:
	docker compose down -v

test:
	go test ./...

vuln:
	govulncheck ./...
