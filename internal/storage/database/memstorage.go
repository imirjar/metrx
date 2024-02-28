package database

import (
	"context"
	"log"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
	"github.com/jackc/pgx/v5"
)

func (m *DB) AddGauges(gauges map[string]float64) error {
	batch := &pgx.Batch{}
	for i, v := range gauges {
		value := strconv.FormatFloat(v, 'f', 6, 64)
		batch.Queue(`
			INSERT INTO metrics (id, type, value)
			VALUES($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET value = $3`, i, "gauge", value)
	}
	return m.db.SendBatch(context.Background(), batch).Close()
}

func (m *DB) AddCounters(counters map[string]int64) error {
	batch := &pgx.Batch{}
	for i, d := range counters {
		delta := strconv.FormatFloat(float64(d), 'f', 6, 64)
		batch.Queue(`
			INSERT INTO metrics (id, type, value)
			VALUES($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET value = EXCLUDED.value + metrics.value`, i, "counter", delta)
		// fmt.Println(delta)
	}
	return m.db.SendBatch(context.Background(), batch).Close()
}

func (m *DB) AddGauge(name string, value float64) (float64, error) {
	mValue := strconv.FormatFloat(value, 'f', 6, 64)

	_, err := m.db.Exec(context.Background(),
		`INSERT INTO metrics (id, type, value) VALUES($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET value = $3`, name, "gauge", mValue,
	)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return value, nil
}

func (m *DB) AddCounter(name string, delta int64) (int64, error) {

	var curDelta int64

	rows := m.db.QueryRow(context.Background(), "SELECT value FROM metrics WHERE type=$1 AND id=$2", "counter", name)
	err := rows.Scan(&curDelta)
	if err != nil {
		log.Println(err)
	}
	curDelta += delta
	mDelta := strconv.FormatInt(curDelta, 10)

	_, err = m.db.Exec(context.Background(),
		`INSERT INTO metrics (id, type, value) VALUES($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET value = $3`, name, "counter", mDelta,
	)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return curDelta, nil
}

func (m *DB) ReadGauge(name string) (float64, bool) {
	var value float64
	rows := m.db.QueryRow(context.Background(), "SELECT value FROM metrics WHERE type=$1 AND id=$2", "gauge", name)

	err := rows.Scan(&value)

	if err != nil {
		log.Println(err)
		return value, false
	}

	return value, true
}

func (m *DB) ReadCounter(name string) (int64, bool) {
	var delta int64
	rows := m.db.QueryRow(context.Background(), "SELECT value FROM metrics WHERE type=$1 AND id=$2", "counter", name)

	err := rows.Scan(&delta)

	if err != nil {
		log.Println(err)
		return delta, false
	}

	return delta, true
}

func (m *DB) ReadAllGauges() (map[string]float64, error) {
	gauges := make(map[string]float64)

	rows, err := m.db.Query(context.Background(), `SELECT id, value FROM metrics WHERE type = $1`, "gauge")
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

func (m *DB) ReadAllCounters() (map[string]int64, error) {
	counters := make(map[string]int64)

	rows, err := m.db.Query(context.Background(), `SELECT id, value FROM metrics WHERE type = $1`, "counter")
	if err != nil {
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
			return counters, err
		}
		counters[m.ID] = *m.Delta
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return counters, err
	}

	// fmt.Println(counters)
	return counters, nil
}
