package main

import (
	"habit/server"
)

func main() {
	server.ListenAndServe("localhost:8080")
}