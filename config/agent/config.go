package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/imirjar/metrx/internal/models"
)

// Create agent app configuration with:
// 1. config file params
// 2. local environment params
// 3. os.Args[] params
func NewAgentConfig() *AgentConfig {
	cfg := AgentConfig{}
	cfg.setFileEnv()
	cfg.setEnv()
	cfg.setFlags()

	return &cfg
}

type AgentConfig struct {
	Addr           string          `env:"ADDRESS" json:"address"`
	Secret         string          `env:"SECRET"`
	CryptoKey      PKey            `env:"CRYPTO_KEY" json:"crypto_key"`
	PollInterval   models.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval models.Duration `env:"REPORT_INTERVAL" json:"report_interval" `

	Host string
}

// set params from file environment
func (ac *AgentConfig) setFileEnv() {
	configFile, err := os.ReadFile("config/agent/config.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse JSON into AgentConfig
	err = json.Unmarshal(configFile, &ac)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
}

// set params from local environment
func (ac *AgentConfig) setEnv() {
	if err := env.Parse(ac); err != nil {
		log.Printf("%+v\n", err)
	}
	log.Print(ac.PollInterval)
}

// set params from os.Args[]
func (ac *AgentConfig) setFlags() {
	a := flag.String("a", "", "api adress")
	p := flag.Int("p", 0, "collect interval")
	r := flag.Int("r", 0, "sending interval")
	k := flag.String("k", "", "SHA-256 hash key")
	// cryptoKey := flag.String("crypto-key", "", "crypto key")

	flag.Parse()

	if *a != "" {
		ac.Addr = fmt.Sprint("http://", *a)
	}
	if *k != "" {
		ac.Secret = *k
	}
	// if *cryptoKey != "" {
	// 	rsa, err := encrypt.GetRSA(*cryptoKey)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	ac.CryptoKey = rsa
	// }

	if *r != 0 {
		ac.ReportInterval.Duration = time.Duration(1_000_000_000 * *r)
	}
	if *p != 0 {
		ac.PollInterval.Duration = time.Duration(1_000_000_000 * *p)
	}
}
