package main

import (
	"time"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/app/server"
)

func main() {
	cfg := config.NewServerConfig()
	serverApp := server.NewServerApp()

	defer serverApp.Service.Storage.Export(cfg.Opts.DumpPath)
	if cfg.Opts.BackupInterval > 0 {
		go func() {
			for {
				time.Sleep(cfg.Opts.BackupInterval)
				serverApp.Service.Storage.Export(cfg.Opts.DumpPath)
			}
		}()
	}

	if err := serverApp.Run(cfg.URL, cfg.Opts.DumpAutoImport, cfg.Opts.DumpPath, cfg.Opts.BackupInterval); err != nil {
		panic(err)
	}
}
