package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
	"github.com/jackc/pgx/v5"
)

func (m *DB) AddGauges(ctx context.Context, gauges map[string]float64) error {
	batch := &pgx.Batch{}
	for i, v := range gauges {
		// value := strconv.FormatFloat(v, 'f', 6, 64)
		value := fmt.Sprint(v)
		batch.Queue(`
			INSERT INTO metrics (id, type, value)
			VALUES($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET value = $3`, i, "gauge", value)
	}
	return m.db.SendBatch(ctx, batch).Close()
}

func (m *DB) AddCounters(ctx context.Context, counters map[string]int64) error {
	batch := &pgx.Batch{}

	for i, d := range counters {
		delta := strconv.FormatFloat(float64(d), 'f', 6, 64)
		batch.Queue(`
			INSERT INTO metrics (id, type, value)
			VALUES($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET value = EXCLUDED.value + metrics.value`, i, "counter", delta)
	}
	return m.db.SendBatch(ctx, batch).Close()
}

func (m *DB) AddGauge(ctx context.Context, name string, value float64) (float64, error) {
	mValue := strconv.FormatFloat(value, 'f', 6, 64)

	_, err := m.db.Exec(ctx,
		`INSERT INTO metrics (id, type, value) VALUES($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET value = $3`, name, "gauge", mValue,
	)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return value, nil
}

func (m *DB) AddCounter(ctx context.Context, name string, delta int64) (int64, error) {
	mDelta := strconv.FormatInt(delta, 10)

	_, err := m.db.Exec(ctx,
		`INSERT INTO metrics (id, type, value) VALUES($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value + metrics.value`, name, "counter", mDelta,
	)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	var result int64
	rows := m.db.QueryRow(ctx, "SELECT value FROM metrics WHERE type=$1 AND id=$2", "counter", name)

	err = rows.Scan(&result)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return result, nil
}

func (m *DB) ReadGauge(ctx context.Context, name string) (float64, bool) {
	var value float64
	rows := m.db.QueryRow(ctx, "SELECT value FROM metrics WHERE type=$1 AND id=$2", "gauge", name)

	err := rows.Scan(&value)

	if err != nil {
		log.Println(err)
		return value, false
	}

	return value, true
}

func (m *DB) ReadCounter(ctx context.Context, name string) (int64, bool) {
	var delta int64
	rows := m.db.QueryRow(ctx, "SELECT value FROM metrics WHERE type=$1 AND id=$2", "counter", name)

	err := rows.Scan(&delta)

	if err != nil {
		log.Println(err)
		return delta, false
	}

	return delta, true
}

func (m *DB) ReadAllGauges(ctx context.Context) (map[string]float64, error) {
	gauges := make(map[string]float64)

	rows, err := m.db.Query(ctx, `SELECT id, value FROM metrics WHERE type = $1`, "gauge")
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.Value)
		if err != nil {
			log.Println(err)
			return gauges, err
		}
		gauges[m.ID] = *m.Value
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return gauges, err
	}

	return gauges, nil
}

func (m *DB) ReadAllCounters(ctx context.Context) (map[string]int64, error) {
	counters := make(map[string]int64)

	rows, err := m.db.Query(ctx, `SELECT id, value FROM metrics WHERE type = $1`, "counter")
	if err != nil {
		log.Println("###ReadAllCounters 1--->", err)
		log.Println(err)
		panic(err)
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.Delta)
		if err != nil {
			log.Println("###ReadAllCounters 2--->", err)
			return counters, err
		}
		counters[m.ID] = *m.Delta
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		log.Println("###ReadAllCounters 3--->", err)
		log.Println(err)
		return counters, err
	}

	// fmt.Println(counters)
	return counters, nil
}
