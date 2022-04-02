package todoapp

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"
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
			"./_build": write: contents: actions.build.contents.output
		}
	}
	actions: {
		deps: docker.#Build & {
			steps: [
				alpine.#Build & {
					packages: {
						bash: {}
						go: {}
					}
				},
				docker.#Copy & {
					contents: client.filesystem."./".read.contents
					dest:     "./habit"
				},
			]
		}

		test: bash.#Run & {
			input:   deps.output
			workdir: "./habit"
			script: contents: #"""
				go test ./...
				"""#
		}

		build: {
			run: bash.#Run & {
				input:   test.output
				workdir: "/src"
				script: contents: #"""
					./generate_proto.sh
					"""#
			}

			contents: core.#Subdir & {
				input: run.output.rootfs
				path:  "/src/build"
			}
		}
	}
}
