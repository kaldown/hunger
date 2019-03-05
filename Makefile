BUILD=./build
PROTO=./proto
PORT=5001

run:
	docker-compose up -d postgres
	PORT=${PORT} go run server/server.go

stop:
	docker-compose down

build:
	rm -rf ${BUILD}
	mkdir -p ${BUILD} ${BUILD}/bin
	protoc --go_out=plugins=grpc:${BUILD} ${PROTO}/*.proto
	go build -o ${BUILD}/bin/server server/server.go
	go build -o ${BUILD}/bin/client client/client.go

install:
	go get -u google.golang.org/grpc
	go install github.com/golang/protobuf/protoc-gen-go
	pip install -U docker-compose

.PHONY: build
