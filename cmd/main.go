package main

import (
	"fmt"
	"habit"
	"habit/stores/pbfilestore"
	// "habit/stores/networkstore"
	"os"
)
func main () {
	defaultPath := "habit"
	// s, err := networkstore.Open(defaultPath)
	s, err := pbfilestore.Open(defaultPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	habit.RunCLI(s)
}