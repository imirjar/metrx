package storage

import (
	"encoding/json"
	"os"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func New() *MemStorage {
	return &MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}

func (m *MemStorage) AddGauge(mName string, mValue float64) {
	m.Gauge[mName] = mValue
}

func (m *MemStorage) AddCounter(mName string, mValue int64) {
	m.Counter[mName] = mValue
}

func (m *MemStorage) ReadAllGauge() map[string]float64 {
	return m.Gauge
}

func (m *MemStorage) ReadAllCounter() map[string]int64 {
	return m.Counter
}

func (m *MemStorage) ReadGauge(mName string) (float64, bool) {
	v, ok := m.Gauge[mName]
	return v, ok
}

func (m *MemStorage) ReadCounter(mName string) (int64, bool) {
	v, ok := m.Counter[mName]
	return v, ok
}

func (m *MemStorage) Export(path string) error {
	// fmt.Println("###")
	// file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
	// if err != nil {
	// 	return err
	// }

	// 	mm, err := json.Marshal(m)
	// 	file.Write(mm)
	// 	file.Close()н

	// 	return nil

	// сериализуем структуру в JSON формат
	data, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		return err
	}
	// сохраняем данные в файл
	return os.WriteFile(path, data, 0666)
}

func (m *MemStorage) Import(path string) error {
	// fmt.Println(path)
	// file, err := os.ReadFile(path)
	// if err != nil {
	// 	return err
	// }
	// // fmt.Println(file)
	// err = json.Unmarshal(file, m)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// return nil
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	return nil
}
