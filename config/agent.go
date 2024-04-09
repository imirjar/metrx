package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	SECRET         string
	URL            string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

// set params from local environment
func (c *AgentConfig) setEnv() {
	//if address in env
	if a := os.Getenv("ADDRESS"); a != "" {
		c.URL = fmt.Sprint("http://", a)
	}
	//if reportInterval in env
	if r := os.Getenv("REPORT_INTERVAL"); r != "" {
		rInt, err := strconv.Atoi(r)
		if err != nil {
			log.Fatal(err)
		}
		c.ReportInterval = time.Duration(1_000_000_000 * rInt) //Nanoseconds to seconds
	}
	//if pollInterval in env
	if p := os.Getenv("POLL_INTERVAL"); p != "" {
		pInt, err := strconv.Atoi(p)
		if err != nil {
			log.Fatal(err)
		}
		c.PollInterval = time.Duration(1_000_000_000 * pInt) //Nanoseconds to seconds
	}
	if k := os.Getenv("SECRET"); k != "" {
		c.SECRET = k
	}
}

// set params from os.Args[]
func (c *AgentConfig) setFlags() {
	a := flag.String("a", "", "api adress")
	p := flag.Int("p", 0, "collect interval")
	r := flag.Int("r", 0, "sending interval")
	k := flag.String("k", "", "SHA-256 hash key")

	flag.Parse()

	if *a != "" {
		c.URL = fmt.Sprint("http://", *a)
	}
	if *r != 0 {
		c.ReportInterval = time.Duration(1_000_000_000 * *r)
	}
	if *p != 0 {
		c.PollInterval = time.Duration(1_000_000_000 * *p)
	}
	if *k != "" {
		log.Println("MY LOVELY KEY", *k)
		c.SECRET = *k
	}

}
