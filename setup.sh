#!/usr/bin/env bash

sleep 1
go get -v -d google.golang.org/grpc
go get -v -d -t github.com/golang/protobuf/...
curl -L https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip -o /tmp/protoc.zip
mkdir -p "$HOME"/protoc
unzip /tmp/protoc.zip -d "$HOME"/protoc
mkdir -p "$HOME"/src && ln -s "$HOME"/protoc "$HOME"/src/protobuf
go get -u github.com/golang/protobuf/protoc-gen-go
export PATH=$HOME/protoc/bin:$GOPATH/bin:$PATH