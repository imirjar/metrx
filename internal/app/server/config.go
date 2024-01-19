package server

import (
	"flag"
	"os"
)

type config struct {
	url string
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
