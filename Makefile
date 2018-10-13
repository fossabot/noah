.PHONY: default build protos
default: build

build:
	protos
	go build

PROTOS_DIRECTORY = ./protos

protos:
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/system $(PROTOS_DIRECTORY)/sequences.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=plugins=grpc:./db/system $(PROTOS_DIRECTORY)/schema.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=plugins=grpc:./db/system $(PROTOS_DIRECTORY)/node.proto
