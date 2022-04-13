package todoapp

import (
	"dagger.io/dagger"
    "universe.dagger.io/docker"
	"universe.dagger.io/alpine"
	"universe.dagger.io/bash"
)

dagger.#Plan & {
	client: {
		filesystem: {
			"./": read: {
				contents: dagger.#FS
			}
			"./": write: contents: actions.build.contents.output
		}
	}
	actions: {
		deps: docker.#Build & {
			steps: [
				alpine.#Build & {
					packages: {
						bash: {}
						go: {}
						protoc: {}
					}
				},
				docker.#Copy & {
					contents: client.filesystem."./".read.contents
					dest:     "./habit"
				},
				bash.#Run & {
					workdir: "./habit"
					script: contents: #"""
							go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
							go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
						"""#
				},
			]
		}


		build: {
			run: bash.#Run & {
				input:   deps.output
				workdir: "./habit"
				script: contents: #"""
						export PATH="$PATH:$(go env GOPATH)/bin"
						protoc proto/*.proto --go_out=proto/
						protoc proto/*.proto --go-grpc_out=proto/
						go test ./...
					"""#
			}
		}	
	}
}
