package config

import "time"

var Testcfg ServerConfig = ServerConfig{
	AppConfig: AppConfig{
		URL: "localhost:8080",
	},
	ServiceConfig: ServiceConfig{
		Interval:   time.Duration(1_000_000_000 * 300), //2s
		FilePath:   "/tmp/metrics-db.json",
		AutoImport: true,
	},
}
