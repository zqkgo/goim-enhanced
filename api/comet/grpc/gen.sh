#! /bin/bash

protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/gogoproto --gogo_out=plugins=grpc:. api.proto
