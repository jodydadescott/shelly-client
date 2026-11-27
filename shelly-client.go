package main

import (
	"errors"
	"os"

	"github.com/jodydadescott/shelly-client/cmd"
)

func main() {
	err := cmd.NewCmd().Execute()
	if errors.Is(err, cmd.ErrDrift) {
		os.Exit(2)
	}
	if err != nil {
		os.Exit(1)
	}
}
