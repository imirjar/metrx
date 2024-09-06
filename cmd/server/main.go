package main

//go:generate go run main.go

import (
	"fmt"
	_ "net/http/pprof"

	app "github.com/imirjar/metrx/internal/server/app"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf(`Build version: %s
Build date: %s
Build commit: %s 
`, buildVersion, buildDate, buildCommit)

	app.Run()
}
