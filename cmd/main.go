package main

import (
	"fmt"
	"habit"
	"habit/stores/pbfilestore"
	"os"
)

func main() {
	store, err := pbfilestore.Open("store")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	tracker := habit.NewTracker(store)
	server, err := habit.NewServer(tracker)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	habit.RunCLI(server.Client())
}
