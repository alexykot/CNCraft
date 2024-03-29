package main

import (
	"fmt"

	"github.com/alexykot/cncraft/core"
	"github.com/alexykot/cncraft/core/control"
)

func main() {
	// TODO instantiate using cobra and with relevant flags. Provide via flags only the params
	//  that require server restart, e.g. host:port. Provide other params via the config.
	conf := control.GetDefaultConfig()
	server, err := core.NewServer(conf)
	if err != nil {
		println(fmt.Errorf("failed to instantiate server: %v", err))
		return
	}
	if err = server.Start(); err != nil {
		println(fmt.Errorf("failed to start server: %v", err))
		return
	}
}
