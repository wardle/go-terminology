
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BINARY=gts
VERSION=0.1.0
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.build=${BUILD}"

default: generate test build

all: generate test build build_all

generate:
	protoc -Ivendor/terminology/protos --go_out=plugins=gprc:snomed vendor/terminology/protos/snomed.proto
	protoc -Ivendor/terminology/protos -Ivendor/terminology/vendor/googleapis --go_out=plugins=grpc:snomed vendor/terminology/protos/server.proto
	protoc -Ivendor/terminology/protos -Ivendor/terminology/vendor/googleapis --grpc-gateway_out=logtostderr=true:snomed vendor/terminology/protos/server.proto
	protoc -Ivendor/terminology/protos -Ivendor/terminology/vendor/googleapis --swagger_out=logtostderr=true:. vendor/terminology/protos/server.proto

test:
	@go test ./...

build:
	@go build $(LDFLAGS) -o ${BINARY}


build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(BINARY)-$(VERSION)-$(GOOS)-$(GOARCH))))

update: 
	@git submodule update --init --recursive

clean:
	@$(RM) ${BINARY}
	@find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

clean-db:
	@$(RM) -r snomed.db