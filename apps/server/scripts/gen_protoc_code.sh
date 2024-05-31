#!/bin/bash
proto_dir=internal/pkg/sim/proto

protoc --proto_path=./${proto_dir} \
  --go_out=./${proto_dir} \
  --go-grpc_out=./${proto_dir} \
  player.proto
