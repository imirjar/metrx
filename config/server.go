package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	URL  string
	Opts DataBackupOptions
}

type DataBackupOptions struct {
	DumpPath       string
	BackupInterval time.Duration
	DumpAutoImport bool
}

// set params from local environment
func (c *ServerConfig) setEnv() {
	//if address in env
	if a := os.Getenv("ADDRESS"); a != "" {
		c.URL = a
	}

	if i := os.Getenv("STORE_INTERVAL"); i != "" {
		c.URL = i
	}

	if f := os.Getenv("FILE_STORAGE_PATH"); f != "" {
		c.URL = f
	}

	if r := os.Getenv("RESTORE"); r != "" {
		rf, err := strconv.ParseBool(r)
		if err != nil {
			panic(err)
		}
		c.Opts.DumpAutoImport = rf
	}
}

// set params from os.Args[]
func (c *ServerConfig) setFlags() {
	a := flag.String("a", "", "executable port")
	i := flag.Int("i", -1, "frequency at which data should be saved")
	f := flag.String("f", "", "path to the file where the data should be saved")
	r := flag.Bool("r", c.Opts.DumpAutoImport, "if true -> load data from the backup File into storage automatically at startup")
	flag.Parse()

	if *a != "" {
		c.URL = *a
	}

	if *i != -1 {
		c.Opts.BackupInterval = time.Duration(1_000_000_000 * *i)
	}

	if *f != "" {

		c.Opts.DumpPath = *f
	}

	if *r != c.Opts.DumpAutoImport {
		c.Opts.DumpAutoImport = *r
	}
}
