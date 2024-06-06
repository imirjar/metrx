package storage

import (
	"context"
	"log"
	"time"

	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/internal/storage/database"
	"github.com/imirjar/metrx/internal/storage/memory"
	"github.com/imirjar/metrx/pkg/ping"
)

type IStorage interface {
	AddGauges(ctx context.Context, gauges map[string]float64) error
	AddCounters(ctx context.Context, counters map[string]int64) error
	AddGauge(ctx context.Context, name string, value float64) (float64, error)
	AddCounter(ctx context.Context, name string, delta int64) (int64, error)
	ReadGauge(ctx context.Context, name string) (float64, bool)
	ReadCounter(ctx context.Context, name string) (int64, bool)
	ReadAllGauges(ctx context.Context) (map[string]float64, error)
	ReadAllCounters(ctx context.Context) (map[string]int64, error)
	//new
	AddMetrics(ctx context.Context, metrics []models.Metrics) error
	AddMetric(ctx context.Context, metric models.Metrics) error
	ReadMetrics(ctx context.Context, mType string) ([]models.Metrics, error)
	ReadMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error)
}

// Create storage layer
// if db connection string from config isn't ok,
// then create "in memory" storage
func NewStorage(DBConn string, filePath string, interval time.Duration, autoImport bool) IStorage {
	if DBConn != "" {
		log.Println("NOTNULLDBCONN")
		db, err := ping.NewDBPool(context.Background(), DBConn)
		if err != nil {
			log.Println(err)
		}
		if err = db.Ping(context.Background()); err != nil {
			log.Println("DBCONN PING ERROR")
		} else {
			return database.NewDB(DBConn)
		}
	}
	store := memory.InMemoryStorage(filePath)
	store.Configure(filePath, autoImport, interval)
	return store
}
