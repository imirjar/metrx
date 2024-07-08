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

func NewConfig() *Config {
	cfg := Config{}
	cfg.setFileEnv()
	cfg.setEnv()
	cfg.setFlags()
	return &cfg
}

type Config struct {
	Addr      string `env:"ADDRESS" json:"address"`
	Secret    string `env:"SECRET"`
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	Storage
}

type Storage struct {
	Interval   models.Duration `env:"STORE_INTERVAL" toml:"backup_interval"`
	FilePath   string          `env:"FILE_STORAGE_PATH" toml:"dump_file"`
	AutoImport bool            `env:"RESTORE" toml:"auto_restore"`
	DBConn     string          `env:"DATABASE_DSN" toml:"conn"`
}

// set params from file environment
func (ac *Config) setFileEnv() {
	configFile, err := os.ReadFile("config/server/config.json")
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

func (s *Config) setEnv() {
	if err := env.Parse(s); err != nil {
		log.Printf("%+v\n", err)
	}
}

// set params from os.Args[]
func (s *Config) setFlags() {
	a := flag.String("a", "", "executable port")
	f := flag.String("f", "", "path to the file where the data should be saved")
	r := flag.Bool("r", s.AutoImport, "if true -> load data from the backup File into storage automatically at startup")
	i := flag.Int("i", -1, "frequency at which data should be saved")
	d := flag.String("d", "", "string postgresql connection")
	k := flag.String("k", "", "SHA-256 hash key")
	cryptoKey := flag.String("crypto-key", "", "crypto key")

	flag.Parse()

	if *a != "" {
		s.Addr = *a
	}
	if *k != "" {
		s.Secret = *k
	}
	if *cryptoKey != "" {
		s.CryptoKey = *k
	}

	if *f != "" {
		s.Storage.FilePath = *f
	}
	if *r != s.AutoImport {
		s.Storage.AutoImport = *r
	}
	if *i != -1 {
		s.Storage.Interval.Duration = time.Duration(1_000_000_000 * *i)
	}
	if *d != "" {
		s.Storage.DBConn = *d
	}
}
