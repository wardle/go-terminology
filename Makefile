
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
	@go generate ./...

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
