package main

import (
	"flag"
	"log"

	"github.com/LeReverandNox/GuessWhat/src/daemon"
)

func parseArgs() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "The host to bind to")
	flag.StringVar(&cfg.Port, "port", "3000", "The port to bind to")

	flag.Parse()
	return cfg
}

func main() {
	cfg := parseArgs()

	err := daemon.Run(cfg)
	if err != nil {
		log.Fatalf("Something went wrong during launch : %v", err)
	}
}
