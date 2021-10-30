docker_tags = v0.0.4
server_main = cmd/server/main.go
server_bin = bin/server
config_path = configs/.env

.PHONY: build docker docker-scratch get mod-tidy mod-download compose compose-pull

build: get
	go build -o $(server_bin) $(server_main)

run: get
	go run $(server_main) --config $(config_path)

get: mod-tidy

mod-tidy: mod-download
	go mod tidy

mod-download:
	go mod download -x

docker:
	docker build ./docker/Dockerfile.alpine --tag $(docker_tags)

docker-scratch:
	docker build ./docker/Dockerfile --tag $(docker_tags)

compose:
	docker compose -f deployments/docker-compose.prod.yml -f deployments/docker-compose.yml up -d -p go-react-app
