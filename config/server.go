package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	AppConfig
	ServiceConfig
	StorageConfig
}

func (s *ServerConfig) setEnv() {
	//if address in env
	if a := os.Getenv("ADDRESS"); a != "" {
		s.AppConfig.URL = a
	}

	if f := os.Getenv("FILE_STORAGE_PATH"); f != "" {
		s.StorageConfig.FilePath = f
	}

	if r := os.Getenv("RESTORE"); r != "" {
		rf, err := strconv.ParseBool(r)
		if err != nil {
			panic(err)
		}
		s.StorageConfig.AutoImport = rf
	}

	if i := os.Getenv("STORE_INTERVAL"); i != "" {
		intI, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		s.ServiceConfig.Interval = time.Duration(1_000_000_000 * intI)
	}

	if d := os.Getenv("DATABASE_DSN"); d != "" {
		s.StorageConfig.DBConn = d
	}
}

// set params from os.Args[]
func (s *ServerConfig) setFlags() {
	a := flag.String("a", "", "executable port")
	f := flag.String("f", "", "path to the file where the data should be saved")
	r := flag.Bool("r", s.AutoImport, "if true -> load data from the backup File into storage automatically at startup")
	i := flag.Int("i", -1, "frequency at which data should be saved")
	d := flag.String("d", "", "string postgresql connection")

	flag.Parse()

	if *a != "" {
		s.AppConfig.URL = *a
	}
	if *f != "" {
		s.StorageConfig.FilePath = *f
	}
	if *r != s.AutoImport {
		s.StorageConfig.AutoImport = *r
	}
	if *i != -1 {
		s.ServiceConfig.Interval = time.Duration(1_000_000_000 * *i)
	}
	if *d != "" {
		s.StorageConfig.DBConn = *d
	}
}

type AppConfig struct {
	URL string
}

type StorageConfig struct {
	FilePath   string
	AutoImport bool
	DBConn     string
}

type ServiceConfig struct {
	Interval time.Duration
}
