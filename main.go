package main

import (
	"errors"
	"os"

	"github.com/nfwGytautas/ezbin/daemon"
	"github.com/nfwGytautas/ezbin/shared"
	"github.com/nfwGytautas/ezbin/user"
)

func main() {
	// Check if we have a `ezbin.yaml` file either in the current directory or in the arguments
	if shared.ArrayContains(os.Args, "--config") {
		// Config only needed for `daemon` mode
		daemon.Entry()
		return
	}

	// Check if we have a config in the current directory
	info, err := os.Stat("ezbin.yaml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			user.Entry()
			return
		}

		// Some other error
		panic(err)
	}

	// Check if the file is a directory
	if info.IsDir() {
		panic(errors.New("ezbin.yaml is a directory"))
	}

	// Config only needed for `daemon` mode
	daemon.Entry()
}
