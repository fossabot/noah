.PHONY: default build protos fresh setup_build_dir
default: build

PROTOS_DIRECTORY = ./protos
BUILD_DIRECTORY = ./bin
PACKAGE = github.com/readystock/noah
EXECUTABLE_NAME = noah

build: protos setup_build_dir
	go build -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

fresh: protos setup_build_dir
	go build -a -x -v -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

setup_build_dir:
	mkdir -p $(BUILD_DIRECTORY)

protos:
	protoc -I=$(PROTOS_DIRECTORY) --go_out=plugins=grpc:./db/system $(PROTOS_DIRECTORY)/schema.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=plugins=grpc:./db/system $(PROTOS_DIRECTORY)/query.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/system $(PROTOS_DIRECTORY)/node.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/system $(PROTOS_DIRECTORY)/account.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/sql $(PROTOS_DIRECTORY)/executor.proto
	ls ./db/sql/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
	ls ./db/system/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'