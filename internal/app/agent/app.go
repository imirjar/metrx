package agent

import (
	"time"

	"github.com/imirjar/metrx/internal/service"
)

// var memStats runtime.MemStats

func Run() error {

	conf := newConfig()
	agent := service.NewAgent()

	poll := time.NewTicker(conf.pollInterval)
	report := time.NewTicker(conf.reportInterval)

	// wg.Add(1)

	// go func() {
	// 	for {
	// 		agent.CollectMetrix()
	// 		time.Sleep(conf.pollInterval)
	// 	}
	// }()

	// time.Sleep(conf.pollInterval)

	// go func() {
	// 	for {
	// 		agent.SendMetrix(conf.url)
	// 		time.Sleep(conf.reportInterval)

	// 	}
	// }()

	// wg.Wait()
	// return nil

	for {
		select {
		case <-poll.C:
			agent.CollectMetrix()
		case <-report.C:
			agent.SendMetrix(conf.url)
		}
	}
}
