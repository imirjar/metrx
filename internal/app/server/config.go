package server

import (
	"flag"
	"log"
	"os"
)

type config struct {
	url string
}

func newConfig() *config {
	cfg := config{
		url: "localhost:8080",
	}
	cfg.setEnv()
	cfg.setFlags()
	log.Print("start on ", cfg.url)
	return &cfg
}

// set params from local environment
func (c *config) setEnv() {
	//if address in env
	if a := os.Getenv("ADDRESS"); a != "" {
		c.url = a
	}
}

// set params from os.Args[]
func (c *config) setFlags() {
	a := flag.String("a", "", "executable port")
	flag.Parse()
	if *a != "" {
		c.url = *a
	}

}
