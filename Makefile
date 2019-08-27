
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BINARY=gts
VERSION=1.1
BUILD=`git rev-list HEAD --max-count=1 --abbrev-commit`
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

bench:
	go test -bench=.  ./terminology

test:
	@go test ./...

test-nc:
	@go test ./... -count=1

build:
	@go build $(LDFLAGS) -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -v -o $(BINARY)-v${VERSION}--$(GOOS)-$(GOARCH))))

pack: build_all
	docker build -t gcr.io/go-terminology/gts:$(VERSION)-$(BUILD) .

push: pack
	docker push gcr.io/go-terminology/gts:$(VERSION)-$(BUILD)

run-container: pack
	docker run gcr.io/go-terminology/gts:$(VERSION)-$(BUILD)

update: 
	@git submodule update --init --recursive

clean:
	@$(RM) ${BINARY}
	@find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

clean-db:
	@$(RM) -r snomed.db
