package app

import (
	"context"
	"log"
	"sync"
	"time"

	config "github.com/imirjar/metrx/config/agent"
	"github.com/imirjar/metrx/internal/agent/client"
	"github.com/imirjar/metrx/internal/agent/system"
	"github.com/imirjar/metrx/internal/models"
)

type Client interface {
	POST(context.Context, []models.Metrics) error
}

type System interface {
	Collect(context.Context) ([]models.Metrics, error)
	// POST(context.Context, models.Batch) error
}

type AgentApp struct {
	client Client
	system System
	sync.Mutex
}

func Run() {
	// Application configuration variables
	cfg := config.NewAgentConfig()

	// HTTP lient for sending metrix to host
	client := client.NewClient(cfg.Secret, cfg.Addr)
	// OS heap
	system := system.NewSystem()

	app := AgentApp{
		client: client,
		system: system,
	}

	if err := app.Start(cfg.PollInterval, cfg.ReportInterval); err != nil {
		panic(err)
	}

}

func (app *AgentApp) Start(p, r time.Duration) error {
	poll := time.NewTicker(p)
	report := time.NewTicker(r)

	var metrics []models.Metrics
	var err error
	var mute sync.RWMutex

	for {
		select {
		case <-poll.C:
			mute.Lock()
			// log.Println("Collect")
			metrics, err = app.system.Collect(context.Background())
			if err != nil {
				log.Println(err)
				return err
			}

			mute.Unlock()
			// log.Print(metrics)
		case <-report.C:
			mute.RLock()
			// log.Println("Send", metrics)

			err := app.client.POST(context.Background(), metrics)
			if err != nil {
				log.Println(err)
				return err
			}
			mute.RUnlock()
		}
	}
	// return nil
}
