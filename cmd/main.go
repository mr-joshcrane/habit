package main

import (
	"fmt"
	"habit"
	"habit/stores/pbfilestore"
	"os"
)
func main () {
	defaultPath := "habit"
	s, err := pbfilestore.Open(defaultPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	habit.RunCLI(s) 
}