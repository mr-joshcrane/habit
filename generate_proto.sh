#!/bin/bash

protoc -I proto/ --go_out=proto/ proto/store.proto

protoc proto/store.proto --go-grpc_out=proto/