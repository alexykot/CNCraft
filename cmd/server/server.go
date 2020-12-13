package main

import (
	"fmt"
	"github.com/alexykot/cncraft/core/control"

	"github.com/fatih/color"

	"github.com/alexykot/cncraft/core"
)

func main() {
	color.NoColor = false

	// TODO instantiate using cobra and with relevant flags. Provide via flags only the params
	//  that require server restart, e.g. host:port. Provide other params via the config.
	server, err := core.NewServer(control.DefaultConfig())
	if err != nil {
		println(fmt.Errorf("failed to start the server: %v", err))
		return
	}
	server.Load()
}
/**/