package main

import (
	"github.com/imirjar/metrx/internal/app/agent"
)

func main() {
	if err := agent.Run(); err != nil {
		panic(err)
	}
}
