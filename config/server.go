package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	URL string
}

// set params from local environment
func (c *ServerConfig) setEnv() {
	//if address in env
	if a := os.Getenv("ADDRESS"); a != "" {
		c.URL = a
	}
}

// set params from os.Args[]
func (c *ServerConfig) setFlags() {
	a := flag.String("a", "", "executable port")
	flag.Parse()
	if *a != "" {
		c.URL = *a
	}
}
