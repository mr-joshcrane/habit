#!/bin/bash

protoc proto/*.proto --go_out=proto/
protoc proto/*.proto --go-grpc_out=proto/

