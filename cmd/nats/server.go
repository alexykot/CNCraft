package main

import natsd "github.com/nats-io/nats-server/server"

func main() {
	opts := &natsd.Options{
		Port: 4222,
	}

	server := natsd.New(opts)
	server.Start()
}
