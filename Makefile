docker_tags = v0.0.4
server_main = cmd/server/main.go
server_bin = bin/server

.PHONY: build docker docker_scratch get

build: get
	go build -o $(server_bin) $(server_main)

get: mod-tidy

mod-tidy: mod-download
	go mod tidy

mod-download:
	go mod download -x

docker:
	docker build ./docker/Dockerfile.alpine --tag $(docker_tags)

docker_scratch:
	docker build ./docker/Dockerfile --tag $(docker_tags)
