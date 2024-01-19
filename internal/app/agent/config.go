package agent

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type config struct {
	url            string
	pollInterval   time.Duration
	reportInterval time.Duration
}

func newConfig() *config {
	cfg := config{
		url:            "http://localhost:8080",
		pollInterval:   time.Duration(1_000_000_000 * 2),  //2s
		reportInterval: time.Duration(1_000_000_000 * 10), //10s
	}
	cfg.setEnv()
	cfg.setFlags()

	return &cfg
}

// set params from local environment
func (c *config) setEnv() {
	//if address in env
	if a := os.Getenv("ADDRESS"); a != "" {
		c.url = fmt.Sprint("http://", a)
	}
	//if reportInterval in env
	if r := os.Getenv("REPORT_INTERVAL"); r != "" {
		rInt, err := strconv.Atoi(r)
		if err != nil {
			log.Fatal(err)
		}
		c.reportInterval = time.Duration(1_000_000_000 * rInt) //Nanoseconds to seconds
	}
	//if pollInterval in env
	if p := os.Getenv("POLL_INTERVAL"); p != "" {
		pInt, err := strconv.Atoi(p)
		if err != nil {
			log.Fatal(err)
		}
		c.pollInterval = time.Duration(1_000_000_000 * pInt) //Nanoseconds to seconds
	}
}

// set params from os.Args[]
func (c *config) setFlags() {
	a := flag.String("a", "", "executable port")
	p := flag.Int("p", 0, "executable port")
	r := flag.Int("r", 0, "executable port")

	flag.Parse()

	if *a != "" {
		c.url = fmt.Sprint("http://", *a)
	}
	if *r != 0 {
		c.reportInterval = time.Duration(1_000_000_000 * *r)
	}
	if *p != 0 {
		c.pollInterval = time.Duration(1_000_000_000 * *p)
	}

}
