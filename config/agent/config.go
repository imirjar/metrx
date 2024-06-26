package config

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env"
)

// Create agent app configuration with:
// 1. config file params
// 2. local environment params
// 3. os.Args[] params
func NewAgentConfig() *AgentConfig {
	cfg := AgentConfig{

		Addr:           "http://localhost:8080",
		PollInterval:   time.Duration(1_000_000_000 * 2),  //2s
		ReportInterval: time.Duration(1_000_000_000 * 10), //10s
	}
	cfg.setEnv()
	cfg.setFlags()

	return &cfg
}

type AgentConfig struct {
	Addr           string        `env:"ADDRESS"`
	Secret         string        `env:"SECRET"`
	CryptoKey      string        `env:"CRYPTO_KEY"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

// set params from local environment
func (c *AgentConfig) setEnv() {
	if err := env.Parse(c); err != nil {
		log.Printf("%+v\n", err)
	}
}

// set params from os.Args[]
func (c *AgentConfig) setFlags() {
	a := flag.String("a", "", "api adress")
	p := flag.Int("p", 0, "collect interval")
	r := flag.Int("r", 0, "sending interval")
	k := flag.String("k", "", "SHA-256 hash key")
	cryptoKey := flag.String("crypto-key", "", "crypto key")

	flag.Parse()

	if *a != "" {
		c.Addr = fmt.Sprint("http://", *a)
	}
	if *k != "" {
		c.Secret = *k
	}
	if *cryptoKey != "" {
		c.CryptoKey = *k
	}

	if *r != 0 {
		c.ReportInterval = time.Duration(1_000_000_000 * *r)
	}
	if *p != 0 {
		c.PollInterval = time.Duration(1_000_000_000 * *p)
	}
}
