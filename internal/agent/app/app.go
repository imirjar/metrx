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

type AgentApp struct {
	client Client
	sync.Mutex
}

// Combine app layers and run app
func Run() {
	// Application configuration variables
	cfg := config.NewAgentConfig()

	// HTTP lient for sending metrix to host
	client := client.New(cfg.Secret, cfg.Addr, &cfg.CryptoKey.Pub)

	app := AgentApp{
		client: client,
	}

	if err := app.Start(cfg.PollInterval.Duration, cfg.ReportInterval.Duration); err != nil {
		panic(err)
	}
}

func (app *AgentApp) Start(p, r time.Duration) error {
	poll := time.NewTicker(p)
	report := time.NewTicker(r)

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	metrics := make(chan []models.Metrics)
	// var err error
	// var mute sync.RWMutex

	for {
		select {
		case <-poll.C:
			go func() {
				// mute.Lock()
				// log.Println("Collect")
				ms, err := system.Collect(context.Background())
				if err != nil {
					log.Println(err)
					return
				}

				log.Println("отправлено в ms", metrics)

				metrics <- ms
			}()
			// mute.Unlock()
			// log.Print(metrics)
		case <-report.C:
			go func() {
				log.Println("принято из ms", metrics)
				ms := <-metrics
				err := app.client.POST(context.Background(), ms)
				if err != nil {
					log.Println(err)
				}
			}()
			// mute.RLock()

			// ms := <-metrics

			// err = app.client.POST(context.Background(), <-metrics)
			// if err != nil {
			// 	log.Println(err)
			// 	return err
			// }
			// mute.RUnlock()
			// case <-sig:
			// 	// сюда попадем, если вызвали cancelFunc()
			// 	fmt.Printf("worker stopped\n")

			// 	time.Sleep(p)
			// 	poll.Stop()

			// 	time.Sleep(r)
			// 	report.Stop()

			// 	return nil

		}
	}
	// return nil
}
