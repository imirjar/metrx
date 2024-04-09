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
}

type AppConfig struct {
	SECRET string
	URL    string
}

type ServiceConfig struct {
	Interval   time.Duration
	FilePath   string
	AutoImport bool
	DBConn     string
}

func (s *ServerConfig) setEnv() {
	//APP START ON ADDRESS
	if a := os.Getenv("ADDRESS"); a != "" {
		s.AppConfig.URL = a
	}

	//BACKUP DATA WHEN MOCK STORAGE
	if f := os.Getenv("FILE_STORAGE_PATH"); f != "" {
		s.ServiceConfig.FilePath = f
	}

	//AUTO RESTORE FROM BACKUP
	if r := os.Getenv("RESTORE"); r != "" {
		rf, err := strconv.ParseBool(r)
		if err != nil {
			panic(err)
		}
		s.ServiceConfig.AutoImport = rf
	}

	//AUTO RESTORE FROM BACKUP
	if i := os.Getenv("STORE_INTERVAL"); i != "" {
		intI, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		s.ServiceConfig.Interval = time.Duration(1_000_000_000 * intI)
	}
	//DATABASE CONNECTION STRING
	if d := os.Getenv("DATABASE_DSN"); d != "" {
		s.ServiceConfig.DBConn = d
	}

	if k := os.Getenv("SECRET"); k != "" {
		s.AppConfig.SECRET = k
	}
}

// set params from os.Args[]
func (s *ServerConfig) setFlags() {
	a := flag.String("a", "", "executable port")
	f := flag.String("f", "", "path to the file where the data should be saved")
	r := flag.Bool("r", s.AutoImport, "if true -> load data from the backup File into storage automatically at startup")
	i := flag.Int("i", -1, "frequency at which data should be saved")
	d := flag.String("d", "", "string postgresql connection")
	k := flag.String("k", "", "SHA-256 hash key")

	flag.Parse()

	if *a != "" {
		s.AppConfig.URL = *a
	}
	if *f != "" {
		s.ServiceConfig.FilePath = *f
	}
	if *r != s.AutoImport {
		s.ServiceConfig.AutoImport = *r
	}
	if *i != -1 {
		s.ServiceConfig.Interval = time.Duration(1_000_000_000 * *i)
	}
	if *d != "" {
		s.ServiceConfig.DBConn = *d
	}
	if *k != "" {
		// log.Println("MY LOVELY KEY", *k)
		s.AppConfig.SECRET = *k
	}
}
