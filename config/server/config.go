package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env"
)

func NewConfig() *Config {
	cfg := Config{
		Addr: "localhost:8080",
		Storage: Storage{
			Interval:   time.Duration(1_000_000_000 * 300), //300s
			FilePath:   "/tmp/metrics-db.json",
			AutoImport: true,
			DBConn:     "",
		},
	}
	cfg.setEnv()
	cfg.setFlags()
	return &cfg
}

type Config struct {
	Addr      string `env:"ADDRESS" toml:"addr"`
	Secret    string `env:"SECRET" toml:"secret"`
	CryptoKey string `env:CRYPTO_KEY toml:"crypto"`
	Storage
}

type Storage struct {
	Interval   time.Duration `env:"STORE_INTERVAL" toml:"backup_interval"`
	FilePath   string        `env:"FILE_STORAGE_PATH" toml:"dump_file"`
	AutoImport bool          `env:"RESTORE" toml:"auto_restore"`
	DBConn     string        `env:"DATABASE_DSN" toml:"conn"`
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
		s.Storage.Interval = time.Duration(1_000_000_000 * *i)
	}
	if *d != "" {
		s.Storage.DBConn = *d
	}
}
