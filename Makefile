# Makefile

.PHONY: generate

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./api/proto/imagedata.proto

.PHONY: docker

docker:
	docker build -f build/docker/client/Dockerfile -t transferspeedclient .
	docker build -f build/docker/server/Dockerfile -t transferspeedserver .
