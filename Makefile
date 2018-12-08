.PHONY: default build protos
default: build

build:
	protos
	go build

PROTOS_DIRECTORY = ./protos

protos:
	#protoc --proto_path=
	# 	$(GOPATH)/src/github.com/gogo/protobuf/protobuf:
	#	$(GOPATH)/src/github.com/gogo/protobuf/protobuf/google/protobuf --gogo_out=.*.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/system $(PROTOS_DIRECTORY)/sequences.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=plugins=grpc:./db/system $(PROTOS_DIRECTORY)/schema.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/system $(PROTOS_DIRECTORY)/node.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/system $(PROTOS_DIRECTORY)/account.proto
	protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/sql $(PROTOS_DIRECTORY)/executor.proto
	#protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/util/log $(PROTOS_DIRECTORY)/log.proto
	#protoc -I=$(PROTOS_DIRECTORY) --go_out=./db/util/tracing $(PROTOS_DIRECTORY)/recorded_span.proto