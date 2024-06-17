package database

import (
	"context"
	"log"

	"github.com/imirjar/metrx/internal/models"
	"github.com/jackc/pgx/v5"
)

func (m *DB) AddMetrics(ctx context.Context, metrics []models.Metrics) error {

	batch := &pgx.Batch{}
	for _, m := range metrics {
		// log.Println("AddGauges-->", i)
		// log.Println("AddGauges v -->", v)
		value, err := m.GetVal()
		if err != nil {
			return err
		}
		// log.Println("AddGauges value -->", value)
		// value := fmt.Sprintln(v)
		batch.Queue(`
			INSERT INTO metrics (id, type, value)
			VALUES($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET value = $3`, m.ID, m.MType, value)
	}

	if err := m.db.SendBatch(ctx, batch).Close(); err != nil {
		return err
	}

	return nil
}

func (m *DB) AddMetric(ctx context.Context, metric models.Metrics) error {

	value, err := metric.GetVal()
	if err != nil {
		return err
	}
	_, err = m.db.Exec(ctx,
		`INSERT INTO metrics (id, type, value) VALUES($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET value = $3`, metric.ID, metric.MType, value,
	)
	if err != nil {
		return errAddGaugeExecError
	}
	return nil
}

func (m *DB) ReadMetrics(ctx context.Context, mType string) ([]models.Metrics, error) {

	metrics := make([]models.Metrics, 30)

	rows, err := m.db.Query(ctx, `SELECT id, value FROM metrics WHERE type = $1`, mType)
	if err != nil {
		return metrics, errReadAllGaugesQueryError
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.Value)
		if err != nil {
			return metrics, errReadAllGaugesScanError
		}
		metrics = append(metrics, m)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (m *DB) ReadMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {

	rows := m.db.QueryRow(ctx, "SELECT value FROM metrics WHERE type=$1 AND id=$2", metric.MType, metric.ID)

	err := rows.Scan(&metric.Value)

	if err != nil {
		log.Println("STORAGE ReadGauge ERROR", err)
		return metric, err
	}

	return metric, nil
}
