package main

import (
	"fmt"
	"os"

	"github.com/jodydadescott/shelly-client/cmd"
)

func main() {
	err := cmd.NewCmd().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
