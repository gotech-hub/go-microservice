#!/bin/bash

# Script để generate protobuf code
set -e

echo "Generating protobuf code..."

# Tạo thư mục gen nếu chưa tồn tại
mkdir -p api/grpc/proto/gen/user

# Generate Go code từ proto file
protoc --go_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       api/grpc/proto/user.proto

echo "Protobuf code generated successfully!" 