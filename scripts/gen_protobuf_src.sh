#!/bin/bash

# INSTALL GRPC:
# apt install -y protobuf-compiler
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
GENDIR=$DIR/../server/rpc

cd $GENDIR && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative msapi.proto
