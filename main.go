package main

import (
	"flag"
	"log"

	"github.com/gurshaan17/redis/config"
	"github.com/gurshaan17/redis/server"
)

func setupFlags() {
	flag.StringVar(
		&config.Host,
		"host",
		"0.0.0.0",
		"host for the dice server",
	)

	flag.IntVar(
		&config.Port,
		"port",
		7379,
		"port for the dice server",
	)

	flag.Parse()
}

func main() {
	setupFlags()

	log.Println("starting thee server")
	server.RunSyncTCPServer()
}
