package database

import (
	"github.com/imirjar/metrx/internal/models"
)

type gauge struct {
	Id    string
	Value float64
}

type counter struct {
	Id    string
	Value int64
}

// func (m *DB) AddGauge(mName string, mValue float64) (float64, error) {
// 	_, ok := m.ReadGauge(mName)

// 	value := strconv.FormatFloat(mValue, 'f', 6, 64)
// 	if !ok {
// 		_, err := m.db.Exec(
// 			"INSERT INTO metrics (id, type, value)"+
// 				" VALUES($1,$2,$3)", mName, "gauge", value)
// 		if err != nil {
// 			return 0, err
// 		}
// 	} else {
// 		_, err := m.db.Exec(
// 			"UPDATE metrics SET value = $1"+
// 				"WHERE id = $2 AND type = $3)", value, mName, "gauge")
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
// 	return mValue, nil
// }

// func (m *DB) AddCounter(mName string, mValue int64) (int64, error) {
// 	curVal, ok := m.ReadCounter(mName)
// 	value := strconv.FormatInt(mValue, 10)
// 	if !ok {
// 		_, err := m.db.Exec(
// 			"INSERT INTO metrics (id, type, value)"+
// 				" VALUES($1,$2,$3)", mName, "counter", value)
// 		if err != nil {
// 			return 0, err
// 		}
// 	} else {
// 		_, err := m.db.Exec(
// 			"UPDATE metrics SET value = $1"+
// 				"WHERE id = $2 AND type = $3)", value, mName, "counter")
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
// 	return curVal + mValue, nil
// }

// func (m *DB) ReadAllGauge() map[string]float64 {
// 	fmt.Println("ReadAllCounter")

// 	gauges := map[string]float64{}

// 	rows, err := m.db.Query(`SELECT id, value FROM metrics WHERE type = $1`, "gauge")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// обязательно закрываем перед возвратом функции
// 	defer rows.Close()

// 	// пробегаем по всем записям
// 	for rows.Next() {
// 		g := gauge{}
// 		err = rows.Scan(&g.Id, &g.Value)
// 		if err != nil {
// 			panic(err)
// 		}
// 		gauges[g.Id] = g.Value
// 	}

// 	// проверяем на ошибки
// 	err = rows.Err()
// 	if err != nil {
// 		return nil
// 	}

// 	return gauges
// }

// func (m *DB) ReadAllCounter() map[string]int64 {
// 	fmt.Println("ReadAllCounter")

// 	counters := map[string]int64{}

// 	rows, err := m.db.Query(`SELECT id, value FROM metrics WHERE type = $1`, "counter")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// обязательно закрываем перед возвратом функции
// 	defer rows.Close()

// 	// пробегаем по всем записям
// 	for rows.Next() {
// 		c := counter{}
// 		err = rows.Scan(&c.Id, &c.Value)
// 		if err != nil {
// 			panic(err)
// 		}
// 		counters[c.Id] = c.Value
// 	}

// 	// проверяем на ошибки
// 	err = rows.Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return counters
// }

// func (m *DB) ReadGauge(mName string) (float64, bool) {
// 	rows := m.db.QueryRow("SELECT value FROM metrics WHERE type=$1 AND id=$2", "gauge", mName)

// 	c := gauge{
// 		Id: mName,
// 	}
// 	err := rows.Scan(&c.Value)

// 	if err != nil {
// 		return 0, false
// 	}
// 	return c.Value, true
// }

// func (m *DB) ReadCounter(mName string) (int64, bool) {
// 	rows := m.db.QueryRow("SELECT value FROM metrics WHERE type=$1 AND id=$2", "counter", mName)

// 	c := counter{
// 		Id: mName,
// 	}
// 	err := rows.Scan(&c.Value)

// 	if err != nil {
// 		return 0, false
// 	}
// 	return c.Value, true
// }

func (m *DB) AddGauge(metric models.Metrics) error                 { return nil }
func (m *DB) AddCounter(name string, delta int64) error            { return nil }
func (m *DB) ReadOne(metric models.Metrics) (models.Metrics, bool) { return models.Metrics{}, true }
func (m *DB) ReadAllGauges() (map[string]float64, error)           { return map[string]float64{}, nil }
func (m *DB) ReadAllCounters() (map[string]int64, error)           { return map[string]int64{}, nil }
func (m *DB) Delete(metric models.Metrics) error                   { return nil }
